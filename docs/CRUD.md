## 第四节 · CRUD

#### 4.1 工具分析

##### 4.1.1 [基础的标准库database/sql包](http://golang.org/pkg/database/sql)
  - 只使用 QueryRowContext() 函数，并传入原始 SQL 查询和一些参数，然后我们将结果扫描到目标变量中。
  - 运行快，编写代码简单,但代码较长
  - 手动将SQL字段映射到变量(易出错，且错误只会在运行时出现)
##### 4.1.2 [gorm - 用于golang的高级对象关系映射库。](http://gorm.io/docs)
  - 强大，便捷，所有的CRUD操作都已经实现，只需要声明模型，并调用 gorm 提供的函数。
  - 需要学习对应CRUD函数使用方法，了解不够使用起来不够灵活；必须学习如何声明关联标签，让gorm理解表之间的关系，；当流量很高时速度较慢（比标准库慢3~5倍）
##### 4.1.3 sqlx
  - 速度与标准库相差无几
  - 字段映射是通过查询文本或结构标签完成的,不必手动映射
  - 代码较长，且错误只会在运行时被捕获
##### 4.1.4 [sqlc](www.sqlc.dev)
  - 速度与标准库相差无几
  - 只需将 db schema 和 SQL 查询传递给 sqlc，每个查询上面都有 1 条注释告诉 sqlc 生成正确的函数签名， 然后sqlc会生成惯用的Golang代码。
  - sql语句通过database/sql去执行；任何错误都会在generate阶段被发现并立即报告；只支持Postgres,mysql处于beta阶段(mysql现已支持，SQLite is planned)
  - 都是生成的文件，我们无需修改他


### 4.2 sqlc当选

#### 4.2.1 sqlc安装
<a style="color:red" href="https://github.com/golang-migrate/migrate/releases">windows将安装包路径配置Path后命令执行即可：Release Downloads</a>
**windows下初始化会失败**

```bash
$ sqlc version

# 创建一个空的 sqlc.yaml 设置文件。
$ sqlc init #根目录生成sqlc.yaml

# 会做两个错误检查，并为我们从 SQL 查询生成 golang 代码。
$ sqlc generate #根据migration下的schema和query SQL生成golang代码
```

#### 4.2.2 sqlc使用

[Go 每日一库之 sqlc](https://blog.csdn.net/darjun/article/details/106866664)
```
-- name: <name> <cmd>
name为生成的方法名，如上面的CreateAuthor/ListAuthors/GetAuthor/DeleteAuthor等，cmd可以有以下取值：

:one：表示 SQL 语句返回一个对象，生成的方法的返回值为(对象类型, error)，对象类型可以从表名得出；
:many：表示 SQL 语句会返回多个对象，生成的方法的返回值为([]对象类型, error)；
:exec：表示 SQL 语句不返回对象，只返回一个error；
:execrows：表示 SQL 语句需要返回受影响的行数。
```

#### 4.2.3 生成类型安全的、地道的 Go 接口代码

##### models.go
schema结构体含json名称 && commit注释

##### db.go
该文件包含 dbtx 接口。
它定义了 sql.DB 和 sql.Tx 对象都有的 4 个常用方法。
这允许我们自由地使用数据库或事务来执行查询。

New() 函数将 DBTX 作为输入并返回一个 Queries 对象。

所以我们可以传入一个 sql.DB 或 sql.Tx 对象,取决于我们是否只想执行 1 个查询、或一组事务中的多个查询。

WithTx允许查询实例与事务相关联。

##### account.sql.go
Sqlc 已明确地将 RETURN * 替换为所有列的名称。查询清晰并避免以不正确的顺序扫描值。
SELECT * 也会自动替换成相应字段。