package logs

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// NewStructuredLogger 创建一个结构化的 Fiber 日志中间件 TODO 需要优化
func NewStructuredLogger(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		fields := logrus.Fields{
			"ip":     c.IP(),
			"method": c.Method(),
			"uri":    c.OriginalURL(),
		}
		if reqID := c.Locals("requestid"); reqID != "" {
			fields["req_id"] = reqID
		}

		logger.WithFields(fields).Info("Request started")

		errHandler := c.App().ErrorHandler

		chainErr := c.Next()

		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				// 使用自定义ErrorHandler时，这段代码是不可达的，因为它的返回值是一个固定的nil
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		res := c.Response()
		fields = logrus.Fields{
			"status":     res.StatusCode(),
			"bytes_sent": len(res.Body()),
			"latency":    time.Since(start).Round(time.Millisecond / 100).String(),
		}

		if chainErr != nil {
			logger.WithFields(fields).Error(chainErr)
		} else {
			logger.WithFields(fields).Info("Request completed")
		}

		return nil
	}
}

//// GetLogEntry 获取当前请求的日志记录器
//func GetLogEntry(c *fiber.Ctx) *logrus.Entry {
//	fields := logrus.Fields{
//		"ip":     c.IP(),
//		"method": c.Method(),
//		"uri":    c.OriginalURL(),
//	}
//
//	if reqID := c.Get("requestid"); reqID != "" {
//		fields["req_id"] = reqID
//	}
//
//	return logrus.StandardLogger().WithFields(fields)
//}
//
//// PanicHandler 处理 Fiber 中的 panic
//func PanicHandler(c *fiber.Ctx, err interface{}) {
//	stack := make([]byte, 4096)
//	length := runtime.Stack(stack, false)
//
//	fields := logrus.Fields{
//		"stack": string(stack[:length]),
//		"panic": fmt.Sprintf("%+v", err),
//	}
//
//	GetLogEntry(c).WithFields(fields).Error("Panic occurred")
//}
