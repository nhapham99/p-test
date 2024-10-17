package paymentRecordControllerRequests

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserTransferRequest struct {
	PaymentRecordId primitive.ObjectID `json:"paymentRecordId" validate:"required"`
	BankAccount     string             `json:"bankAccount" validate:"required"`
	AccountHolder   string             `json:"accountHolder" validate:"required"`
	BankName        string             `json:"bankName" validate:"required"`
}
