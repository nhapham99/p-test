package repositories

import (
	"context"
	"payment-module/configs"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	"payment-module/jwtchecker"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getPaymentRecordCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, configs.GetConfigDatabaseName(), configs.GetPaymentRecordCollection())
}

func InsertOnePaymentRecord(c *fiber.Ctx, paymentRecord models.PaymentRecord) (*mongo.InsertOneResult, models.PaymentRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()
	now := time.Now()
	newObject := paymentRecord
	newObject.Id = primitive.NewObjectID()
	newObject.RecordInfo.CreatedAt = &now
	newObject.RecordInfo.CreatedBy = models.ForeignKey{
		Id:   jwtchecker.GetUserIdFromJwt(c),
		Name: jwtchecker.GetUserNameFromJwt(c),
	}

	result, err := getPaymentRecordCollection().InsertOne(ctx, newObject)
	return result, newObject, err
}

func CountPaymentRecords(filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()

	result, err := getPaymentRecordCollection().CountDocuments(ctx, filter)
	return result, err
}

func FindOnePaymentRecord(filter bson.M) (*models.PaymentRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	var object models.PaymentRecord
	defer cancel()

	err := getPaymentRecordCollection().FindOne(ctx, filter).Decode(&object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func FindListPaymentRecords(filter bson.M) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()

	result, err := getPaymentRecordCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdatePaymentRecords(c *fiber.Ctx, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()

	now := time.Now()
	update["recordInfo.updatedAt"] = &now
	update["recordInfo.updatedBy"] = models.ForeignKey{
		Id:   jwtchecker.GetUserIdFromJwt(c),
		Name: jwtchecker.GetUserNameFromJwt(c),
	}
	result, err := getPaymentRecordCollection().UpdateMany(ctx, filter, bson.M{"$set": update})
	return result, err
}

func FindListPaymentRecordsWithPipeline(pipeline []bson.M) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	defer cancel()

	result, err := getPaymentRecordCollection().Aggregate(ctx, pipeline)
	return result, err
}
