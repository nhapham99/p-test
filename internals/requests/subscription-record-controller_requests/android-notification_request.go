package paymentRecordControllerRequests

type AndroidClientPurchaseNotificationRequest struct {
	PurchaseToken  string `json:"purchaseToken" validate:"required"`
	BasePlanId     string `json:"basePlanId" validate:"required"`
	ProductId      string `json:"productId" validate:"required"`
	PurchaseStatus string `json:"purchaseStatus" validate:"required"`
}
