package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	selfCheckRouter "payment-module/router/health-check"

	healthCheckRouter "payment-module/internals/router/health-check"
	notificationRouter "payment-module/internals/router/notification"
	paymentRecordRouter "payment-module/internals/router/payment-record"
	subscriptionRecordRouter "payment-module/internals/router/subscription-record"
	telegramRouter "payment-module/internals/router/telegram"
	vnpayRouter "payment-module/vnpay/router"
)

var BASE_PATH string = "/services/payment"

func SetupRoutes(app *fiber.App) {
	app.Get(BASE_PATH, func(c *fiber.Ctx) error {
		claimData := c.Locals("jwtClaims")
		if claimData == nil {
			return c.SendString("Jwt was bypassed")
		}
		return c.JSON(claimData)
	})

	// Group api calls with param '/tintuc/api'
	api := app.Group(BASE_PATH, logger.New())

	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	// Setup note routes, can use same syntax to add routes for more models
	selfCheckRouter.SetupRoutes(app)
	healthCheckRouter.SetupRoutes(v1)
	notificationRouter.SetupRoutes(v1)
	telegramRouter.SetupRoutes(v1)
	paymentRecordRouter.SetupRoutesV1(v1)
	vnpayRouter.SetupRoutesV1(v1)
	subscriptionRecordRouter.SetupRoutesV1(v1)
	// catch 404 and forward to error handler
	// app.use(function(req, res, next) {
	// 	next(createError(404));
	// });

}
