## 作业
- 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API注册，以及 main 函数对于服务的注册和启动，
信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。
Week04作业提交地址：https://github.com/Go-000/Go-000/issues/76

# Week04学习笔记

## 目录结构设计

### cmd
- 每个应用程序的目录名应该与你想要的可执行文件的名称相匹配
  (例如: /cmd/myapp/main.go,这样编译出来的应用程序名称就是myapp了)

### internal
- 私有应用程序和库代码
```
 这个布局模式是由 Go 编译器本身执行的
 并不局限于 顶级 internal 目录。在项目树的任何级别上都可以有多个内部目录。 
 你可以选择向 internal 包中添加一些额外的结构，以分隔共享和非 共享的内部代码。
 
 你的实际应用程序代码可以放在 /internal/app 目录下(例如 /internal/app/myapp)
 
 这些 应用程序共享的代码可以放在 /internal/pkg 目录下(例如 /internal/pkg/myprivlib)。
 因为我们习惯把相关的服务，比如账号服务，内部有 rpc、job、 admin 等，相关的服务整合一起后，
 需要区分 app。单一的服务， 可以去掉 /internal/myapp。
``` 
#### internal子目录 v1
- model 放各种结构体,映射了mysql表里的一个结构体
- dao (database access object 数据访问层) 
  import model层的结构体
  访问数据库的的一些方法
  通常是一个文件里就对应着一张数据库里的一张表或是mysql中的键值对儿或集合
- service 依赖dao层的接口
  上层是service,下层是dao,依赖下层的时候不能依赖细节，要依赖抽象, 这就是依赖倒置了

定义：高层模块儿不能依赖低层模块儿，而是大家都依赖于抽象
  抽象不能依赖实现，而是实现依赖抽象
  依赖倒置其实就是高层来定义抽象，底层来实现抽象，这也就是依赖倒置了
- server grpc的监听启动代码

#### internal子目录 v2
- biz: 业务逻辑的组装层,类似 DDD 的 domain 层,
       data类DDD的 repo, repo接口在这里定义,
       使用依赖倒置的原则。

- data: 业务数据访问，包含 cache、db 等封装，  
        实现了 biz 的 repo 接口. 我们可能会把 data 
        与 dao 混淆在一起,data偏重业务的含义,
        它所要做的是将领域对象重新拿出来。
        我们去掉了DDD 的 infra(基础设施)层

- service: 实现了 api 定义的服务层，类似 DDD 的 
        application 层，处 理 DTO 到 biz 领域实体的转换(DTO -> DO),
        同时协同各类 biz 交互,但是不应处理复杂逻辑。

- PO(持久化对象) 它跟持久层（通常是关系型数据库）的数据结构
               形成一一对应的映射关系，如果持久层是关系型数据库，
               那么数据表中的每个字段（或若干个）就对应 PO 的一
               个（或若干个）属性

### pkg
- 外部应用程序可以使用的库代码(例如 /pkg/mypubliclib)
```
   /pkg 目录内，可以参考 go 标准库的组织方式，按照功能分类。 
   /internla/pkg 一般用于项目内的 跨多个应用的公共共享代码，
   但其作用域仅在单个项目工程内
```

### kit 
- 为不同的微服务建立一个统一的kit工具包项目(基础库/框架) 和 app 项目

- 统一、标准库方式布局、高度抽象、支持插件

### api
- API 协议定义目录，xxapi.proto protobuf 文件，以及生成的 go 文件。
  通常把 api 文档直接在 proto 文件中描述。

### configs 配置文件模板或默认配置

### test
```
  额外的外部测试应用程序和测试数据。你可以随时根据需求构造 /test 目录。
  对于较大的项目，有一个数据子目录是有意义的。
  例如，你可以使用 /test/data 或 /test/testdata (如果你需要忽略目录中的内容)。
  请注意，Go 还会忽略以“.”或“_”开头的目录或文件，
  因此在如何命名测试数据目录方面有更大的灵活性。
```

## Service Applicaton Project

### app目录内每个微服务的命名按照全局唯一的名称
- 例如:"account.service.vip" 三段式来命名

### 微服务中的app服务类型分为4类
    interface,service,job,admin
1. interface: 对外的 BFF 服务，接受来自用户的请求， 比如暴露了 HTTP/gRPC 接口
2. service: 对内的微服务，仅接受来自内部其他服务或 者网关的请求，比如暴露了gRPC 接口只对内服务
3. admin：区别于 service，更多是面向运营测的服务， 通常数据权限更高，隔离带来更好的代码级别安全 
4. job: 流式任务处理的服务，上游一般依赖 message broker
5. task: 定时任务，类似 cronjob，部署到 task 托管平 台中

- cmd应用目录负责程序的 启动、关闭、配置初始化


