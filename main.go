package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
	"github.com/weedien/countdown-server/server"
	"github.com/weedien/countdown-server/store"
)

func main() {
	// 初始化数据库
	store.InitDB()

	// 初始化Snowflake节点
	store.InitSnowFlakeNode()

	// 初始化路由
	server.InitRoutes()

	// 从命令行中获取端口信息
	port := flag.String("port", "8080", "HTTP服务器端口")

	// 启动HTTP服务器
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
