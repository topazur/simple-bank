## 第13节 · Mock DB

#### 1. 优点分析

- 更轻松地编写独立的测试：因为每个测试都会使用自己单独的模拟数据库来存储数据，所以他们之间不会有冲突。
- 测试将运行得更快：因为他们不必花时间与数据库交谈并等待查询运行。所有操作都将在内存中执行，并在同一进程中执行。
- 允许我们编写实现 100% 覆盖率的测试：使用模拟数据库，我们可以轻松设置和测试一些边缘情况。例如意外错误或连接丢失，如果我们使用真正的数据库，这是不可能实现的。

#### 2. How to mock?

##### 假数据库
自己实现假的数据库切片，通过映射来读取数据

##### github.com/golang/mock

- mockgen 命令初始化

```bash
  -destination string
        Output file; defaults to stdout.  默认情况下，mockgen 会将生成的代码写入 stdout。
```

###### 结构体

1. MockStore - 是实现 Store 接口所有必需功能的结构体。

2. MockStoreMockRecorder - 是 MockStore 的模拟记录器，会生成与MockStore相同的函数用来记录


> 生成gomock.Controller，构建NewMockStore，通过http.NewRequest和httptest.NewRecorder发起请求并记录响应

> 测试api即可以仅测试单个情况，也可以并发测试所用情况；测试时数据库操作是调用store中的方法，但是返回结果的情况是根据自己的api接口情况得到的

```
·测试开发编写流程：
api响应情况 --> 产生需要测试的情况 ---> 根据每个情况调用store方法操作数据库；对响应code和body使用testify校验
```
