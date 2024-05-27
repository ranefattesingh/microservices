package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ranefattesingh/pkg/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)

		status := c.Writer.Status()

		log.Logger().Info(
			c.Request.URL.Path,
			zap.Duration("latency", latency),
			zap.Int("http_status", status),
			zap.String("method", c.Request.Method),
		)
	}
}
