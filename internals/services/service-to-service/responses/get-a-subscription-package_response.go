package serviceToServiceResponses

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServiceGetSubscriptionPackageResponse struct {
	Id       primitive.ObjectID `json:"_id"`
	Name     string             `json:"name"`
	Active   bool               `json:"active"`
	Price    int64              `json:"price"`
	Duration int16              `json:"duration"`
}
