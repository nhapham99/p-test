package controllers

import (
	"net/http"
	"payment-module/appstore"
	_error "payment-module/error"
	"payment-module/internals/constants"
	subscriptionRecordControllerRequests "payment-module/internals/requests/subscription-record-controller_requests"
	"payment-module/internals/responses"
	"payment-module/internals/services"
	serviceToService "payment-module/internals/services/service-to-service"
	"payment-module/jwtchecker"
	"payment-module/logger"
	"payment-module/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AndroidClientPurchaseNotification(c *fiber.Ctx) error {
	logger.Info("AndroidClientPurchaseNotification starting....")

	userId := jwtchecker.GetUserIdFromJwt(c)
	userInfo, err := serviceToService.GetUserInfo(userId)
	if err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	if userInfo.PremiumTo != nil && userInfo.PremiumTo.After(time.Now()) && userInfo.LastPayment != nil && userInfo.LastPayment.Platform != constants.PLATFORMS.Googleplay {
		return c.Status(http.StatusBadRequest).JSON(responses.BaseResponse{Status: http.StatusBadRequest, Message: "User already has premium subscription on other platform", Data: nil})
	}

	var object subscriptionRecordControllerRequests.AndroidClientPurchaseNotificationRequest
	if err := c.BodyParser(&object); err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	record, err := services.HandleGooglePlayClientPurchaseNotification(c, userId, userInfo.PremiumFrom, userInfo.PremiumTo, object)
	logger.Info("record: ", record)
	if err != nil {
		logger.Error("HandleGooglePlayClientPurchaseNotification error: ", err)
		return _error.HandleSystemError(c, _error.New(err))
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "Google play subscription success", Data: record})
}

func CheckGooglePlayRenew(c *fiber.Ctx) error {
	logger.Info("CheckGooglePlayRenew starting....")
	userId := jwtchecker.GetUserIdFromJwt(c)
	user, err := serviceToService.GetUserInfo(userId)
	if err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	if user.PremiumTo != nil && user.LastPayment != nil && user.LastPayment.Platform == constants.PLATFORMS.Googleplay && user.LastPayment.SubscriptionToken != nil {
		isRenew, err := services.CheckGooglePlayRenew(c, userId, *user.PremiumFrom, *user.PremiumTo, *user.LastPayment.SubscriptionToken)
		if err != nil {
			return _error.HandleSystemError(c, _error.New(err))
		}

		if isRenew {
			return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "User has renew subscription", Data: isRenew})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "User not renew subscription", Data: false})
}

func IOSClientPurchaseNotification(c *fiber.Ctx) error {
	logger.Info("IOSClientPurchaseNotification starting....")

	userId := jwtchecker.GetUserIdFromJwt(c)
	userInfo, err := serviceToService.GetUserInfo(userId)
	if err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	if userInfo.PremiumTo != nil && userInfo.PremiumTo.After(time.Now()) && userInfo.LastPayment != nil && userInfo.LastPayment.Platform != constants.PLATFORMS.Appstore {
		return c.Status(http.StatusBadRequest).JSON(responses.BaseResponse{Status: http.StatusBadRequest, Message: "User already has premium subscription on other platform", Data: nil})
	}

	var object subscriptionRecordControllerRequests.IOSClientPurchaseNotificationRequest
	if err := c.BodyParser(&object); err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	record, err := services.HandleAppStoreClientPurchaseNotification(c, userId, object)
	if err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "Appstore client request save success, waiting for approve from appstore server", Data: record})
}

func AppStoreServerNotification(c *fiber.Ctx) error {
	logger.Info("AppStoreServerNotification starting....")

	clientIp := utils.ReadUserIP(c)

	if !appstore.IsAppStoreWhitelistedIp(clientIp) {
		return c.Status(http.StatusForbidden).JSON(responses.BaseResponse{Status: http.StatusForbidden, Message: "IP address not whitelisted", Data: nil})
	}

	var object appstore.AppStoreServerRequest

	if err := c.BodyParser(&object); err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	record, err := services.HandleAppStoreServerNotification(c, object.SignedPayload)
	if err != nil {
		return _error.HandleSystemError(c, _error.New(err))
	}

	return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "Transaction successfully", Data: record})
}
