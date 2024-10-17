package models

import (
	"payment-module/appstore"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppstoreServerNotificationRecord struct {
	Id               primitive.ObjectID       `json:"_id" bson:"_id,omitempty"`
	UserId           primitive.ObjectID       `json:"userId,omitempty" bson:"userId,omitempty"`
	AppAccountToken  string                   `json:"appAccountToken,omitempty" bson:"appAccountToken,omitempty"`
	NotificationType string                   `json:"notificationType,omitempty" bson:"notificationType,omitempty"`
	SubType          string                   `json:"subType,omitempty" bson:"subType,omitempty"`
	TransactionInfo  appstore.TransactionInfo `json:"transactionInfo,omitempty" bson:"transactionInfo,omitempty"`
	Version          string                   `json:"version,omitempty" bson:"version,omitempty"`
	CreatedAt        time.Time                `json:"createdAt" bson:"createdAt,omitempty"`
}
