## test

#### main_test.go
在TestMain内连接数据库

#### sql_test.go
- 要检查测试结果，我建议使用 testify 包 `go get github.com/stretchr/testify`
- 仅使用子包require即可完成测试需求
- 不满足测试条件将自动使测试失败。
```
NotZero() 函数将断言一个值不能是其类型的零值。
```

#### db transaction 数据库事务

- 是一个单一的工作单元
- 通常由多个数据库操作组成。

##### 好处
- 即使在系统故障的情况下，也能保证工作单元可靠且一致
- 在程序之间提供隔离，并发访问数据库

##### ACID Property
> 数据库事务必须满足 ACID 属性。
- Atomicity 原子性：意味着交易的所有操作都成功完成，或者整个交易失败，一切都回滚，并且数据库不变。
- Consistency 一致性：意味着数据库状态在事务执行后应该保持有效，更准确地说，所有写入数据库的数据必须是有效的。（根据预定义的规则，包括约束、级联和触发器。）
- Isolation 隔离：意味着并发运行的所有事务不应相互影响。(隔离级别,个事务所做的更改何时可以对其他人可见。)
- Durability 耐久性：它基本上意味着成功事务写入的所有数据必须留在持久存储中，并且不会丢失，即使在系统故障的情况下。

##### How to run SQL TX
我们用 BEGIN 语句开始一个事务。然后我们编写一系列正常的 SQL 查询（或操作）。
如果他们都成功了，我们提交交易以使其永久化，数据库将更改为新状态。
否则，如果任何查询失败，我们回滚交易，因此之前交易查询所做的所有更改都将消失，并且数据库保持与事务之前相同。

##### SQL TX in Golang



# MiaoSha-Yiletian
拉勾专栏《打造千万级流量秒杀系统》源码


https://github.com/lagoueduCol/MiaoSha-Yiletian.git

https://github.com/letian0805/seckill


### build
> compile packages and dependencies(编译文件得到二进制可执行文件)

- 跨平台编译：`env GOOS=linux GOARCH=amd64 go build`

### install
> compile and install packages and dependencies

- 也是编译，与build最大的区别是编译后会将输出文件打包成库放在pkg下
- 常用于本地打包编译的命令: go install

### get
> download and install packages and dependencies
- 用于获取go的第三方包，通常会默认从git repo上pull最新的版本
- 常用命令(u表示最新版本)如: `go get -u github.com/go-sql-driver/mysql` (从github上获取mysql的driver并安装至本地)

### fmt
> go fmt (reformat) package sources
- 类似于C中的lint，统一代码风格和排版;不只有检测，还有修复功能
- 常用命令如:go fmt

### test
> test packages
- 运行当前包目录下的tests
- 常用命令(v会显示详细信息)如: `go test` 或 `go test -v` 等
- Go的test一般以 `XXX_test.go` 为文件名; package名相同即可
- XXX的部分一般为XXX_test.go所要测试的代码文件名。注:Go并没有特别要求XXX的部分必须是要测试的文件名。目的仅仅是相呼应

- test文件下的每一个test case函数均必须以Test开头并且符合 `TestXxx` 形式，否则go test会直接跳过测试不执行 （一般不满足此格式的函数作为被调用函数）;test文件可以含多个test case函数
- test case的入参为 `t *testing.T` 或 `b *testing.B`(benchmark测试性能)
- t.Errorf为打印错误信息，并且当前test case会被跳过

- `t.SkipNow()` 为跳过当前test，并且直接按PASS处理继续下一个test; 必须放在test case第一行才生效
- Go的test不会保证多个TestXxx是顺序执行，但是通常会按顺序执行
- 使用`t.Run(key string, TestCase func )`来执行 subtests 可以做到控制test输出以及test的顺序

- 使用TestMain（参数为`m *testing.M`）作为初始化test，并且使用m.Run()来调用其他tests之前，可以完成一些需要初始化操作的testing，比如数据库连接，文件打开，REST服务登录等
- 如果没有在TestMain中调用m.Run()则除了TestMain以外的其他tests都不会被执行

### Test之benchmark
- benchmark函数—般以Benchmark开头
- benchmark的case—般会跑b.N次，而且每次执行都会如此
- 在执行过程中会根据实际case的执行时间是否稳定会增加b.N的次数以达到稳态才会停止
- benchmark同样是test case，也受TestMain的控制
- `go test -bench=.`命令只会执行带benchmark的test case

```go
func aaa(n int) int {
  for n > 0 {
    n--
  }
  return n
}

func BenchmarkAll(b *testing.B) {
	// b.N 变化的;每次都会增加，但是aaa随着n的增加执行时间也在增加，并不会区域稳态，导致test不会停止
	for n := 0; n < b.N; n++ {
		aaa()
	}
}
```

## 热启动 `rizla main.go`
`go get github.com/kataras/rizla`