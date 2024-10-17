package HealthCheckRouter

import (
	"payment-module/internals/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	news := router.Group("/config")

	news.Get("/", controllers.GetAllConfig)
	news.Post("/", controllers.CreateConfig)
	news.Put("/:objectId", controllers.EditAConfig)
	news.Delete("/:objectId", controllers.DeleteAConfig)
}
