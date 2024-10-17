package services

import (
	_error "payment-module/error"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	"payment-module/internals/repositories"
	serviceToService "payment-module/internals/services/service-to-service"
	serviceToServiceRequests "payment-module/internals/services/service-to-service/requests"
	"payment-module/jwtchecker"
	"payment-module/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	paymentRecordControllerRequests "payment-module/internals/requests/payment-record-controller_requests"
	serviceResponse "payment-module/internals/services/service-to-service/responses"
)

func FindOnePaymentRecordById(id primitive.ObjectID) (*models.PaymentRecord, error) {
	filter := bson.M{"_id": id}
	result, err := repositories.FindOnePaymentRecord(filter)
	return result, err
}

// Tao moi PaymentRecord de luu vao mongodb
func createPaymentRecord(c *fiber.Ctx, request paymentRecordControllerRequests.CreateNewPaymentRecordRequest,
	paymentMethodInfo *serviceResponse.ServiceGetPaymentMethodResponse,
	subscriptionPackageInfo *serviceResponse.ServiceGetSubscriptionPackageResponse) models.PaymentRecord {
	return models.PaymentRecord{
		User: models.ForeignKey{
			Id:   jwtchecker.GetUserIdFromJwt(c),
			Name: jwtchecker.GetUserNameFromJwt(c),
		},
		PaymentMethod: models.ForeignKey{
			Id:   request.PaymentMethodId,
			Name: paymentMethodInfo.Name,
		},
		SubscriptionPackage: models.ForeignKey{
			Id:   request.SubscriptionPackageId,
			Name: subscriptionPackageInfo.Name,
		},
		Price:    subscriptionPackageInfo.Price,
		Duration: subscriptionPackageInfo.Duration,
		Status:   constants.PAYMENT_RECORD_STATUS_CREATED,
	}
}

// Xu ly yeu cau thanh toan & luu vao DB
func CreatePaymentRecord(c *fiber.Ctx, request paymentRecordControllerRequests.CreateNewPaymentRecordRequest,
	paymentMethodInfo *serviceResponse.ServiceGetPaymentMethodResponse,
	subscriptionPackageInfo *serviceResponse.ServiceGetSubscriptionPackageResponse) (*string, *_error.SystemError) {
	newObject := createPaymentRecord(c, request, paymentMethodInfo, subscriptionPackageInfo)

	createPaymentRecord, newPaymentRecord, err := repositories.InsertOnePaymentRecord(c, newObject)
	utils.UnUsed(createPaymentRecord, newPaymentRecord)
	if err != nil {
		return nil, _error.New(err)
	}

	paymentRecordIdString := newPaymentRecord.Id.Hex()
	return &paymentRecordIdString, nil
}

// Xac nhan thanh toan thanh cong
func AdmTransferConfirm(c *fiber.Ctx, recordId string) *_error.SystemError {
	paymentRecordId, err := primitive.ObjectIDFromHex(recordId)
	if err != nil {
		return _error.New(err)
	}

	paymentRecord, err := FindOnePaymentRecordById(paymentRecordId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return _error.New(err)
		}

		// Nếu ko tìm thấy bản ghi nào
		return _error.NewErrorByString(_error.PAYMENT_E001_001)
	}

	// Nếu loại thanh toán chuyển khoản & trạng thái payment không phải mới tạo
	if paymentRecord.PaymentMethod.Id.Hex() == constants.PAYMENT_METHOD_TRANSFER_ID && paymentRecord.Status != constants.PAYMENT_RECORD_STATUS_CREATED {
		return _error.NewErrorByString(_error.PAYMENT_E001_002)
	}

	var platform int
	if paymentRecord.PaymentMethod.Id.Hex() == constants.PAYMENT_METHOD_TRANSFER_ID {
		platform = constants.PLATFORMS.DirectTransfer
	}

	if paymentRecord.PaymentMethod.Id.Hex() == constants.PAYMENT_METHOD_VNPAY_ID {
		platform = constants.PLATFORMS.Vnpay
	}

	// Gửi cập nhật premium user sang auth service
	updateUserPremiumRequest := serviceToServiceRequests.ServiceUpdateUserPremiumRequest{
		UserId:    paymentRecord.User.Id,
		PaymentId: paymentRecord.Id,
		Duration:  paymentRecord.Duration,
		Platform:  platform,
	}

	updateUserPremiumResponse, err := serviceToService.UpdateUserPremium(updateUserPremiumRequest)
	if err != nil {
		return _error.New(err)
	}

	// Cập nhật trạng thái thành đã được xác nhận thanh toán, cập nhật thời hạn premium vào bản ghi
	filter := bson.M{"_id": paymentRecordId}
	update := bson.M{
		"status":                  constants.PAYMENT_RECORD_STATUS_PAID,
		"paymentConfirmationDate": time.Now(),
		"oldPremiumFrom":          updateUserPremiumResponse.OldPremiumFrom,
		"oldPremiumTo":            updateUserPremiumResponse.OldPremiumTo,
		"newPremiumFrom":          updateUserPremiumResponse.NewPremiumFrom,
		"newPremiumTo":            updateUserPremiumResponse.NewPremiumTo,
	}

	updatePaymentRecord, err := repositories.UpdatePaymentRecords(c, filter, update)
	utils.UnUsed(updatePaymentRecord)

	if err != nil {
		return _error.New(err)
	}

	// gửi record sang tracking service, premiumFrom là hiện tại, premiumTo là hiện tại + duration
	err = serviceToService.AddPaymentToTracking(serviceToServiceRequests.AddPaymentToTrackingRequest{
		User: serviceToServiceRequests.ForeignKey{
			Id:   paymentRecord.User.Id.Hex(),
			Name: paymentRecord.User.Name,
		},
		PaymentId: paymentRecord.Id.Hex(),
		PaymentMethod: serviceToServiceRequests.ForeignKey{
			Id:   paymentRecord.PaymentMethod.Id.Hex(),
			Name: paymentRecord.PaymentMethod.Name,
		},
		PremiumFrom: time.Now(),
		PremiumTo:   time.Now().AddDate(0, 0, int(paymentRecord.Duration)),
	})

	if err != nil {
		return _error.New(err)
	}

	return nil
}
