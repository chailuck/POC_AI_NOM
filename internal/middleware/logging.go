// internal/middleware/logging.go
package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLogging(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			// Log request details
			logger.Infow("Request started",
				"method", req.Method,
				"uri", req.RequestURI,
				"remote_addr", req.RemoteAddr,
			)

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// Log response details
			logger.Infow("Request completed",
				"method", req.Method,
				"uri", req.RequestURI,
				"status", res.Status,
				"duration", time.Since(start),
			)

			return err
		}
	}
}
