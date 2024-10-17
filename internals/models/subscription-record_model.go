package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionRecord struct {
	Id              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId          primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	PurchaseToken   string             `json:"purchaseToken,omitempty" bson:"purchaseToken,omitempty"`
	ProductId       string             `json:"productId,omitempty" bson:"productId,omitempty"`
	BasePlanId      string             `json:"basePlanId,omitempty" bson:"basePlanId,omitempty"`
	PurchaseStatus  string             `json:"purchaseStatus,omitempty" bson:"purchaseStatus,omitempty"`
	AppAccountToken string             `json:"appAccountToken,omitempty" bson:"appAccountToken,omitempty"`
	PremiumFrom     time.Time          `json:"premiumFrom,omitempty" bson:"premiumFrom,omitempty"`
	PremiumTo       time.Time          `json:"premiumTo,omitempty" bson:"premiumTo,omitempty"`
	Platform        int                `json:"platform,omitempty" bson:"platform"` // 0: Android, 1: iOS
	CreatedAt       time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
