package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionRequestRecord struct {
	Id              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId          primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	PurchaseToken   string             `json:"purchaseToken,omitempty" bson:"purchaseToken,omitempty"`
	BasePlanId      string             `json:"basePlanId,omitempty" bson:"basePlanId,omitempty"`
	PurchaseStatus  string             `json:"purchaseStatus,omitempty" bson:"purchaseStatus,omitempty"`
	Platform        int                `json:"platform" bson:"platform"` // 0: Android, 1: iOS
	AppAccountToken string             `json:"appAccountToken,omitempty" bson:"appAccountToken,omitempty"`
	Status          int                `json:"status" bson:"status,omitempty"` // 0: pending, 1: success, 2: failed
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
}
