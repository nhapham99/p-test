package controllers

import (
	"net/http"
	"payment-module/internals/responses"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: time.Now().String()})
}
