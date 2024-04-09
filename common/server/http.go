package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/config"
	"github.com/weedien/notify-server/common/logs"
	"github.com/weedien/notify-server/common/result"
	"strconv"
	"strings"
)

func RunHTTPServer(createHandler func(router *fiber.Router) error) {
	port := config.Config().Server.Port
	RunHTTPServerOnAddr(":"+strconv.Itoa(port), createHandler)
}

func RunHTTPServerOnAddr(addr string, createHandler func(router *fiber.Router) error) {
	app := fiber.New(customHttpServerConfig())
	setMiddlewares(app)

	// 统一前缀 /api
	api := app.Group("/api")

	err := createHandler(&api)
	if err != nil {
		logrus.WithError(err).Panic("Unable to create handler")
	}

	logrus.Info("Starting HTTP server")
	err = app.Listen(addr)
	if err != nil {
		logrus.WithError(err).Panic("Unable to start HTTP server")
	}
}

func setMiddlewares(app *fiber.App) {
	app.Use(requestid.New())

	env := config.Config().Server.Env
	if env == "prod" {
		// 生产环境使用logrus日志便于统一管理
		app.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	} else {
		app.Use(logger.New())
	}

	app.Use(recover.New())
	addCorsMiddleware(app)
}

func addCorsMiddleware(app *fiber.App) {
	allowedOrigins := strings.Split(config.Config().Security.CorsAllowedOrigins, ",")
	if len(allowedOrigins) == 0 {
		return
	}

	corsConfig := cors.ConfigDefault
	corsConfig.AllowOrigins = strings.Join(allowedOrigins, ",")
	corsConfig.AllowHeaders = "Origin,Accept,Authorization,Content-Type,X-CSRF-Token"
	corsConfig.AllowMethods = "GET,POST,PUT,DELETE,OPTIONS"
	corsConfig.MaxAge = 300

	app.Use(cors.New(corsConfig))
}

func customHttpServerConfig() fiber.Config {
	return fiber.Config{
		// 自定义错误处理 支持返回json
		// 主要是为了使用fiber.NewError进行错误处理，便于fiber日志进行记录
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = ctx.Status(code).JSON(result.Fail(code, err.Error()))
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(result.Fail(code, "Internal Server Error"))
			}
			return nil
		},
	}
}
