package NotificationRouter

import (
	"payment-module/internals/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	news := router.Group("/notification")

	news.Get("/", controllers.GetAllNotificationConfig)
	news.Post("/", controllers.CreateNotificationConfig)
	news.Put("/:objectId", controllers.EditANotificationConfig)
	news.Delete("/:objectId", controllers.DeleteANotificationConfig)
}
