package controllers

import (
	"net/http"
	_error "payment-module/error"
	"payment-module/internals/responses"
	"payment-module/internals/services"
	"payment-module/logger"

	"github.com/gofiber/fiber/v2"
)

// SendMessage ... Send message to tele/mail/...
// @Summary Send Message
// @Description end message to tele/mail/...
// @Tags News
// @Success 200 {array} models.News
// @Failure 404 {object} object
// @Router / [get]
func SendMessage(c *fiber.Ctx) error {
	logger.Info("SendMessage starting....")
	result, err := services.SendMessage(c)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	if result {
		return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: result})
	} else {
		return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "false", Data: result})
	}
}
