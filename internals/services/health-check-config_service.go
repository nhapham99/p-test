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

// var categoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "category")
func getConfigCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetConfigCollection())
}

var validate = validator.New()

func CreateConfig(c *fiber.Ctx) (*mongo.InsertOneResult, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	userName := jwt.GetUserNameFromJwt(c)
	var object models.Config
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&object); err != nil {
		return nil, _error.New(err)
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&object); validationErr != nil {
		return nil, _error.New(validationErr)
	}

	newObject := models.Config{
		Id:          primitive.NewObjectID(),
		Code:        object.Code,
		Name:        object.Name,
		Description: object.Description,
		Status:      object.Status,
		Services:    object.Services,
		Telegram:    object.Telegram,
		CreatedAt:   time.Now().Unix(),
		UpdateAt:    time.Now().Unix(),
		CreateUser:  userName,
		UpdateUser:  userName,
	}

	result, err1 := getConfigCollection().InsertOne(ctx, newObject)
	if err1 != nil {
		return nil, _error.New(err1)
	}

	return result, nil
}

func EditAConfig(c *fiber.Ctx) (models.Config, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	objectId := c.Params("objectId")
	userName := jwt.GetUserNameFromJwt(c)
	var object models.Config
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(objectId)

	//validate the request body
	if err := c.BodyParser(&object); err != nil {
		return models.Config{}, _error.New(err)
	}

	//todo: lay object cu va tim rowStatus & RowStatusTst
	var oldObject models.Config
	getConfigCollection().FindOne(ctx, bson.M{"_id": objId}).Decode(&oldObject)
	update := bson.M{
		"code":        object.Code,
		"name":        object.Name,
		"description": object.Description,
		"status":      object.Status,
		"services":    object.Services,
		"telegram":    object.Telegram,
		"createdAt":   oldObject.CreatedAt,
		"updateAt":    time.Now().Unix(),
		"createUser":  oldObject.CreateUser,
		"updateUser":  userName,
	}

	result, err3 := getConfigCollection().UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err3 != nil {
		return models.Config{}, _error.New(err3)
	}

	//get updated object details
	var updatedObject models.Config
	if result.MatchedCount == 1 {
		err4 := getConfigCollection().FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedObject)
		if err4 != nil {
			return models.Config{}, _error.New(err4)
		}
	}

	return updatedObject, nil
}

func DeleteAConfig(c *fiber.Ctx) (string, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	objectId := c.Params("objectId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(objectId)

	filter := bson.M{"_id": objId}
	result, err2 := getConfigCollection().DeleteOne(ctx, filter)
	if err2 != nil {
		return "", _error.New(err2)
	}

	if result.DeletedCount < 1 {
		return "", _error.NewErrorByString("Config NOT FOUND OR YOU DON'T HAVE PERMISSION")
	}

	return "Config successfully deleted!", nil
}

// func GetAConfigById(c *fiber.Ctx, hospitalCode string, objectId string) (models.Config, *_error.SystemError) {
// 	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
// 	var result models.Config
// 	defer cancel()
// 	objId, _ := primitive.ObjectIDFromHex(objectId)

// 	collection, err := GetConfigCollection(c, hospitalCode)
// 	if err != nil {
// 		return models.Config{}, _error.NewDatabaseNotFound()
// 	}
// 	err2 := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
// 	if err2 != nil {
// 		return models.Config{}, _error.New(err2)
// 	}
// 	return result, err
// }

// func GetAConfigByCode(c *fiber.Ctx, hospitalCode string, code string, objectId string) ([]models.Config, *_error.SystemError) {
// 	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
// 	var objects []models.Config
// 	defer cancel()
// 	filter := bson.M{}
// 	if code != "" {
// 		filter["code"] = code
// 	}
// 	if objectId != "" {
// 		objId, _ := primitive.ObjectIDFromHex(objectId)
// 		filter["title"] = bson.M{"_id": objId}
// 	}

// 	collection, err := GetConfigCollection(c, hospitalCode)
// 	if err != nil {
// 		return nil, _error.NewDatabaseNotFound()
// 	}
// 	results, err2 := collection.Find(ctx, filter)
// 	if err2 != nil {
// 		return nil, _error.New(err2)
// 	}
// 	for results.Next(ctx) {
// 		var singleObject models.Config
// 		if err2 = results.Decode(&singleObject); err != nil {
// 			logger.Warnf("GetAConfigByCode parse object fail: %s", err2.Error())
// 			continue
// 		}

// 		objects = append(objects, singleObject)
// 	}
// 	return objects, err
// }

func GetAllConfig() ([]models.Config, *_error.SystemError) {
	logger.Info("GetAllConfig starting....")
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	var objects []models.Config
	defer cancel()

	filter := bson.M{}
	results, err2 := getConfigCollection().Find(ctx, filter)
	if err2 != nil {
		return nil, _error.New(err2)
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleObject models.Config
		if err := results.Decode(&singleObject); err != nil {
			logger.Warnf("GetAllConfig parse object fail: %s", err.Error())
			continue
		}

		objects = append(objects, singleObject)
	}

	return objects, nil
}
