## 作业
- 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API注册，以及 main 函数对于服务的注册和启动，
信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

## 准备工作
### 安装protoc
```shell
    go get -u google.golang.org/protobuf/cmd/protoc-gen-go
              google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### 创建pb
``` 
    cd ./api/user/
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto;
```

### 安装wire
```
    go get -u github.com/google/wire/cmd/wire
```

## TODO list
- <input type="checkbox" checked> API层，通过proto文件生成grpc服务接口  
- <input type="checkbox"> 使用wire构建依赖  
- <input type="checkbox"> 在main中实现服务的注册和启动
- <input type="checkbox"> 服务信号处理
- <input type="checkbox"> 分层架构


## 业务分层

### service层职责
1. 接口层，实现grpc server接口。
2. 应用层，同时实现biz业务逻辑、其他微服务接口的组装，DTO到领域对象的转换。
3. 消息队列消费

### biz层职责
业务逻辑，各业务逻辑应互相独立。如果有依赖放到service层去组装。
repo接口定义。

### data层职责
实现repo接口，从db获取数据并实现缓存控制。

## 实现
api实现一个profile服务，提供一个查询用户profile的接口。返回用户昵称和用户的金币数。
biz实现一个user业务提供用户基础数据，一个coin业务提供用户的金币数据。
data实现user、coin表的读取与缓存

