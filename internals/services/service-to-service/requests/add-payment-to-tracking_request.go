package serviceToServiceRequests

import (
	"time"
)

type ForeignKey struct {
	Id      string `json:"id" bson:"id,omitempty"`
	Name    string `json:"name,omitempty" bson:"name"`
	Version string `json:"version,omitempty" bson:"version"`
}

type AddPaymentToTrackingRequest struct {
	User          ForeignKey `json:"user,omitempty" bson:"user"`
	PaymentId     string     `json:"paymentId,omitempty" bson:"paymentId"`
	PaymentMethod ForeignKey `json:"paymentMethod,omitempty" bson:"paymentMethod"`
	PremiumFrom   time.Time  `json:"premiumFrom" bson:"premiumFrom"`
	PremiumTo     time.Time  `json:"premiumTo" bson:"premiumTo"`
}
