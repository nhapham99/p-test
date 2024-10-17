package paymentRecordControllerRequests

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateNewPaymentRecordRequest struct {
	// Thông tin chung
	PaymentMethodId       primitive.ObjectID `json:"paymentMethodId" validate:"required"`
	SubscriptionPackageId primitive.ObjectID `json:"subscriptionPackageId" validate:"required"`
	// Thông tin liên quan đến thanh toán bằng chuyển khoản
	BankName    string `json:"bankName"`
	BankAccount string `json:"bankAccount"`
	Message     string `json:"message"`
}
