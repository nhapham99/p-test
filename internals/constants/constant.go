package constants

import (
	"time"
)

const AUTHORIZATION string = "Authorization"
const APPLICATION_JSON string = "application/json"

const TIME_OUT_EXECUTE time.Duration = 10    //sec
const TIME_OUT_CONNECTION time.Duration = 10 //sec

const HomeDocumentId string = "5c7452c7aeb4c97e0c2738"
const LimitChuongTrinhHot int = -1
const LimitGoiKham int = 2
const LimitTinTuc int = 5

// Định dạng date string
const YYYYMMDDHHMISS string = "20060102150405"
const DD_slash_MM_slash_YYYY string = "02/01/2006"

// JWT constant
const JWT_KEY_USER_ID string = "userId"
const JWT_KEY_USER_NAME string = "sub"
const JWT_KEY_PHONE string = "phone"
const JWT_KEY_EMAIL_VALIDATED string = "emailValidated"
const JWT_KEY_FIRST_LOGIN string = "firstLogin"
const JWT_TOKEN_TYPE string = "bearer"
const RJWT_KEY_ROLE string = "roles"

const HOSPITAL_ID_JWT_KEY string = "hisId"
const HOSPITAL_CODE_JWT_KEY string = "hospitalCode"
const USER_ID_JWT_KEY string = "userID"

// DB Constant
const CATEGORY_COLLECTION_NAME string = "category"
const NEWS_COLLECTION_NAME string = "news"
const HOME_COLLECTION_NAME string = "home"

// Row status
const ROW_STATUS_DISABLE int = 0
const ROW_STATUS_ENABLE int = 1

// Warning Type
const WARNING_TYPE_TELEGRAM = 1

// OTP status
const OTP_LENGTH = 6
const OTP_MAX_CREATE_COUNT = 3
const OTP_MAX_INPUT_COUNT = 3
const OTP_CREATE_WAIT = 60 // giây
const OTP_LIFE_TIME = 300  // giây (5 phút)
const OTP_STATUS_CREATE = 0
const OTP_STATUS_CONFIRMED = 1
const OTP_STATUS_USED = 2
const RESET_PWD_MIN_WAIT = 30            // ngày
const OTP_RESET_PWD_EMAIL_LIFE_TIME = 24 // giờ
const RESET_PWD_EMAIL_PATH string = "%s/password-service/%s/%s"

// PaymentRecord status
const PAYMENT_RECORD_STATUS_CREATED = 0
const PAYMENT_RECORD_STATUS_PAID = 1
const PAYMENT_RECORD_STATUS_TRANSFER_REJECTED = 2
const PAYMENT_RECORD_STATUS_PAY_ERROR = 3

// PaymentType
var BASE_PLANS = map[string]int64{
	"subscription1month":  30,
	"subscription6months": 180,
	"subscription1year":   365,
}

type Platforms struct {
	Others         int
	Googleplay     int
	Appstore       int
	DirectTransfer int
	Combo          int
	Voucher        int
}

var PLATFORMS = Platforms{
	Others:         -1,
	Googleplay:     1,
	Appstore:       2,
	DirectTransfer: 3,
	Combo:          4,
	Voucher:        5,
}

var APPSTORE_HANDLED_TYPES = []string{"SUBSCRIBED", "DID_RENEW", "DID_CHANGE_RENEWAL_PREF"}
