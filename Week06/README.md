# 作业
```
  参考 Hystrix 实现一个滑动窗口计数器。

	以上作业，要求提交到 GitHub 上面，Week06 作业提交地址：
	https://github.com/Go-000/Go-000/issues/81

	请务必按照示例格式进行提交，不要复制其他同学的格式，以免格式错误无法抓取作业。
```

# 微服务可用性设计

## 隔离
- 定义：隔离，本质上是对系统或资源进行分割，从而实现当系统发生故障时能限定
  传播范围和影响范围，即发生故障后只有出问题的服务不可用，保证其 他服务仍然可用。

1. 服务隔离 
   • 动静分离、读写分离 
   例如:CQRS，也就是Command Query Responsibility Segregation，故名思义是将 command 与 query 分离的一种模式。
   当 command 系统完成数据更新的操作后，会通过「领域事件」的方式通知 query 系统。query 系统在接受到事件之后更新自己
   的数据源。所有的查询操作都通过 query 系统暴露的接口完成
   - 正向索引  解析文档内的单词，然后建立从文档到词组的映射关系
       解析每个文档出现的单词，然后建立从文档 (document) 到词组 (words) 的映射关系，这就是正向索引
   - 反向索引  建立从单词 (word) 到文档 (document lsit) 的映射关系
       反向索引方向则是正向索引的逆向，建立从单词 (word) 到文档 (document lsit) 的映射关系
- 服务隔离: 
   1.动静隔离: 
    ```
      例如: cpu的cacheline false sharing
      数据库 mysql 表设计中避免 bufferpool(缓冲池) 频繁过期, 隔离 动静表
      大到架构设计中的图片、静态资源等缓 存加速
    ```
      本质上都体现的一样的思路，即加速/缓 存访问变换频次小的
      
      CDN 场景中，将静态 资源和动态 API 分离，也是体现了隔离的思路
      例如:CDN的场景，CDN基本都支持边缘计算，在每个CDN节点都部署一个Agent，把流量聚合后再打回到源站，这样会使源的qps大降的。

    ```
      • 降低应用服务器负载，静态文件访问负载全部通过 CDN。 
      • 对象存储存储费用最低。 
      • 海量存储空间，无需考虑存储架构升级。 
      • 静态CDN带宽加速，延迟低。
    ```

- InnoDB的缓冲池(bufferpool)
      缓存表数据与索引数据，把磁盘上的数据加载到缓冲池，
      避免每次访问都进行磁盘IO，起到加速访问的作用。

- 管理与淘汰缓冲池
      1. 预读: 磁盘读写，并不是按需读取，而是按页读取，一次至少读一页数据（一般是4K），
      如果未来要读取的数据就在页中，就能够省去后续的磁盘IO，提高效率。
      
2. 轻重隔离 
   • 核心、快慢、热点
   - 业务按照Level进行资源池分集
   - 
3. 物理隔离 
   • 线程、进程(容器隔离)、集群、机房
     
4. 热点隔离
   • 热点即经常访问的数据