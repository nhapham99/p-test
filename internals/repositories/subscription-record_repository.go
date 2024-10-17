package repositories

import (
	"context"
	"payment-module/configs"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	subscriptionRecordControllerRequests "payment-module/internals/requests/subscription-record-controller_requests"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// func getSubscriptionRecordCollection() *mongo.Collection {
// 	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetSubscriptionRecordCollection())
// }

func getSubscriptionRecordRequestCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetSubscriptionRecordRequestCollection())
}

func getSubscriptionRecordCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetSubscriptionRecordCollection())
}

func getAppStoreServerNotificationCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetAppStoreServerNotificationCollection())
}

func InsertOneAndroidSubscriptionRequestRecord(c *fiber.Ctx, userId primitive.ObjectID, object subscriptionRecordControllerRequests.AndroidClientPurchaseNotificationRequest, platform int) (models.SubscriptionRequestRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()
	newObject := models.SubscriptionRequestRecord{
		Id:             primitive.NewObjectID(),
		UserId:         userId,
		PurchaseToken:  object.PurchaseToken,
		BasePlanId:     object.BasePlanId,
		PurchaseStatus: object.PurchaseStatus,
		Platform:       platform,
		CreatedAt:      time.Now(),
	}

	_, err := getSubscriptionRecordRequestCollection().InsertOne(ctx, newObject)
	return newObject, err
}

func InsertOneIOSSubscriptionRequestRecord(c *fiber.Ctx, userId primitive.ObjectID, object subscriptionRecordControllerRequests.IOSClientPurchaseNotificationRequest, platform int) (models.SubscriptionRequestRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()
	newObject := models.SubscriptionRequestRecord{
		Id:              primitive.NewObjectID(),
		UserId:          userId,
		AppAccountToken: strings.ToLower(object.AppAccountToken),
		Platform:        platform,
		Status:          0,
		CreatedAt:       time.Now(),
	}

	_, err := getSubscriptionRecordRequestCollection().InsertOne(ctx, newObject)
	return newObject, err
}

func InsertOneGooglePlaySubscriptionRecord(c *fiber.Ctx, userId primitive.ObjectID, object subscriptionRecordControllerRequests.AndroidClientPurchaseNotificationRequest, platform int, premiumFrom, premiumTo time.Time) (models.SubscriptionRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()
	newObject := models.SubscriptionRecord{
		Id:             primitive.NewObjectID(),
		UserId:         userId,
		PurchaseToken:  object.PurchaseToken,
		BasePlanId:     object.BasePlanId,
		ProductId:      object.ProductId,
		PurchaseStatus: object.PurchaseStatus,
		PremiumFrom:    premiumFrom,
		PremiumTo:      premiumTo,
		Platform:       platform,
		CreatedAt:      time.Now(),
	}

	_, err := getSubscriptionRecordCollection().InsertOne(ctx, newObject)
	return newObject, err
}

func InsertOneAppStoreSubscriptionRecord(c *fiber.Ctx, userId primitive.ObjectID, appAccountToken string, platform int, premiumFrom, premiumTo time.Time) (models.SubscriptionRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()
	newObject := models.SubscriptionRecord{
		Id:              primitive.NewObjectID(),
		UserId:          userId,
		AppAccountToken: appAccountToken,
		PremiumFrom:     premiumFrom,
		PremiumTo:       premiumTo,
		Platform:        platform,
		CreatedAt:       time.Now(),
	}

	_, err := getSubscriptionRecordCollection().InsertOne(ctx, newObject)
	return newObject, err
}

func InsertOneAppStoreServerNotificationRecord(c *fiber.Ctx, payload models.AppstoreServerNotificationRecord) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()

	_, err := getAppStoreServerNotificationCollection().InsertOne(ctx, payload)
	return payload, err
}

func FindUserIdByAppAccountToken(c *fiber.Ctx, appAccountToken string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()

	var record models.SubscriptionRequestRecord
	err := getSubscriptionRecordRequestCollection().FindOne(ctx, primitive.M{"appAccountToken": strings.ToLower(appAccountToken)}).Decode(&record)
	return record.UserId, err
}
