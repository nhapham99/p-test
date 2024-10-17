package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRecord struct {
	Id                      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	User                    ForeignKey         `json:"user,omitempty" bson:"user"`
	PaymentMethod           ForeignKey         `json:"paymentMethod,omitempty" bson:"paymentMethod"`
	SubscriptionPackage     ForeignKey         `json:"subscriptionPackage,omitempty" bson:"subscriptionPackage"`
	Price                   int64              `json:"price,omitempty" bson:"price"`
	Duration                int16              `json:"duration,omitempty" bson:"duration"`
	Status                  int16              `json:"status,omitempty" bson:"status"`
	PaymentConfirmationDate *time.Time         `json:"paymentConfirmationDate,omitempty" bson:"paymentConfirmationDate"`
	OldPremiumFrom          *time.Time         `json:"oldPremiumFrom,omitempty" bson:"oldPremiumFrom"`
	OldPremiumTo            *time.Time         `json:"oldPremiumTo,omitempty" bson:"oldPremiumTo"`
	NewPremiumFrom          *time.Time         `json:"newPremiumFrom,omitempty" bson:"newPremiumFrom"`
	NewPremiumTo            *time.Time         `json:"newPremiumTo,omitempty" bson:"newPremiumTo"`
	BankInfo                *BankInfo          `json:"bankInfo,omitempty" bson:"bankInfo"`
	RecordInfo              RecordInfo         `json:"recordInfo,omitempty" bson:"recordInfo"`
}
