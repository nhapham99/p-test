package utils

import "github.com/gofiber/fiber/v2"

func ReadUserIP(c *fiber.Ctx) string {
	IPAddress := c.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = c.Get("X-Forwarded-For")
	}
	return IPAddress
}
