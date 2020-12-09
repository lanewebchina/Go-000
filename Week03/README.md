# Week03学习笔记

## 本周的作业
- 基于 errgroup 实现一个http server的启动和关闭,
  以及 linux signal 信号的注册和处理,
  要保证能够一个退出，全部注销退出。

## 要看的文章
- Effective Go see https://github.com/bingohuang/effective-go-zh-en
- https://golang.org/ref/mem
- https://www.jianshu.com/p/5e44168f47a3

## Goroutine
- 并行不意味着是并发，并行指的是不同的执行单元。并发指代的不同的执行单元同时执行.
- log.Fatal() 底层调用了os.Exit会导致Goroutine的 defer函数不会被调用到
- Only user log.Fatal from main.main or init functions

### 重要的是记住如下三点
- 把并发交给调用者, 意思是一定是调用者来决定是否启动Goroutine，而不是在函数内部启动
- 搞清楚goroutine什么时候退出，要管控它的生命周期
  一个goroutine的生命周期应该是你自己来管理的,一定要搞清楚它啥时候结束，如何让它结束掉
  要有手段能让它退出
- 能够控制这个goroutine退出，包括context超时, goroutine监听channel的消息

```
  //Tracker knows how to track events for the application
  type Tracker struct {
	  ch   chan string
	  stop chan struct{}
  }

  func main() {
	  tr := NewTracker()
    go tr.Run()
    _ = tr.Event(context.Background(), "test")
    _ = tr.Event(context.Background(), "test")
    _ = tr.Event(context.Background(), "test")
    ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
    defer cancel()
    tr.ShutDown(ctx)
  }

  func NewTracker() *Tracker {
    return &Tracker{
      ch: make(chan string, 10),
    }
  }

  func (t *Tracker) Event(ctx context.Context, data string) error {
    select {
    case t.ch <- data:
      return nil
    case <-ctx.Done():
      return ctx.Err()
    }
  }

  func (t *Tracker) Run() {
    for data := range t.ch {
      time.Sleep(1 * time.Second)
      fmt.Println(data)
    }
    t.stop <- struct{}{}
  }

  func (t *Tracker) ShutDown(ctx context.Context) {
    close(t.ch)
    select {
    case <-t.stop:
    case <-ctx.Done():
    }
  }
```

- 要避免goroutine泄漏


## Memory model 内存模型

### Go的内存模型就是要说明白变量的赋值谁先谁后的问题

- 读 https://golang.org/ref/mem
  https://www.jianshu.com/p/5e44168f47a3(翻译) 通看这篇文章

- 内存重排 为了提供读写内存的系效率，会对读写指令进行重新排列，这就是所谓的内存重拍，英文为MemoryRecordering
          CPU的设计者会对读写指令进行重新排列，其实还有编译器重排

- 内存屏障 
  CPU提供的锁的机制, atomic compare-and-swap 都用到了这套机制
  对于多线程的程序，所有的 CPU 都会提供“锁” 支持，称之为 barrier，或者 fence。
  它要求：barrier 指令要求所有对内存的操作都必须要“扩散”到 memory 之后才能继续执行其他对 memory 的操作。
  因此，我们可以用高级点的 atomic compare-and-swap，或者 直接用更高级的锁，通常是标准库提供。

cpu l1缓存store buffer计算完了后来不及写入内存的例子
```
 先执行 (1) 和 (3)，将他们直接写入 store buffer， 
      接着执行 (2) 和 (4)。“奇迹”要发生了：
      (2) 看了 下 store buffer，并没有发现有 B 的值，于是从 Memory 读出了 0，
      (4) 同样从 Memory 读出了 0。 最后，打印出了 00
```

- data race 会引起两个问题 
  1.原子  2.可见性
  所以没有安全的data race
```
  查看汇编的命令
  go tool compile -S main.go
  查看代码中是否含有data race的命令
  go build -race main.go
```

## Package sync
- go提供的底层的同步语意的,CAS指令,原子赋值

- 原子赋值
  策略: 在临界区最晚加锁，最早释放, 锁里面的代码越简单越好
