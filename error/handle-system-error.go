package errors

import (
	"net/http"
	"payment-module/configs"
	"payment-module/internals/responses"

	"github.com/gofiber/fiber/v2"
)

// Handle lỗi hệ thống, luu y: phai chac chan "err != nil"
func HandleSystemError(c *fiber.Ctx, err *SystemError) error {
	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BaseResponse{Status: http.StatusInternalServerError, Message: configs.GetSystemErrorMessage(), Data: ""})
	}
	if err.ErrorCode() == DB_NOT_FOUND {
		return c.Status(http.StatusForbidden).JSON(responses.BaseResponse{Status: http.StatusForbidden, Message: err.ErrorMessage(), Data: ""})
	} else if err.ErrorCode() == CHECKSUM_INVALID {
		return c.Status(http.StatusUnauthorized).JSON(responses.BaseResponse{Status: http.StatusUnauthorized, Message: err.ErrorMessage(), Data: err.ErrorMessage()})
	} else if err.ErrorCode() == HANDLER_ERROR_WITH_MESSAGE {
		return c.Status(http.StatusBadRequest).JSON(responses.BaseResponse{Status: http.StatusBadRequest, Message: err.ErrorMessage(), Data: err.ErrorMessage()})
	} else if err.ErrorMessage() != "" {
		return c.Status(http.StatusBadRequest).JSON(responses.BaseResponse{Status: http.StatusBadRequest, Message: err.ErrorMessage(), Data: err.ErrorMessage()})
	} else {
		return c.Status(http.StatusInternalServerError).JSON(responses.BaseResponse{Status: http.StatusInternalServerError, Message: err.ErrorMessage(), Data: ""})
	}
}
