package middleware

import (
	"github.com/gin-gonic/gin"
	"go-storage/pkg/logger"
	"time"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := logger.WithLogger(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		logger.FromContext(ctx).Info("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
