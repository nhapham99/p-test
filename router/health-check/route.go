package healthCheckRoutestes

import (
	"payment-module/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	HealthCheck := router.Group("/health-check")

	HealthCheck.Get("/", controllers.HealthCheck)
}
