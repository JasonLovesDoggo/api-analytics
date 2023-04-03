package analytics

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tom-draper/api-analytics/analytics/go/core"
)

func Analytics(apiKey string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		data := core.RequestData{
			Hostname:     c.Hostname(),
			Path:         c.Path(),
			IPAddress:    c.IP(),
			UserAgent:    string(c.Request().Header.UserAgent()),
			Method:       c.Method(),
			Status:       c.Response().StatusCode(),
			ResponseTime: time.Since(start).Milliseconds(),
			CreatedAt:    start.Format(time.RFC3339),
		}

		core.LogRequest(apiKey, data, "Fiber")

		return err
	}
}
