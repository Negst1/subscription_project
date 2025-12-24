package middleware

import (
	"time"

	"subscribe_project/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestLogger := logger.Log.WithFields(logrus.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
			"request_id": c.Get("X-Request-ID"),
			"start_time": start.Format(time.RFC3339),
		})

		requestLogger.Debug("Started processing request")

		err := c.Next()

		duration := time.Since(start)

		fields := logrus.Fields{
			"method":      c.Method(),
			"path":        c.Path(),
			"status":      c.Response().StatusCode(),
			"ip":          c.IP(),
			"user_agent":  c.Get("User-Agent"),
			"duration_ms": duration.Milliseconds(),
			"duration":    duration.String(),
		}

		if requestID := c.Get("X-Request-ID"); requestID != "" {
			fields["request_id"] = requestID
		}

		status := c.Response().StatusCode()
		switch {
		case status >= 500:
			if err != nil {
				fields["error"] = err.Error()
			}
			requestLogger.WithFields(fields).Error("Request failed with server error")
		case status >= 400:
			if err != nil {
				fields["error"] = err.Error()
			}
			requestLogger.WithFields(fields).Warn("Request failed with client error")
		case status >= 300:
			requestLogger.WithFields(fields).Info("Request redirected")
		default:
			requestLogger.WithFields(fields).Info("Request completed successfully")
		}

		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"error":       err.Error(),
				"method":      c.Method(),
				"path":        c.Path(),
				"status":      c.Response().StatusCode(),
				"duration_ms": duration.Milliseconds(),
			}).Error("Request handler returned error")
		}

		return err
	}
}
