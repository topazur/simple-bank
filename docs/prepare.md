# 准备工作

## 第一节 · 数据库设计

### 编辑区
> [使用dbdiagram.io设计SQL数据库结构](https://dbdiagram.io
) ➡ 可导出sql、pdf、png

```python
//// -- LEVEL 1
//// -- Tables and References

// as 别名
Table accounts as A {
  // bigserial 大的自动递增的整数（8字节/64位）；pk 主键
  // increment: auto-increment
  id bigserial [pk]
  owner varchar [not null]
  // 余额
  balance bigint [not null]
  // 货币名称
  currency Currency [not null]
  // timestamp 时间戳；timestamptz 带时区
  created_at timestamptz [not null, default: `now()`]
  
  // 索引列表：query
  Indexes {
    owner
  }
}

// 记录账户余额的所有更改
Table entries {
  id bigserial [pk]
  account_id bigint [not null, ref: > A.id]
  amount bigint [not null, note: '收入or支出']
  created_at timestamptz [not null, default: `now()`]
  
   Indexes {
    account_id
  }
 }

// 转账：记录2个账户之间的所有汇款
Table transfers {
  id bigserial [pk]
  // 外键引用
  from_account_id bigint [not null, ref: > A.id]
  to_account_id bigint [not null, ref: > A.id]
  amount bigint [not null, note: '支出']
  created_at timestamptz [not null, default: `now()`]
  
   Indexes {
    from_account_id
    to_account_id
    (from_account_id, to_account_id)
  }
 }

// 货币枚举
// 或者使用内置类型varchar，并在程序内处理值验证
Enum Currency {
  USD
  EUR
}

// 声明外键引用的方法
// Creating references
// You can also define relaionship separately
// > many-to-one; < one-to-many; - one-to-one
// Ref: U.country_code > countries.code  
// Ref: merchants.country_code > countries.code
// Ref: order_items.product_id > products.id
```


## 第二节 · 环境配置 (数据库服务，建表，GUI)

### docker
> [下载地址和hub镜像](https://hub.docker.com/editions/community/docker-ce-desktop-windows)

> Docker Desktop 是 Docker 在 Windows 10 和 macOS 操作系统上的官方安装方式，这个方法依然属于先在虚拟机中安装 Linux 然后再安装 Docker 的方法。

> 由于开发环境变动，或不易配置，我使用云服务器搭建环境。防止关键信息泄露，通过环境变量方式导入字段。

### PostgresSQL
> [镜像地址](https://hub.docker.com/_/postgres)

> TablePlus ➡ mac电脑GUI工具

#### 镜像相关操作
- 拉取镜像：`docker pull postgres:13-alpine`
- 创建容器：`docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine` 通过环境变量设置用户名(user默认为postgres)密码。
- sql进入容器：`docker exec -it postgres13 psql -U root` 以psql命令方式进入；Postgres容器在本地设置信任身份验证，无需密码;运行`select now()`命令测试；`\q`退出sql。
- shell进入容器：`docker exec -it postgres13 /bin/sh` 以shell命令方式进入；Postgres容器在本地设置信任身份验证，无需密码;运行`createdb --username=root --owner=root simple_bank
`命令创建数据库（user用户身份连接，owner创建数据库归属）；`psql simple_bank`psql访问数据库；`dropdb simple_bank`删除数据库；`exit`退出shell。
- 不进入容器直接运行容器内的命令：`docker exec -it postgres13 createdb --username=root --owner=root simple_bank`；`docker exec -it postgres13 psql -U root simple_bank`


#### windows安装错误情况
[端口占用、data下存在pid、data文件读写权限](https://blog.csdn.net/international24/article/details/89710703) | 
[pg_ctl没有指定数据目录并且没有设置PGDATA环境变量](https://blog.csdn.net/IToBeNo_1/article/details/79808817) 
| [开放防火墙5432端口](https://blog.csdn.net/weixin_40598838/article/details/111875617) 
| [pg_ctl 命令](https://www.cnblogs.com/hello-wei/p/10150883.html) 
| [在Windows平台上安装与运行PostgreSQL的常见问题与解答](https://wiki.postgresql.org/wiki/%E5%9C%A8Windows%E5%B9%B3%E5%8F%B0%E4%B8%8A%E5%AE%89%E8%A3%85%E4%B8%8E%E8%BF%90%E8%A1%8CPostgreSQL%E7%9A%84%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98%E4%B8%8E%E8%A7%A3%E7%AD%94) 
| 初始化服务`pg_ctl init -D "C:\Program Files\PostgreSQL\13\data"`、
[postgresql windows 服务启动失败](https://www.cnblogs.com/wolbo/p/11551686.html)、
[启动服务](https://blog.csdn.net/weixin_32227927/article/details/113477719)


## 第三节 · 数据库模式迁移(golang-migrate)

### golang-migrate安装
[migrate](https://github.com/golang-migrate/migrate) | 
[migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) | 
<a style="color:red" href="https://github.com/golang-migrate/migrate/releases">windows将安装包路径配置Path后命令执行即可：Release Downloads</a>

### 使用示例（加入Makefile）
```bash
$ mkdir -p db/migrate

# 文件扩展名为sql；存放目录是‘db/migration’；`-seq`标志来生成迁移文件的**连续**版本号；`init_schema`是自定义的迁移名称
# 生成前缀相同而后缀分别为up和down的两个文件：运行up-script（将按其前缀版本的相同顺序依次运行）以对模式进行向前更改；想恢复up-script所做的更改，则运行down-script（将按其前缀版本的相反顺序依次运行）。
$ migrate create -ext sql -dir db/migration -seq init_schema

# path指定迁移文件dir；database指定数据库服务器的URL(postgres容器默认不启用SSL)；verbose要求migrate打印详细日志记录；最后使用“up”参数告诉migrate运行`migrate up`
$ migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
# 生成规定的表 以及 schema_migrations表存储最新应用的迁移版本；dirty代表当前版本是干净的还是脏的。f或0表示干净（即没有出问题）。

$ migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
```

### Makefile

[windows通过mingw安装mingw32-make包使用make](https://www.cnblogs.com/TonyJia/p/13212110.html)、
[Make 命令教程 - 阮一峰的网络日志](http://www.ruanyifeng.com/blog/2015/02/make.html)、
[报错：`Makefile:13: ***  。 停止`](http://blog.sciencenet.cn/blog-1470666-873932.html)

