package serviceToServiceResponses

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServiceGetPaymentMethodResponse struct {
	Id     primitive.ObjectID `json:"_id"`
	Name   string             `json:"name"`
	Active bool               `json:"active"`
}
