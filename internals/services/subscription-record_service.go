package services

import (
	"encoding/json"
	"errors"
	"payment-module/appstore"
	googleplay "payment-module/googleplay"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	"payment-module/internals/repositories"
	srcRequests "payment-module/internals/requests/subscription-record-controller_requests"
	serviceToService "payment-module/internals/services/service-to-service"
	serviceToServiceResponses "payment-module/internals/services/service-to-service/responses"
	"payment-module/logger"
	"payment-module/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleGooglePlayClientPurchaseNotification(c *fiber.Ctx, userId primitive.ObjectID, currentPremiumFrom, currentPremiumTo *time.Time, object srcRequests.AndroidClientPurchaseNotificationRequest) (models.SubscriptionRecord, error) {
	logger.Info("HandleGooglePlayClientPurchaseNotification starting....")
	var record models.SubscriptionRecord

	_, err := SaveAndroidClientPurchaseRequest(c, userId, object)
	if err != nil {
		return record, errors.New("save googleplay client request failed")
	}

	res, isPurchaseValid, err := googleplay.CheckGoogleAuth(object.PurchaseToken)
	if err != nil {
		return record, errors.New("check googleplay client purchase failed")
	}

	if !isPurchaseValid {
		return record, errors.New("purchase token is not valid")
	} else {
		if object.PurchaseStatus == "purchased" {
			if _, ok := constants.BASE_PLANS[object.BasePlanId]; !ok {
				return record, errors.New("base plan id is not valid")
			}

			subscriptionInfo, ok := (*res).(map[string]interface{})
			if !ok {
				return record, errors.New("failed to convert res to map[string]interface{}")
			}
			startTimeString, ok := subscriptionInfo["startTime"].(string)
			if !ok {
				return record, errors.New("failed to convert startTime to string")
			}
			startTime, err := time.Parse(time.RFC3339, startTimeString)
			if err != nil {
				return record, err
			}
			lineItems, ok := subscriptionInfo["lineItems"].([]interface{})
			if !ok {
				return record, errors.New("failed to convert lineItems to []interface{}")
			}
			var expiryTime time.Time
			for _, item := range lineItems {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					return record, errors.New("failed to convert item to map[string]interface{}")
				}
				if itemMap["productId"] == "premium" {
					expiryTimeString, ok := itemMap["expiryTime"].(string)
					if !ok {
						return record, errors.New("failed to convert expiryTime to string")
					}
					expiryTime, err = time.Parse(time.RFC3339, expiryTimeString)
					if err != nil {
						return record, err
					}

					break
				}
			}

			// update subscription to user
			_, err = serviceToService.UpdateUserSubscription(userId, startTime, expiryTime, &object.PurchaseToken, constants.PLATFORMS.Googleplay)
			if err != nil {
				return record, err
			}

			// save subscription record
			record, err := repositories.InsertOneGooglePlaySubscriptionRecord(c, userId, object, constants.PLATFORMS.Googleplay, startTime, expiryTime)
			if err != nil {
				return record, err
			}

			// acknowledge subscription
			err = googleplay.SubscriptionAcknowledge(object.PurchaseToken, object.BasePlanId)
			if err != nil {
				return record, err
			}

			return record, nil
		} else {
			return record, errors.New("purchase status is not purchased")
		}
	}
}

func CheckGooglePlayRenew(c *fiber.Ctx, userId primitive.ObjectID, currentPremiumFrom, currentPremiumTo time.Time, lastPurchaseToken string) (bool, error) {
	res, isPurchaseValid, err := googleplay.CheckGoogleAuth(lastPurchaseToken)
	if err != nil {
		return false, errors.New("check googleplay client purchase failed")
	}

	if !isPurchaseValid {
		return false, errors.New("purchase token is not valid")
	} else {
		if res != nil {
			subscriptionInfo, ok := (*res).(map[string]interface{})
			if !ok {
				return false, errors.New("failed to convert res to map[string]interface{}")
			}
			lineItems, ok := subscriptionInfo["lineItems"].([]interface{})
			if !ok {
				return false, errors.New("failed to convert lineItems to []interface{}")
			}
			startTimeString, ok := subscriptionInfo["startTime"].(string)
			if !ok {
				return false, errors.New("failed to convert startTime to string")
			}
			startTime, err := time.Parse(time.RFC3339, startTimeString)
			if err != nil {
				return false, err
			}
			for _, item := range lineItems {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					return false, errors.New("failed to convert item to map[string]interface{}")
				}
				if itemMap["productId"] == "premium" {
					expiryTimeString, ok := itemMap["expiryTime"].(string)
					if !ok {
						return false, errors.New("failed to convert expiryTime to string")
					}
					expiryTime, err := time.Parse(time.RFC3339, expiryTimeString)
					if err != nil {
						return false, err
					}

					if expiryTime.After(currentPremiumTo.Add(10 * time.Second)) {
						// update subscription to user
						_, err = serviceToService.UpdateUserSubscription(userId, startTime, expiryTime, &lastPurchaseToken, constants.PLATFORMS.Googleplay)
						if err != nil {
							return false, err
						}

						// save subscription record
						_, err := repositories.InsertOneGooglePlaySubscriptionRecord(c, userId, srcRequests.AndroidClientPurchaseNotificationRequest{
							PurchaseToken:  lastPurchaseToken,
							PurchaseStatus: "RENEWED",
						}, constants.PLATFORMS.Googleplay, startTime, expiryTime)
						if err != nil {
							return false, err
						}
						return true, nil
					} else {
						return false, nil
					}
				}
			}
		} else {
			return false, errors.New("purchase token is not valid")
		}
	}

	return false, nil
}

