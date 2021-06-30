package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/topaz-h/go-simple-bank/util"
)

/*
Error: 无法连接到数据库：未知驱动程序“mysql”。

[Go Mysql驱动](https://www.cnblogs.com/vincenshen/p/9427798.html)
`go get github.com/go-sql-driver/mysql`

Go本身（"database/sql"）不提供具体数据库驱动，只提供驱动接口和管理。
为了与特定的数据库引擎交谈。需要配合数据库驱动使用。各个数据库驱动需要第三方实现，并且注册到Go中的驱动管理中。

因为我们实际上并没有直接在代码中调用 驱动 的任何函数，仅通过引入使用其init方法加载引擎，必须通过空白导入告诉 go formatter 保留它。
运行`go mod tidy`来清理依赖项; 未导入显示的`indirect`会被去掉
*/
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	/// 传入数据库driver和连接地址 连接数据库
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer testDB.Close()

	/// 使用 New & conn 来创建新的 Queries 对象。
	testQueries = New(testDB)

	/* m.Run() 之前的准备 ↑↑↑ */
	os.Exit(m.Run())
	/// start test ...
	// 此函数将返回一个退出代码，它告诉我们测试是通过还是失败。
	// 然后我们应该通过 os.Exit 命令将它报告回测试运行器。
}
