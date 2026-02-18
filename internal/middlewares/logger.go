package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start)

		logEvent := logger.Info()
		if c.Writer.Status() >= 400 {
			logEvent = logger.Error()
		}

		logEvent.
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("latency", duration).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent())

		if raw != "" {
			logEvent.Str("query", raw)
		}

		if c.Writer.Size() > 0 {
			logEvent.Int("size", c.Writer.Size())
		}

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logEvent.Err(err)
			}
		}

		logEvent.Msg("HTTP request")
	}
}
