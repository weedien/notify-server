package main

import (
	_ "github.com/go-sql-driver/mysql" // MySQL驱动
	"github.com/gofiber/fiber/v2"
	"github.com/weedien/notify-server/common/logs"
	"github.com/weedien/notify-server/common/server"
	"github.com/weedien/notify-server/template/ports"
	"github.com/weedien/notify-server/template/service"
)

func main() {
	logs.Init()

	application := service.NewTemplateApp()

	server.RunHTTPServer(func(app *fiber.Router) error {
		ports.RegisterHandlers(*app, ports.NewHttpServer(application))
		return nil
	})
}
