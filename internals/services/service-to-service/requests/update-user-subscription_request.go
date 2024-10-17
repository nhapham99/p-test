package serviceToServiceRequests

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceUpdateUserSubscriptionRequest struct {
	UserId            primitive.ObjectID `json:"userId"`
	PremiumFrom       time.Time          `json:"premiumFrom"`
	PremiumTo         time.Time          `json:"premiumTo"`
	LastPurchaseToken *string            `json:"lastPurchaseToken"`
	Platform          int                `json:"platform"`
}
