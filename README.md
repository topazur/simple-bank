# go-simple-bank

[学习来源 simplebank](https://github.com/techschool/simplebank.git)

## Angular git commit message 规范
<details>
<summary>点击查看 <code>Angular git commit</code></summary>
<code>

```bash
# 注意冒号后面的空格
$ feat: 新功能（feature）
$ fix: 修补bug
$ docs: 文档（documentation）
$ style: 格式（不影响代码运行的变动）
$ refactor: 重构（即不是新增功能，也不是修改bug的代码变动）
$ chore: 构建过程或辅助工具的变动
$ revert: 撤销，版本回退
$ perf: 性能优化
$ test：测试
$ improvement: 改进
$ build: 打包
$ ci: 持续集成
```
</code>
</details>

<br />

## Go Modules
[语雀go modules笔记](https://www.yuque.com/topazur/golang/axng7n)
<details>
<summary>点击查看 <code>Go Modules</code></summary>
<code>

```sh
$ go env -w GOPROXY=https://goproxy.cn,direct
# 开启go modules, 默认值为auto
$ go env -w GO111MODULE=on
$ go mod init [模块名,或者github路径]  # 初始化当前文件夹, 创建go.mod文件
$ go mod download    # 下载依赖的module到本地cache（默认为$GOPATH/pkg/mod目录）
$ go mod edit -fmt   # 编辑go.mod文件(选项有-json、-require和-exclude)
$ go mod graph       # 以文本模式打印模块依赖图
$ go mod tidy        # 增加缺少的module，删除无用的module
$ go mod vendor      # 在根目录生成vendor目录,将依赖复制到vendor下
$ go mod verify      # 校验依赖是否正确
$ go mod why         # 查找依赖
# 总结: 获取第三方依赖时(下载到GOPATH/pkg中)
$ go get -u -v package # v显示下载详细信息,u将会升级到最新的次要版本或者修订版本(x.y.z, z是修订版本号， y是次要版本号)
```
</code>
</details>

<details>
<summary>点击查看 <code>Go Command</code></summary>
<code>

```sh
# 像执行脚本文件一样执行Go代码
$ go run
# go install表示安装的意思:
# 先编译源代码得到可执行文件,然后将生成的可执行文件移动到GOPATH的bin目录下。
# 因为我们的环境变量中配置了GOPATH下的bin目录，所以我们就可以在任意地方直接访问/执行可执行文件了。
$ go install
# 将源代码编译成可执行文件。生成的执行文件都在当前执行命令的目录下。
$ go build
```
</code>
</details>
