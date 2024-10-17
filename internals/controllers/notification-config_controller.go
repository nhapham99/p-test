package controllers

import (
	"net/http"
	_error "payment-module/error"
	"payment-module/internals/responses"
	"payment-module/internals/services"
	"payment-module/logger"

	"github.com/gofiber/fiber/v2"
)

// GetAllNotificationConfig ... GetAllNotificationConfig
// @Summary GetAllNotificationConfig
// @Description GetAllNotificationConfig
// @Tags Notification
// @Success 200 {array} models.Notification
// @Failure 404 {object} object
// @Router / [get]
func GetAllNotificationConfig(c *fiber.Ctx) error {
	logger.Info("GetAllNotificationConfig starting....")
	object, err := services.GetAllNotificationConfig()
	if err != nil {
		return _error.HandleSystemError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: object})
}

// CreateNotificationConfig ... CreateNotificationConfig
// @Summary Create new Notification Config based on paramters
// @Description Create new Notification Config
// @Tags NotificationConfig
// @Accept json
// @Param Config body models.Config true "Config Data"
// @Success 200 {object} object
// @Failure 400,500 {object} object
// @Router / [post]
func CreateNotificationConfig(c *fiber.Ctx) error {
	result, err := services.CreateNotificationConfig(c)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	return c.Status(http.StatusCreated).JSON(responses.BaseResponse{Status: http.StatusCreated, Message: "success", Data: result})
}

func EditANotificationConfig(c *fiber.Ctx) error {
	result, err := services.EditANotificationConfig(c)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: result})
}

func DeleteANotificationConfig(c *fiber.Ctx) error {
	result, err := services.DeleteANotificationConfig(c)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	return c.Status(http.StatusOK).JSON(
		responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: result},
	)
}
