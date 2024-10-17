package jwtchecker

// Path that don't check Authorization
var PATH_UNCHECK_AUTHORIZATION = []string{
	"/",
	"/health-check",
	"/monitor/api/v1/telegram/sendMessage",
	"/services/payment/v1/vnpay/ipn",
	"/services/payment/v1/user/subscription-record/app-store-server-notification",
}

// things := []string{"foo", "bar", "baz"}
