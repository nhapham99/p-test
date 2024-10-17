package services

import (
	"context"
	"payment-module/configs"
	_error "payment-module/error"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	"payment-module/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// var categoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "category")
func getVersionCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetVersionCollection())
}

func GetCurrentVersion() int64 {
	version, err := getVersionConfig()
	if err != nil {
		return -1
	}
	return version.Version
}

func GetCurrentInterval() int {
	version, err := getVersionConfig()
	if err != nil {
		return -1
	}
	return version.Interval
}

func UpdateVersion() {
	deleteAllVersion()
	createVerion()
}

func getVersionConfig() (models.Version, *_error.SystemError) {
	logger.Info("GetVersionConfig starting....")
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()
	filter := bson.M{}
	results, err2 := getVersionCollection().Find(ctx, filter)
	if err2 != nil {
		return models.Version{}, _error.New(err2)
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleObject models.Version
		if err := results.Decode(&singleObject); err != nil {
			logger.Warnf("GetVersionConfig parse object fail: %s", err.Error())
			continue
		}
		logger.Infof("GetVersionConfig:%o", singleObject.Version)
		return singleObject, nil
	}

	return models.Version{}, nil
}

func createVerion() (*mongo.InsertOneResult, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	userName := "system"
	defer cancel()
	newObject := models.Version{
		Id:         primitive.NewObjectID(),
		Version:    time.Now().Unix(),
		CreatedAt:  time.Now().Unix(),
		UpdateAt:   time.Now().Unix(),
		CreateUser: userName,
		UpdateUser: userName,
	}

	result, err1 := getVersionCollection().InsertOne(ctx, newObject)
	if err1 != nil {
		return nil, _error.New(err1)
	}
	return result, nil
}

func editAVersionStatus(objectId string, version int64) (models.Version, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	userName := "admin"
	var object models.Config
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(objectId)
	filter := bson.M{"_id": objId}
	err := getVersionCollection().FindOne(ctx, filter).Decode(&object)
	if err != nil {
		return models.Version{}, _error.NewErrorByString("Version ISN'T EXISTS")
	}

	update := bson.M{
		"version":    version,
		"createdAt":  object.CreatedAt,
		"createUser": object.CreateUser,
		"updateAt":   time.Now().Unix(),
		"updateUser": userName,
	}

	result, err := getVersionCollection().UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return models.Version{}, _error.New(err)
	}

	//get updated object details
	var updatedObject models.Version
	if result.MatchedCount == 1 {
		err := getVersionCollection().FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedObject)
		if err != nil {
			return models.Version{}, _error.New(err)
		}
	}
	return updatedObject, nil
}

func deleteAllVersion() (string, *_error.SystemError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()
	filter := bson.M{}
	result, err2 := getVersionCollection().DeleteMany(ctx, filter)
	if err2 != nil {
		return "", _error.New(err2)
	}

	if result.DeletedCount < 1 {
		return "", _error.NewErrorByString("Verion NOT FOUND OR YOU DON'T HAVE PERMISSION")
	}

	return "Config successfully deleted!", nil
}
