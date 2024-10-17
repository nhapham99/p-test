package controllers

import (
	"net/http"
	_error "payment-module/error"
	"payment-module/internals/responses"
	"payment-module/internals/services"
	"payment-module/logger"

	"github.com/gofiber/fiber/v2"
)

// GetNews ... Get all news
// @Summary Get all news
// @Description get all news
// @Tags News
// @Success 200 {array} models.News
// @Failure 404 {object} object
// @Router / [get]
func GetAllConfig(c *fiber.Ctx) error {
	logger.Info("GetAllNews starting....")
	object, err := services.GetAllConfig()
	if err != nil {
		return _error.HandleSystemError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: object})
}

// CreateConfig ... Create Config
// @Summary Create new Config based on paramters
// @Description Create new Config
// @Tags Config
// @Accept json
// @Param Config body models.Config true "Config Data"
// @Success 200 {object} object
// @Failure 400,500 {object} object
// @Router / [post]
func CreateConfig(c *fiber.Ctx) error {
	result, err := services.CreateConfig(c)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	services.UpdateVersion()
	return c.Status(http.StatusCreated).JSON(responses.BaseResponse{Status: http.StatusCreated, Message: "success", Data: result})
}

func EditAConfig(c *fiber.Ctx) error {
	result, err := services.EditAConfig(c)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	services.UpdateVersion()
	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: result})
}

func DeleteAConfig(c *fiber.Ctx) error {
	result, err := services.DeleteAConfig(c)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	services.UpdateVersion()
	return c.Status(http.StatusOK).JSON(
		responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: result},
	)
}

// GetConfigByID ... Get the Config by id
// @Summary Get one Config
// @Description get Config by ID
// @Tags Config
// @Param id path string true "Config ID"
// @Success 200 {object} models.Config
// @Failure 400,404 {object} object
// @Router /{id} [get]
// func GetAConfig(c *fiber.Ctx) error {
// 	objectId := c.Params("objectId")
// 	hospitalCode := c.Get("hospitalCode")
// 	object, err := services.GetAConfigById(c, hospitalCode, objectId)
// 	if err != nil {
// 		return _error.HandleSystemError(c, err)
// 	}

// 	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: object})
// }
