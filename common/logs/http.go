package logs

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// NewStructuredLogger 创建一个结构化的 Fiber 日志中间件
func NewStructuredLogger(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		fields := logrus.Fields{
			"remote_addr": c.IP(),
			"method":      c.Method(),
			"uri":         c.OriginalURL(),
		}

		if reqID := c.Get("X-Request-Id"); reqID != "" {
			fields["req_id"] = reqID
		}

		logger := logger.WithFields(fields)
		logger.Info("Request started")

		err := c.Next()

		res := c.Response()
		fields = logrus.Fields{
			"resp_status":       res.StatusCode(),
			"resp_bytes_length": len(res.Body()),
			"resp_elapsed":      time.Since(start).Round(time.Millisecond / 100).String(),
		}

		if err != nil {
			logger.WithFields(fields).Error(err)
		} else {
			logger.WithFields(fields).Info("Request completed")
		}

		return err
	}
}

// GetLogEntry 获取当前请求的日志记录器
func GetLogEntry(c *fiber.Ctx) *logrus.Entry {
	fields := logrus.Fields{
		"remote_addr": c.IP(),
		"method":      c.Method(),
		"uri":         c.OriginalURL(),
	}

	if reqID := c.Get("X-Request-Id"); reqID != "" {
		fields["req_id"] = reqID
	}

	return logrus.StandardLogger().WithFields(fields)
}

// PanicHandler 处理 Fiber 中的 panic
func PanicHandler(c *fiber.Ctx, err interface{}) {
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)

	fields := logrus.Fields{
		"stack": string(stack[:length]),
		"panic": fmt.Sprintf("%+v", err),
	}

	GetLogEntry(c).WithFields(fields).Error("Panic occurred")
}
