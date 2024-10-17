package serviceToServiceRequests

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServiceUpdateUserPremiumRequest struct {
	UserId    primitive.ObjectID `json:"userId"`
	PaymentId primitive.ObjectID `json:"paymentId"`
	Duration  int16              `json:"duration"`
	Platform  int                `json:"platform"`
}
