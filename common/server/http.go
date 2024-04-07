package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func RunHTTPServer(createHandler func(app *fiber.App) error) {
	RunHTTPServerOnAddr(":"+os.Getenv("PORT"), createHandler)
}

func RunHTTPServerOnAddr(addr string, createHandler func(app *fiber.App) error) {
	app := fiber.New()
	setMiddlewares(app)

	err := createHandler(app)
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
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	addCorsMiddleware(app)

	//app.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
}

func addCorsMiddleware(app *fiber.App) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsConfig := cors.ConfigDefault
	corsConfig.AllowOrigins = strings.Join(allowedOrigins, ",")
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = "Origin,Accept,Authorization,Content-Type,X-CSRF-Token"
	corsConfig.AllowMethods = "GET,POST,PUT,DELETE,OPTIONS"
	corsConfig.MaxAge = 300

	app.Use(cors.New(corsConfig))
}
