# Week05学习笔记

## commont系统实例

### 读的核心逻辑
1. Cache-Aside 模式，先读取缓存，再读取存储。

2. 早期 cache rebuild 是做到服务里的，对于重建逻辑，
  一般会使用 read ahead 的思路，即预读，用户访 问了第一页，
  很有可能访问第二页，所以缓存会超 前加载，避免频繁 cache miss。

3. 当缓存抖动的时候，特别容易引起集群 hundering herd (惊群)现象,
   大量的请求会触发 cache rebuild,因为使用了预加载,容易导致服务OOM。
   所以我们开到回源的逻辑里，使用消息队列来进行逻辑异步化，对于当前请求
   只返回 mysql 中部分数据即停止。

#### hundering herd (惊群)现象
- 大量的进程在等待某个事件的发生
- 在某个时间点这个事件发生了
- 大量的进程都被叫醒了,但是只有一个进程能获得资源(比如抢一个锁)
- 剩下的进程都重新进入等待状态
- 这个过程不断循环直到最后没有等待的进程了

### 写的核心逻辑
1. 写和读相比较，写可以认为是透穿到存储层的，系统的瓶颈往往就来自于存储层，或者有状态层。

2. 对于写的设计上，我们认为刚发布的评论有极短的延迟(通常小于几 ms)对用户可见是可接受的，
   把对存储的直接冲击下放到消息队列，按照消息反压的思路，即如果存储 latency 升高，消费
   能力就下降，自然消息容易堆积，系统始终以最大化方式消费。 

3. Kafka 是存在 partition(分区) 概念的，可以认为是物理上 的一个小队列。

   一个 topic 是由一组 partition 组成 的，所以 Kafka 的吞吐模型理解为: 
   
   全局并行，局部串行的生产消费方式。
   
   对于入队的消息，可以按照 
   hash(comment_subject) % N(partitions) 的方式进行分发。
   那么某个 partition 中的 评论主题的 数据一定都在一起，这样方便我们串行消费。

4. 同样的，处理回源消息也是类似的思路.

### 表设计的优化

1. 内容表和索引表分开设计，方便 mysql datapage 缓存更多的 row，如果和 context 耦合，
   会导致更大 的 IO。长远来看 content 信息可以直接使用 KV storage 存储。

### 缓存的设计

1. redis sortedest
 - Sorted Set有点像Set和Hash的结合体。
 - 和Set一样，它里面的元素是唯一的，类型是String，所以它可以理解为就是一个Set。
 - 但是Set里面的元素是无序的，而Sorted Set里面的元素都带有一个浮点值，叫做分数（score），
   所以这一点和Hash有点像，因为每个元素都映射到了一个值。
 - Sorted Set是有序的，规则如下：
    如果A.score > B.score，那么A > B。
    如果A.score == B.score，那么A和B的大小就通过比较字符串来决定了，而A和B的字符串是不
    会相等的，因为Sorted Set里面的值都是唯一的。

2. 根据数据的lasttime查询特定条数的数据输出给客户端

3. shardingKey分区Key (Shard Key 碎片Key)
 - 当将数据存储进行分隔成分区时，需靠考虑将哪些数据应该放在哪个分区。一个分区通常需要检索一
   个或者多个一定范围内的数据属性，就由这些属性构成Shard的Key(有时称为分区Key)
   Shard的Key应该是静态的。它不应该基于可能改变的数据。

### 可用性设计 - Singleflight & Doubleflight

1. 起因:对于热门的主题，如果存在缓存穿透的情况，会导致 大量的同进程、跨进程的数据回源到
   存储层，可能会 引起存储过载的情况，如何只交给同进程内，一个程序去做加载存储?

2. 使用归并回源的思路
```
   https://pkg.go.dev/golang.org/x/sync/singleflight

   同进程只交给一个程序去获取 mysql 数据，然后批量返回。
   
   同时这个 lease owner 投递一个 kafka 消息，做 index cache 的 recovery 操作。
   
   这样可以大大减少 mysql 的压力，以及大量透穿导致的密集写 kafka 的问题。
   更进一步的，后续连续的请求，仍然可能会短时 cache miss，我们可以在进程内设置一个
   short-lived flag，标记最近有一个人投递了 cache rebuild 的消息， 直接drop。
```
