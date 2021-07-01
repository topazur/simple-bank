package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	// 文件夹名与package名不同：需显示指定package名
	db "github.com/topaz-h/go-simple-bank/db/sqlc"
)

// Server 结构体提供实现 HTTP API 服务器的方法.
type Server struct {
	// 它将允许我们在处理来自客户端的 API 请求时与数据库进行交互。
	store *db.Store
	// 该路由器将帮助我们将每个 API请求发送到正确的处理程序进行处理
	router *gin.Engine
}

/// NewServer 创建一个新的 HTTP 服务器并设置路由.
func NewServer(store *db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}
	// 默认使用了2个中间件Logger(), Recovery()
	router := gin.Default()

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server, nil
}

/// Start 在指定地址上运行 HTTP 服务器.
func (server *Server) Start(address string) error {
	// TODO: 添加一些优雅的关闭逻辑...
	return server.router.Run(address)
}

/// gin.H 本质上是 map[string]interface 的map
/// 统一处理异常时api返回的响应体
func errorResponse(err error) gin.H {
	fmt.Println(err.Error())
	return gin.H{"error": err.Error()}
}
