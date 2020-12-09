## 作业
- 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么？应该怎么做请写出代码

#### 我的理解
- 应该是使用Warp把错误向上抛,然后在面向横向切面功能的位置来记录这些错误。但在某些需要降级的场景下，适合在server层吐掉错误。

#### 实验结果为
```
original error: *errors.errorString record not found
stack trace:
record not found
find user by userid: 100 
week02/dao.(*User).GetByUserID
	/Users/lanecao/Documents/GoAdvanced/Go-000/Week02/dao/dao.go:24
week02/service.(*UserService).GetByUserID
	/Users/lanecao/Documents/GoAdvanced/Go-000/Week02/service/service.go:9
main.main
	/Users/lanecao/Documents/GoAdvanced/Go-000/Week02/main.go:11
runtime.main
	/usr/local/Cellar/go/1.15.3/libexec/src/runtime/proc.go:204
runtime.goexit
	/usr/local/Cellar/go/1.15.3/libexec/src/runtime/asm_amd64.s:1374
User={0 }
```

### Error Type

- Sentinel Error

  **预定义的特定错误**

  调用方需要使用 == 将结果与预先声明的值进行等值比较
  
  在两个包之间建立了源代码依赖关系

  sentinel错误处理策略会增大包的表面积，应该尽可能避免sentinel errors

- Error Types

  **实现了Error接口的自定义类型**

  可以包装底层的错误以提供更多上下文

  调用者要使用类型断言和类型switch，就要让自定义的error 变为public，这会导致和调用者产生强耦合

  应该避免使用，或者至少避免它们作为公共API的一部分

- Opaque Errors

  **不透明错误处理**

  只需返回错误而不假设其内容

  作为调用者，关于操作的结果所知道的就是成功或者失败，没有能力看到错误的内部

  可以断言错误实现了特定的行为，而不是断言错误是特定的类型或值

  它要求代码和调用者之间的耦合最少，是最灵活的错误处理策略


### Handling Error

#### Wrap errors

- 使用github.com/pkg/errors包，可以向错误值添加上下文

- 底层的方法应该是warp往上抛，统一在面向切面横向编程的插件中记录错误就好了 

示例代码
```
  func readfile(path string)([]byte,error) {
     f,err := os.Open(path)
     if err!= nil {
         return nil,errors.Wrap(err,"open failed")
     }
     defer f.close()
  }

  func main() {
      _,err := ReadConfig()
      if err!= nil {
          fmt.Printf("original error: %T %v\n", errors.Cause(err),errors.Cause(err))
          fmt.Printf("stack trace:\n%+v\n",err)
      }
  }
```

使用Wrap的场景
- 基础库(Kit库)不适合使用Warp，它只能返回跟因
- 如果和其他库进行协作或者是dao层的方法，应使用errors.Wrap或者errors.Wrapf保存堆栈信息
- 一旦处理了错误,那么错误就不再是错误，比如降级就应该返回降级数据 
- 在程序的顶部或者是工作的goroutine 顶部(请求入口)，使用%+v把堆栈详情记录

#### Is & As

go1.13 errors 包中包含两个用于检查错误的新函数：Is 和 As

errors.Is函数的行为类似于(sentinel error)的判等操作

errors.As函数的行为类似于类型断言(type assertion)

在处理包装错误(包含其他错误的错误)时，这些函数会考虑错误链中的所有错误