package middleware

import (
	"github.com/gin-gonic/gin"
	"go-storage/internal/config"
)

func Config(cnf config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := config.WithConfig(c.Request.Context(), cnf)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
