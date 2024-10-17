package configs

import (
	"context"
	"fmt"
	"payment-module/internals/constants"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	fmt.Println("Try to connect to MongoDB:" + EnvMongoURI())
	mongodbOption := options.Client().ApplyURI(EnvMongoURI()).SetTimeout(constants.TIME_OUT_CONNECTION * time.Second)
	client, err := mongo.NewClient(mongodbOption)
	if err != nil {
		fmt.Println(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), constants.TIME_OUT_EXECUTE*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}
