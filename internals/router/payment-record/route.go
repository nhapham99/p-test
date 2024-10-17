package PaymentRecordRouter

import (
	"payment-module/internals/controllers"
	"payment-module/internals/secure"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutesV1(router fiber.Router) {
	userCollection := router.Group("/user/payment-record")
	userCollection.Put("/create-payment-record", secure.CheckUserRolePermission(secure.ROLE_USER), controllers.CreatePaymentRecord)
	userCollection.Get("/get-a-payment-record/:objectId", secure.CheckUserRolePermission(secure.ROLE_USER), controllers.UserGetAPaymentRecord)

	adminCollection := router.Group("/admin/payment-record")
	adminCollection.Get("/get-list-payment-records", secure.CheckUserRolePermission(secure.ROLE_SITE_ADMIN), controllers.AdminGetAllPaymentRecordsPagination)
	adminCollection.Post("/adm-transfer-confirm/:objectId", secure.CheckUserRolePermission(secure.ROLE_SITE_ADMIN), controllers.AdmTransferConfirm)
}
