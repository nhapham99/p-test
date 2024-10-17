package configs

import (
	"os"
	"payment-module/logger"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGOURI")
}

func GetEnv(envName string) string {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env name:" + envName)
	}

	return os.Getenv(envName)
}

func GetEnvFromOS(envName string) string {
	result := os.Getenv(envName)
	logger.Debugf("GetEnvFromOS with name:%s result:%s", envName, result)
	if len(result) <= 0 {
		return GetEnv(envName)
	}
	return result
}

func GetInt64Env(envName string) int64 {
	val := GetEnv(envName)
	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logger.Fatal("Error loading .env name:" + envName)
	}
	return ret
}

// -------------- Start JWT ----------------------------
func GetJwtPublicKeyPath() string {
	return GetEnv("JWT_PUBLIC_KEY_PATH")
}

// -------------- End JWT ----------------------------

// Start get Product Common Config
func GetProductPort() string {
	return GetEnv("PRODUCT_PORT")
}

// End get Product Common Config

// Start get error message
func GetSystemErrorMessage() string {
	return GetEnv("ERROR_SYSTEM")
}

func GetDBNotFoundMessage() string {
	return GetEnv("ERROR_DATABASE_NOT_FOUND")
}

// End get error message

// Start get Database Config
func GetConfigDatabaseName() string {
	return GetEnv("CONFIG_DATABASE_NAME")
}

func GetConfigCollection() string {
	return "health-check-config"
}

func GetNotificationConfigCollection() string {
	return "notification-config"
}

func GetVersionCollection() string {
	return "version"
}

func GetPaymentRecordCollection() string {
	return "payment-record"
}

func GetSubscriptionRecordCollection() string {
	return "subscription-record"
}

func GetSubscriptionRecordRequestCollection() string {
	return "subscription-request-record"
}

func GetAppStoreServerNotificationCollection() string {
	return "appstore-server-notification"
}

func GetTradingCollection() string {
	return "trading"
}

func GetTradedCollection() string {
	return "traded"
}

func GetVNPAYIPNCollection() string {
	return "vnpay_ipn"
}

// End get Database Config

// Start get service to service config

func GetServiceCategoryBaseUrl() string {
	return GetEnvFromOS("SERVICE_CATEGORY_BASEURL")
}

func GetServiceAuthBaseUrl() string {
	return GetEnvFromOS("SERVICE_AUTH_BASEURL")
}

func getServiceTrackingBaseUrl() string {
	return GetEnvFromOS("SERVICE_TRACKING_BASEURL")
}

func GetGetAPaymentMethodUrl() string {
	return GetServiceCategoryBaseUrl() + GetEnvFromOS("SERVICE_CATEGORY_GETAPAYMENTMETHOD_URL")
}

func GetGetASubscriptionPackageUrl() string {
	return GetServiceCategoryBaseUrl() + GetEnvFromOS("SERVICE_CATEGORY_GETASUBSCRIPTIONPACKAGE_URL")
}

func GetUpdateUserPremiumUrl() string {
	return GetServiceAuthBaseUrl() + GetEnvFromOS("SERVICE_AUTH_UPDATEUSERPAYMENT_URL")
}

func GetAddPaymentToTrackingUrl() string {
	return getServiceTrackingBaseUrl() + GetEnvFromOS("SERVICE_TRACKING_ADDPAYMENTTOTRACKING_URL")
}

func GetGetUserInfoUrl() string {
	return GetServiceAuthBaseUrl() + "/services/auth-internal/account/get-user-info"
}

func GetUpdateUserSubscriptionUrl() string {
	return GetServiceAuthBaseUrl() + "/services/auth-internal/account/update-user-subscription"
}

// End get service to service config

// Start get VNPAY enviroment
func GetVNPWhiteListIP() string {
	return GetEnvFromOS("VNP_WHITELIST_IP")
}
func GetVNPUrl() string {
	return GetEnvFromOS("VNP_URL")
}
func GetVNPCreateOrderPath() string {
	return GetEnvFromOS("VNP_CREATE_ORDER_PATH")
}
func GetVNPTMNCode() string {
	return GetEnvFromOS("VNP_TMNCODE")
}
func GetVNPOrderInfo() string {
	return GetEnvFromOS("VNP_ORDERINFO")
}
func GetVNPReturnUrl() string {
	return GetEnvFromOS("VNP_RETURNURL")
}
func GetVNPHashSecret() string {
	return GetEnvFromOS("VNP_HASHSECRET")
}

// End get VNPAY enviroment
