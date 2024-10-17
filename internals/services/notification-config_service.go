package services

import (
	"context"
	"payment-module/configs"
	_error "payment-module/error"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	jwt "payment-module/jwtchecker"
	"payment-module/logger"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getNotificationConfigCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetNotificationConfigCollection())
}

var validateNotification = validator.New()

func CreateNotificationConfig(c *fiber.Ctx) (*mongo.InsertOneResult, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	userName := jwt.GetUserNameFromJwt(c)
	var object models.Notification
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&object); err != nil {
		return nil, _error.New(err)
	}

	//use the validator library to validate required fields
	if validationErr := validateNotification.Struct(&object); validationErr != nil {
		return nil, _error.New(validationErr)
	}

	newObject := models.Notification{
		Id:          primitive.NewObjectID(),
		Code:        object.Code,
		Name:        object.Name,
		Description: object.Description,
		ApiKey:      object.ApiKey,
		Secretkey:   object.Secretkey,
		Status:      object.Status,
		Telegram:    object.Telegram,
		CreatedAt:   time.Now().Unix(),
		UpdateAt:    time.Now().Unix(),
		CreateUser:  userName,
		UpdateUser:  userName,
	}

	result, err1 := getNotificationConfigCollection().InsertOne(ctx, newObject)
	if err1 != nil {
		return nil, _error.New(err1)
	}

	return result, nil
}

func EditANotificationConfig(c *fiber.Ctx) (models.Notification, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	objectId := c.Params("objectId")
	userName := jwt.GetUserNameFromJwt(c)
	var object models.Notification
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(objectId)

	//validate the request body
	if err := c.BodyParser(&object); err != nil {
		return models.Notification{}, _error.New(err)
	}

	//todo: lay object cu va tim rowStatus & RowStatusTst
	var oldObject models.Notification
	getNotificationConfigCollection().FindOne(ctx, bson.M{"_id": objId}).Decode(&oldObject)
	update := bson.M{
		"code":        object.Code,
		"name":        object.Name,
		"description": object.Description,
		"apiKey":      object.ApiKey,
		"secretKey":   object.Secretkey,
		"status":      object.Status,
		"telegram":    object.Telegram,
		"createdAt":   oldObject.CreatedAt,
		"updateAt":    time.Now().Unix(),
		"createUser":  oldObject.CreateUser,
		"updateUser":  userName,
	}

	result, err3 := getNotificationConfigCollection().UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err3 != nil {
		return models.Notification{}, _error.New(err3)
	}

	//get updated object details
	var updatedObject models.Notification
	if result.MatchedCount == 1 {
		err4 := getNotificationConfigCollection().FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedObject)
		if err4 != nil {
			return models.Notification{}, _error.New(err4)
		}
	}

	return updatedObject, nil
}

func DeleteANotificationConfig(c *fiber.Ctx) (string, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	objectId := c.Params("objectId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(objectId)

	filter := bson.M{"_id": objId}
	result, err2 := getNotificationConfigCollection().DeleteOne(ctx, filter)
	if err2 != nil {
		return "", _error.New(err2)
	}

	if result.DeletedCount < 1 {
		return "", _error.NewErrorByString("Config NOT FOUND OR YOU DON'T HAVE PERMISSION")
	}

	return "Config successfully deleted!", nil
}

func GetAllNotificationConfig() ([]models.Notification, *_error.SystemError) {
	logger.Debug("GetAllNotificationConfig starting....")
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	var objects []models.Notification
	defer cancel()

	filter := bson.M{}
	results, err2 := getNotificationConfigCollection().Find(ctx, filter)
	if err2 != nil {
		return nil, _error.New(err2)
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleObject models.Notification
		if err := results.Decode(&singleObject); err != nil {
			logger.Warnf("GetAllNotificationConfig parse object fail: %s", err.Error())
			continue
		}

		objects = append(objects, singleObject)
	}

	return objects, nil
}

func GetANotificationConfigByApiKey(c *fiber.Ctx, apiKey string) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	var objects []models.Notification
	defer cancel()
	filter := bson.M{"apiKey": apiKey}

	results, err := getNotificationConfigCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for results.Next(ctx) {
		var singleObject models.Notification
		if err = results.Decode(&singleObject); err != nil {
			logger.Warnf("GetAConfigByCode parse object fail: %s", err.Error())
			continue
		}

		objects = append(objects, singleObject)
	}
	return objects, err
}