func SaveAndroidClientPurchaseRequest(c *fiber.Ctx, userId primitive.ObjectID, request srcRequests.AndroidClientPurchaseNotificationRequest) (models.SubscriptionRequestRecord, error) {
	record, err := repositories.InsertOneAndroidSubscriptionRequestRecord(c, userId, request, constants.PLATFORMS.Googleplay)
	return record, err
}

func HandleAppStoreClientPurchaseNotification(c *fiber.Ctx, userId primitive.ObjectID, object srcRequests.IOSClientPurchaseNotificationRequest) (models.SubscriptionRequestRecord, error) {
	record, err := repositories.InsertOneIOSSubscriptionRequestRecord(c, userId, object, constants.PLATFORMS.Appstore)
	if err != nil {
		return record, errors.New("save appstore client request failed")
	}

	return record, nil
}

func HandleAppStoreServerNotification(c *fiber.Ctx, signedPayload string) (any, error) {
	var record serviceToServiceResponses.ServiceUpdateUserPremiumResponse

	_, payload, _, err := utils.ParseJWS(signedPayload)
	if err != nil {
		return record, err
	}
	decodedBytes, err := utils.Base64UrlDecode(payload)
	if err != nil {
		return record, err
	}
	var notificationPayload appstore.NotificationPayload
	err = json.Unmarshal(decodedBytes, &notificationPayload)
	if err != nil {
		return record, err
	}

	_, transactionPayload, _, err := utils.ParseJWS(notificationPayload.Data.SignedTransactionInfo)
	if err != nil {
		return record, err
	}

	decodedTransactionBytes, err := utils.Base64UrlDecode(transactionPayload)
	if err != nil {
		return record, err
	}

	var transactionInfo appstore.TransactionInfo
	err = json.Unmarshal(decodedTransactionBytes, &transactionInfo)
	if err != nil {
		return record, err
	}

	logger.Info("signedPayload starting....", signedPayload)
	logger.Info("transactionInfo.AppAccountToken starting....", transactionInfo.AppAccountToken)
	userId, err := repositories.FindUserIdByAppAccountToken(c, transactionInfo.AppAccountToken)
	logger.Info("userId starting....", userId)
	if err != nil {
		logger.Error(err)
		return record, err
	}

	var appstoreServerNotificationRecord = models.AppstoreServerNotificationRecord{
		UserId:           userId,
		AppAccountToken:  transactionInfo.AppAccountToken,
		NotificationType: notificationPayload.NotificationType,
		SubType:          notificationPayload.Subtype,
		TransactionInfo:  transactionInfo,
		Version:          notificationPayload.Version,
		CreatedAt:        time.Now(),
	}
	_, err = repositories.InsertOneAppStoreServerNotificationRecord(c, appstoreServerNotificationRecord)
	if err != nil {
		return record, err
	}

	if !utils.Contains(constants.APPSTORE_HANDLED_TYPES, notificationPayload.NotificationType) {
		return record, nil
	}

	var premiumFrom, premiumTo time.Time
	if notificationPayload.NotificationType == "SUBSCRIBED" {
		premiumFrom = time.Now()
	} else if notificationPayload.NotificationType == "DID_RENEW" || notificationPayload.NotificationType == "DID_CHANGE_RENEWAL_PREF" {
		premiumFrom = time.Unix(int64(transactionInfo.PurchaseDate)/1000, 0)
	}
	premiumTo = time.Unix(int64(transactionInfo.ExpiresDate)/1000, 0)
	_, err = serviceToService.UpdateUserSubscription(userId, premiumFrom, premiumTo, &transactionInfo.AppAccountToken, constants.PLATFORMS.Appstore)
	if err != nil {
		return record, err
	}

	_, err = repositories.InsertOneAppStoreSubscriptionRecord(c, userId, transactionInfo.AppAccountToken, constants.PLATFORMS.Appstore, premiumFrom, premiumTo)
	if err != nil {
		return record, errors.New("save appstore subscription record failed")
	}

	return record, err
}
