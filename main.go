package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/topaz-h/go-simple-bank/api"
	db "github.com/topaz-h/go-simple-bank/db/sqlc"
	"github.com/topaz-h/go-simple-bank/util"
)

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	/// 传入数据库driver和连接地址 连接数据库
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// defer conn.Close()

	/// 通过*sql.DB实现store&server
	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