- 解决方案
  1. sync.Mutex 大锁，互斥锁
             fast path 
             slow path 会有 goroutine的切换
  
  go 1.8中还存在锁饥饿的情况在之后的runtime版本中做了更新
  ```
    解决锁饥饿的新原理

    首先，goroutine1将获得锁并休眠100ms, 当goroutine2试图获取锁时，它将被添加到锁的
    队列中-FIFO(实际就是park the goroutine),goroutine2将进去等待状态。
    然后，当goroutine1完成它的工作时，它将释放锁，并通知队列唤醒goroutine2,它将被标记
    为可运行的，并且等待Go的runtime来调度它。
  ```

  2. sync.RwMutex  读写锁 它也会有goroutine直接的切换

  3. sync.atomic 最轻量 它跑的是连续的
    atomic.value
    cfg作为包级全局对象，在这个例子中被多个goroutine同时访问

  4. Copy-On-Write(俗称 COW) 思路在微服务降级或者local cache场景中经常使用。
                写时复制指的是，写操作时复制全量老数据到一个新的对象中,
                携带上本次新写的数据，之后利用原子替换atomic.Value,更新调用者的变量。来完成无锁访问共享数据。
                旧版本的数据要等到没有任何人引用它的时候才会被GC掉.

    使用场景: 降级数据, 配置文件, local cache 读多写少的场景，in process cache(进程内缓存)

    使用实例: Redis是怎么实现的bgsafe的(把redis里的key-value dump到磁盘)？
           假如dump的同时有人来查Redis，作为一个单线程的程序redis是不是就无法工作了。Redis的实现非常简单，它就是
           fork了一个进程,一开始fork进程的时候,它的地址空间，也就是它的parents指向的地址空间是一样的。
           老的进程肯定还会有用户不断的往里面去更新数据，那么新的进程是不会收到影响的，因为内存有一个叫做Copy-On-Write
           的功能，哪一个内存页被写了，它会标记为dirty(也就是脏的),就会把它拷出来指向另外一个地方，所以新老相互不影响的。      
 
  5. 判断哪个锁最优的策略，一定要写Benchmark来进行实践
  即便我们知道可能在Mutext vs Atomic的情况里，Mutex相对更重。
  因为涉及到更多的goroutine之间的上下文切换pack blocking goroutine以及
  唤醒goroutine

     使用方法: go test -bench=. config_test.go
  ```
   $go test -bench=. config_test.go
   goos: darwin
   goarch: amd64
   BenchmarkAtomic-8       294080270            4.21 ns/op
   BenchmarkMutext-8       1161136              1331 ns/op
   PASS
   ok      command-line-arguments  3.893s
  ```     

- errgroup see https://pkg.go.dev/golang.org/x/sync/errgroup
  
  1. 定义: 我们把一个复杂的任务，尤其是依赖多个微服务rpc需要聚合数据的任务，分解为依赖和并行。
     依赖的意思是:需要上游a的数据才能访问下游b的数据进行组合.
     并行的意思是:分解微多个小任务并行执行,最终等全部执行完毕.
  2. 原理: 利用sync.Waitgroup管理并执行的 groutine
     - 并行工作流
     - 错误处理或优雅降级
     - context传播和取消
     - 利用局部变量+闭包
  

- sync.Pool 场景是用来保存和服用临时对象，以减少内存分配，降低GC压力(Request-Driven特别合适)
  
## Package context
  1. Request-scoped context
     定义: 在Go的服务器中，每个传入的请求都在其自己的goroutine中处理。
     请求处理程序通常启动额外的goroutine来访问其他后端，如数据库和RPC服务。
     处理请求的goroutine通常需要访问特定于请求(request-specific context)的值。
     例如最终用户的身份、授权令牌和请求的截止日期(deadline)。
     当一个请求被取消或超时时，处理该请求的所有goroutine都应该退出(fail fast), 这样系统就可以
     回收它们正在使用的任何资源
  
  2. 在Go 1.7引入了context包，它使得跨API边界的请求范围元数据、取消信号和截止日期很容易传递给
     处理请求所涉及的所有goroutine(显示传递)
  
  3. 将context集成到API中的要点, context的作用域是请求级别的
     - 方法1: 首参数传递context对象，比如, 参考net包Dialer.DialContext 
       次函数执行正常的Dial操作，但可以通过context对象取消函数调用
       ```
         func (d *Dialer) DialContext(ctx context.Context,network,address string) (Conn,error)
       ```
     - 方法2: 在一个request对象中携带一个可选的context对象.
       例如: net/http库的Request.WithContext,通过携带给定的context对象，返回一个新的Request对象。
       ```
         func (r *Request) WithContext(ctx context.Context) *Request
       ```
  4. context.WithValue
     - 基于valueCtx实现 
       ```
        type valueCtx struct { 
           Context  //Context放在这
	         key, val interface{} //key,value放在这
        } 
        为了实现不断的WithValue，构建新的context, 内部在查找key的时候，使用递归方式不断从当前，从父节点寻
        找匹配的key,直到root context(Background和TODO Value函数会返回nil)
       ```
     - 调用这个方法的使用实际是创建了一个新的,parent是不会改的
     - context.Value的数据是面向请求的原数据，不应该作为函数的可选参数来使用
       (比如context里面挂了一个sqlTx对象,传递到Dao层使用),因为元数据相对函数
       参数更加是隐含的，面向请求的。而参数是更加显示的。
     - 使用场景: Tracking(链路追踪的信息),debug信息,调度的原数据(染色信息,API重要性)
    

### Final Notes
- 使用WithCancel,WIthDeadline,WithTiemout, WIthValue替换一个Context
- 级联取消，派生自parent的Context会级联取消
- 不要在Context携带业务逻辑的数据进去
- 所有耗时很长的操作一定要传递Context便于可以取消
- Context.Value 不应该影响你的业务代码
- 不要用Context.Value做一些业务逻辑的事情
- 不要把业务的代码比如用户id放到context.Value里面

see https://talks.golang.org/2014/gotham-context.slide#1

## channel
- 一定要发送者来close channel

