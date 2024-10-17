package main

import (
	"fmt"
	"net/http"
	"os"
	"payment-module/configs"
	_ "payment-module/docs"
	"payment-module/internals/responses"
	"payment-module/internals/services"
	"payment-module/jwtchecker"
	"payment-module/logger"
	"payment-module/router"
	schedule "payment-module/schedule"
	configVnPay "payment-module/vnpay/config"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"golang.org/x/exp/slices"
)

var _currentVersion int64

// @title           Swagger Monitor API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /monitor/api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	logger.Init()
	logger.Infof("Starting on time:%s", time.Now().Format("2006-January-02"))

	// https://github.com/gofiber/fiber/issues/623
	// app := fiber.New(&fiber.Settings{
	// 	StrictRouting: true,
	// })
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendString(fmt.Sprintf("Http Code: %d", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})

	app.Use(cors.New())

	// Default middleware config
	// app.Use(logger.New(logger.Config{
	// 	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	// }))
	logger.Infof("get public key")
	key, err2 := jwtchecker.GetPublicKeyFromFile(configs.GetJwtPublicKeyPath())
	if err2 != nil {
		panic("This panic is caught by settings jwt-public-key-path fail")
	}
	app.Use(jwtchecker.New(jwtchecker.Config{
		Secret: key,
		Filter: func(c *fiber.Ctx) bool {
			// match, err := regexp.MatchString("/monitor/api/v1/*", c.Path())
			return slices.Contains(jwtchecker.PATH_UNCHECK_AUTHORIZATION, c.Path())
		},
	}))
	// app.Use(jwtchecker.New(jwtchecker.Config{
	// 	Filter: func(c *fiber.Ctx) bool {
	// 		// if X-Secret-Pass header matches then we skip the jwt validation
	// 		return c.Get("X-Secret-Pass") == "c5cacd9002a3"
	// 	},
	// }))

	logger.Infof("connect db start")
	//run database
	configs.ConnectDB()

	configVnPay.InitEnviroment(
		configs.GetVNPWhiteListIP(),
		configs.GetVNPUrl(),
		configs.GetVNPCreateOrderPath(),
		configs.GetVNPTMNCode(),
		configs.GetVNPOrderInfo(),
		configs.GetVNPReturnUrl(),
		configs.GetVNPHashSecret(),
	)

	logger.Infof("setup route start")
	// Setup the router
	router.SetupRoutes(app)

	app.Get("/services/payment/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(responses.BaseResponse{Status: http.StatusOK, Message: "Payment service DEMO", Data: time.Now().String()})
	})
	// startGetVersion()
	// startHealthCheck()

	// productPort := configs.GetProductPort()
	productPort := os.Getenv("PRODUCT_PORT")
	if productPort == "" {
		productPort = "8014"
	}
	err := app.Listen(":" + productPort)
	if err != nil {
		logger.Fatalf("fiber.Listen failed %s", err)
	}
}

func startGetVersion() {
	logger.Infof("Starting get version")
	_currentVersion = services.GetCurrentVersion()
	_interval := services.GetCurrentInterval()
	if _interval <= 0 {
		_interval = 5
	}
	logger.Infof("Starting with version:%o", _currentVersion)
	schedule.NewCheckVersionScheduler(_interval,
		func() {
			logger.Infof("Starting version check")
			version := services.GetCurrentVersion()
			logger.Infof("current-version: %o new-version:%o", _currentVersion, version)
			if version > _currentVersion {
				logger.Infof("########################### RESTART ALL SCHEDULE ###########################")
				restartHealthCheck()
				restartVersionSchedule()
				_currentVersion = version
			}
		},
	)
}

func restartVersionSchedule() {
	stopVersionSchedule()
	startGetVersion()
}

func stopVersionSchedule() {
	logger.Infof("Stop version-check")
	schedule.GetAllJobVersionCheck()
	schedule.ClearAllVersionCheckScheduler()
	schedule.GetAllJobVersionCheck()
}

func restartHealthCheck() {
	stopHealthCheck()
	startHealthCheck()
}

func stopHealthCheck() {
	logger.Infof("Stop health-check")
	schedule.GetAllJobHealthCheck()
	schedule.ClearAllHealthCheckScheduler()
	schedule.GetAllJobHealthCheck()
}

func startHealthCheck() {
	logger.Infof("Starting health-check")
	configs, err := services.GetAllConfig()
	if err != nil {
		logger.Fatalf("Get all config failed %s", err)
	} else {
		for i := 0; i < len(configs); i++ {
			logger.Infof("config:%s", configs[i])
			for j := 0; j < len(configs[i].Services); j++ {
				schedule.NewHealthCheckScheduler(configs[i].Code, configs[i].Services[j], configs[i].Telegram)
			}
		}
	}
}

func getGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}
