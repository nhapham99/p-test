package NotificationRouter

import (
	"payment-module/internals/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	news := router.Group("/telegram")
	news.Post("/sendMessage", controllers.SendMessage)
}
