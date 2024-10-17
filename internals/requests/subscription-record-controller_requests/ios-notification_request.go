package paymentRecordControllerRequests

type IOSClientPurchaseNotificationRequest struct {
	AppAccountToken string `json:"appAccountToken" validate:"required"`
}
