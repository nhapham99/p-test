package services

import (
	_error "payment-module/error"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	"payment-module/internals/secure"
	"payment-module/logger"
	_util "payment-module/utils"

	"github.com/gofiber/fiber/v2"
)

func SendMessage(c *fiber.Ctx) (bool, *_error.SystemError) {
	logger.Info("SendMessage starting....")
	apiKey := c.Get(constants.AUTHORIZATION)

	configs, err := GetANotificationConfigByApiKey(c, apiKey)
	if err != nil {
		return false, _error.NewHandlerErrorWithMessage("Không thể lấy được thông tin cấu hình của ApiKey.")
	}
	if configs == nil || len(configs) != 1 {
		return false, _error.NewHandlerErrorWithMessage("Kiểm tra lại ApiKey và cấu hình.")
	}
	secretKey := configs[0].Secretkey
	token := configs[0].Telegram.Token
	chatId := configs[0].Telegram.ChatId
	var object models.Message

	//validate the request body
	if err := c.BodyParser(&object); err != nil {
		return false, _error.New(err)
	}

	secured, err2 := secure.ValidateChecksum(secretKey, object)
	if err2 != nil || !secured {
		return secured, err2
	}

	return _util.SendMessageTelegram(token, chatId, object.Message)
}
