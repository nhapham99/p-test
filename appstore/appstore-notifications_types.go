package appstore

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AppStoreServerNotification struct {
	// appleRootCert   string
	Payload         *NotificationPayload
	TransactionInfo *TransactionInfo
	RenewalInfo     *RenewalInfo
	IsValid         bool
}

type AppStoreServerRequest struct {
	SignedPayload string `json:"signedPayload"`
}

type NotificationHeader struct {
	Alg string   `json:"alg"`
	X5c []string `json:"x5c"`
}

type NotificationPayload struct {
	Data             NotificationData `json:"data"`
	NotificationType string           `json:"notificationType"`
	NotificationUUID string           `json:"notificationUUID"`
	SignedDate       int              `json:"signedDate"`
	Subtype          string           `json:"subtype"`
	Version          string           `json:"version"`
	CreatedAt        time.Time        `json:"createdAt,omitempty"`
}

type NotificationSummary struct {
	RequestIdentifier      string   `json:"requestIdentifier"`
	AppAppleId             string   `json:"appAppleId"`
	BundleId               string   `json:"bundleId"`
	ProductId              string   `json:"productId"`
	Environment            string   `json:"environment"`
	StoreFrontCountryCodes []string `json:"storefrontCountryCodes"`
	FailedCount            int64    `json:"failedCount"`
	SucceededCount         int64    `json:"succeededCount"`
}

type NotificationData struct {
	AppAppleId            int    `json:"appAppleId"`
	BundleId              string `json:"bundleId"`
	BundleVersion         string `json:"bundleVersion"`
	Environment           string `json:"environment"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
	Status                int32  `json:"status"`
}

type TransactionInfo struct {
	AppAccountToken             string `json:"appAccountToken"`
	Currency                    string `json:"currency,omitempty"`
	BundleId                    string `json:"bundleId"`
	Environment                 string `json:"environment"`
	ExpiresDate                 int    `json:"expiresDate"`
	InAppOwnershipType          string `json:"inAppOwnershipType"`
	IsUpgraded                  bool   `json:"isUpgraded"`
	OfferIdentifier             string `json:"offerIdentifier"`
	OfferType                   int32  `json:"offerType"`
	OriginalPurchaseDate        int    `json:"originalPurchaseDate"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	Price                       int32  `json:"price,omitempty"`
	ProductId                   string `json:"productId"`
	PurchaseDate                int    `json:"purchaseDate"`
	Quantity                    int32  `json:"quantity"`
	RevocationDate              int    `json:"revocationDate"`
	RevocationReason            int32  `json:"revocationReason"`
	SignedDate                  int    `json:"signedDate"`
	StoreFront                  string `json:"storefront"`
	StoreFrontId                string `json:"storefrontId"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
	TransactionId               string `json:"transactionId"`
	TransactionReason           string `json:"transactionReason"`
	Type                        string `json:"type"`
	WebOrderLineItemId          string `json:"webOrderLineItemId"`
}

type RenewalInfo struct {
	jwt.StandardClaims
	AutoRenewProductId          string `json:"autoRenewProductId"`
	AutoRenewStatus             int32  `json:"autoRenewStatus"`
	Environment                 string `json:"environment"`
	ExpirationIntent            int32  `json:"expirationIntent"`
	GracePeriodExpiresDate      int    `json:"gracePeriodExpiresDate"`
	IsInBillingRetryPeriod      bool   `json:"isInBillingRetryPeriod"`
	OfferIdentifier             string `json:"offerIdentifier"`
	OfferType                   int32  `json:"offerType"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	PriceIncreaseStatus         int32  `json:"priceIncreaseStatus"`
	ProductId                   string `json:"productId"`
	RecentSubscriptionStartDate int    `json:"recentSubscriptionStartDate"`
	RenewalDate                 int    `json:"renewalDate"`
	SignedDate                  int    `json:"signedDate"`
}
