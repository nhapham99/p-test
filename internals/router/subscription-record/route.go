package SubscriptionRecordRouter

import (
	"payment-module/internals/controllers"
	"payment-module/internals/secure"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutesV1(router fiber.Router) {
	userCollection := router.Group("/user/subscription-record")

	// Android
	userCollection.Post("/android-client-purchase-notification", secure.CheckUserRolePermission(secure.ROLE_USER), controllers.AndroidClientPurchaseNotification)
	userCollection.Get("/check-google-play-renew", secure.CheckUserRolePermission(secure.ROLE_USER), controllers.CheckGooglePlayRenew)

	// IOS
	userCollection.Post("/ios-client-purchase-notification", secure.CheckUserRolePermission(secure.ROLE_USER), controllers.IOSClientPurchaseNotification)
	userCollection.Post("/app-store-server-notification", controllers.AppStoreServerNotification)
}
