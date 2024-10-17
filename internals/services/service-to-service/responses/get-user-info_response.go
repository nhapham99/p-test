package serviceToServiceResponses

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountInfo struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type LastPayment struct {
	PaymentId         *primitive.ObjectID `json:"paymentId" bson:"paymentId"`
	SubscriptionToken *string             `json:"subscriptionToken" bson:"subscriptionToken"`
	VoucherId         *string             `json:"voucherId" bson:"voucherId"`
	Platform          int                 `json:"platform" bson:"platform"`
	UpdatedAt         time.Time           `json:"updatedAt" bson:"updatedAt"`
}

type ServiceGetUserInfoResponse struct {
	AccountInfo AccountInfo  `json:"accountInfo" bson:"accountInfo"`
	PremiumFrom *time.Time   `json:"premiumFrom"`
	PremiumTo   *time.Time   `json:"premiumTo"`
	LastPayment *LastPayment `json:"lastPayment" bson:"lastPayment"`
}
