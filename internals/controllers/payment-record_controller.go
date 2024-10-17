package controllers

import (
	"net/http"
	_error "payment-module/error"
	"payment-module/internals/constants"
	paymentRecordControllerRequests "payment-module/internals/requests/payment-record-controller_requests"
	"payment-module/internals/responses"
	"payment-module/internals/services"
	baseService "payment-module/internals/services/base"
	serviceToService "payment-module/internals/services/service-to-service"
	"payment-module/jwtchecker"
	"payment-module/logger"
	vnpayService "payment-module/vnpay/services"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var validatePayment = validator.New()

func CreatePaymentRecord(c *fiber.Ctx) error {
	logger.Info("CreatePaymentRecord starting....")

	userId := jwtchecker.GetUserIdFromJwt(c)
	userInfo, err1 := serviceToService.GetUserInfo(userId)
	if err1 != nil {
		return _error.HandleSystemError(c, _error.New(err1))
	}

	if userInfo.PremiumTo != nil && userInfo.PremiumTo.After(time.Now()) && userInfo.LastPayment != nil && (userInfo.LastPayment.Platform == constants.PLATFORMS.Googleplay || userInfo.LastPayment.Platform == constants.PLATFORMS.Appstore) {
		return c.Status(http.StatusBadRequest).JSON(responses.BaseResponse{Status: http.StatusBadRequest, Message: "User already has premium subscription on other platform", Data: nil})
	}

	var object paymentRecordControllerRequests.CreateNewPaymentRecordRequest

	if err := c.BodyParser(&object); err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	if validationErr := validatePayment.Struct(&object); validationErr != nil {
		return _error.HandleSystemError(c, _error.New(validationErr))
	}

	paymentMethodInfo, subscriptionInfo, err := services.CheckPaymentMethod(c, object)
	if err != nil {
		return _error.HandleSystemError(c, err)
	}
	if paymentMethodInfo.Id.Hex() == constants.PAYMENT_METHOD_TRANSFER_ID {
		result, err := baseService.CreatePaymentRecord(c, object, paymentMethodInfo, subscriptionInfo)
		if err != nil {
			return _error.HandleSystemError(c, err)
		}
		return c.Status(http.StatusCreated).JSON(responses.BaseResponse{Status: http.StatusCreated, Message: "success", Data: result})
	} else if paymentMethodInfo.Id.Hex() == constants.PAYMENT_METHOD_VNPAY_ID {
		result, err := vnpayService.CreateVnPayTransaction(c, object, paymentMethodInfo, subscriptionInfo)
		if err != nil {
			return _error.HandleSystemError(c, err)
		}
		return c.Status(http.StatusSeeOther).JSON(responses.BaseResponse{Status: http.StatusSeeOther, Message: "success", Data: result})
	} else {
		return _error.HandleSystemError(c, _error.NewErrorByString(_error.PAYMENT_E001_005))
	}
}

func UserGetAPaymentRecord(c *fiber.Ctx) error {
	logger.Info("UserGetAPaymentRecord starting....")

	result, err := services.UserGetAPaymentRecord(c)

	if err != nil {
		return _error.HandleSystemError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: result})
}

func AdmTransferConfirm(c *fiber.Ctx) error {
	logger.Info("AdmTransferConfirm starting....")

	objectId := c.Params("objectId")

	err := baseService.AdmTransferConfirm(c, objectId)

	if err != nil {
		return _error.HandleSystemError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: nil})
}

func AdminGetAllPaymentRecordsPagination(c *fiber.Ctx) error {
	logger.Info("AdminGetAllPaymentRecordsPagination starting....")

	result, err := services.AdminGetAllPaymentRecordsPagination(c)

	if err != nil {
		return _error.HandleSystemError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "success", Data: result})
}
