package services

import (
	"context"
	_error "payment-module/error"
	"payment-module/internals/constants"
	"payment-module/internals/models"
	"payment-module/internals/repositories"
	paymentRecordControllerRequests "payment-module/internals/requests/payment-record-controller_requests"
	"payment-module/internals/responses"
	paymentRecordControllerResponses "payment-module/internals/responses/payment-record-controller_responses"
	serviceToService "payment-module/internals/services/service-to-service"
	serviceToServiceRequests "payment-module/internals/services/service-to-service/requests"
	serviceResponse "payment-module/internals/services/service-to-service/responses"
	"payment-module/jwtchecker"
	"payment-module/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckPaymentMethod(c *fiber.Ctx,
	request paymentRecordControllerRequests.CreateNewPaymentRecordRequest) (*serviceResponse.ServiceGetPaymentMethodResponse,
	*serviceResponse.ServiceGetSubscriptionPackageResponse, *_error.SystemError) {
	paymentMethodInfo, err := serviceToService.GetAPaymentMethod(request.PaymentMethodId)
	if err != nil {
		return nil, nil, _error.New(err)
	}

	subscriptionPackageInfo, err := serviceToService.GetASubscriptionPackage(request.SubscriptionPackageId)
	if err != nil {
		return nil, nil, _error.New(err)
	}

	if !paymentMethodInfo.Active {
		return nil, nil, _error.NewErrorByString(_error.PAYMENT_E002_002)
	} else if !subscriptionPackageInfo.Active {
		return nil, nil, _error.NewErrorByString(_error.PAYMENT_E002_003)
	} else if subscriptionPackageInfo.Duration <= 0 {
		return nil, nil, _error.NewErrorByString(_error.PAYMENT_E002_004)
	} else if subscriptionPackageInfo.Price <= 0 {
		return nil, nil, _error.NewErrorByString(_error.PAYMENT_E002_005)
	}
	return paymentMethodInfo, subscriptionPackageInfo, nil
}

func CreateTransferPayment(c *fiber.Ctx, request paymentRecordControllerRequests.CreateNewPaymentRecordRequest,
	paymentMethodInfo *serviceResponse.ServiceGetPaymentMethodResponse,
	subscriptionPackageInfo *serviceResponse.ServiceGetSubscriptionPackageResponse) (*primitive.ObjectID, *_error.SystemError) {
	if request.BankAccount == "" || request.BankName == "" || request.Message == "" {
		return nil, _error.NewErrorByString(_error.PAYMENT_E001_006)
	}

	newObject := createPaymentRecord(c, request, paymentMethodInfo, subscriptionPackageInfo)

	createPaymentRecord, newPaymentRecord, err := repositories.InsertOnePaymentRecord(c, newObject)
	utils.UnUsed(createPaymentRecord, newPaymentRecord)
	if err != nil {
		return nil, _error.New(err)
	}

	return &newPaymentRecord.Id, nil
}

func AdminGetAllPaymentRecordsPagination(c *fiber.Ctx) (*responses.PaginationResponse, *_error.SystemError) {
	page, size, _, _, sortParams := utils.GetPaginationParams(c)
	filterStr := c.Query("filter", "")
	statusStr := c.Query("status", "")
	paymentMethodIdStr := c.Query("paymentMethodId", "")
	subscriptionPackageIdStr := c.Query("subscriptionPackageId", "")
	status := []int64{}
	paymentMethodId := []primitive.ObjectID{}
	subscriptionPackageId := []primitive.ObjectID{}

	if statusStr != "" {
		for _, v := range strings.Split(statusStr, ",") {
			intValue, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, _error.New(err)
			}
			status = append(status, intValue)
		}
	}

	if paymentMethodIdStr != "" {
		for _, v := range strings.Split(paymentMethodIdStr, ",") {
			if objectID, err := primitive.ObjectIDFromHex(v); err != nil {
				return nil, _error.New(err)
			} else {
				paymentMethodId = append(paymentMethodId, objectID)
			}
		}
	}

	if subscriptionPackageIdStr != "" {
		for _, v := range strings.Split(subscriptionPackageIdStr, ",") {
			if objectID, err := primitive.ObjectIDFromHex(v); err != nil {
				return nil, _error.New(err)
			} else {
				subscriptionPackageId = append(subscriptionPackageId, objectID)
			}
		}
	}

	paymentRecordCount, result, err := adminGetAllPaymentRecordsPagination(page, size, sortParams, filterStr, paymentMethodId, subscriptionPackageId, status)
	if err != nil {
		return nil, _error.New(err)
	}

	paymentRecords := []models.PaymentRecord{}
	if paymentRecordCount > 0 {
		err = result.All(context.TODO(), &paymentRecords)
		if err != nil {
			return nil, _error.New(err)
		}

		if len(paymentRecords) == 0 {
			paymentRecords = []models.PaymentRecord{}
		}
	}

	response := responses.PaginationResponse{
		RecordCount: paymentRecordCount,
		PageCount:   utils.CalculatePageCount(paymentRecordCount, size),
		CurrentPage: page,
		PageSize:    size,
		Records:     paymentRecords,
	}

	return &response, nil
}

func UserGetAPaymentRecord(c *fiber.Ctx) (*paymentRecordControllerResponses.UserGetPaymentRecordResponse, *_error.SystemError) {
	userId := jwtchecker.GetUserIdFromJwt(c)
	recordId := c.Params("objectId")
	paymentRecordId, err := primitive.ObjectIDFromHex(recordId)
	if err != nil {
		return nil, _error.New(err)
	}

	paymentRecord, err := findOnePaymentRecordById(paymentRecordId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, _error.New(err)
		}

		// Nếu ko tìm thấy bản ghi nào
		return nil, _error.NewErrorByString(_error.PAYMENT_E001_001)
	}

	if paymentRecord.User.Id != userId {
		return nil, _error.NewErrorByString(_error.PAYMENT_E001_001)
	}

	// Request hợp lệ: cập nhật trạng thái thành đã được user xác nhận chuyển khoản
	response := paymentRecordControllerResponses.UserGetPaymentRecordResponse{
		Status: paymentRecord.Status,
	}

	return &response, nil
}

func createPaymentRecord(c *fiber.Ctx, request paymentRecordControllerRequests.CreateNewPaymentRecordRequest,
	paymentMethodInfo *serviceResponse.ServiceGetPaymentMethodResponse,
	subscriptionPackageInfo *serviceResponse.ServiceGetSubscriptionPackageResponse) models.PaymentRecord {
	var bankInfo *models.BankInfo = nil

	if request.PaymentMethodId.Hex() == constants.PAYMENT_METHOD_TRANSFER_ID {
		bankInfo = &models.BankInfo{
			BankAccount: request.BankAccount,
			Message:     request.Message,
			BankName:    request.BankName,
		}
	}

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
		BankInfo: bankInfo,
	}
}

func AdmTransferConfirm(c *fiber.Ctx, recordId string) *_error.SystemError {
	paymentRecordId, err := primitive.ObjectIDFromHex(recordId)
	if err != nil {
		return _error.New(err)
	}

	paymentRecord, err := findOnePaymentRecordById(paymentRecordId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return _error.New(err)
		}

		// Nếu ko tìm thấy bản ghi nào
		return _error.NewErrorByString(_error.PAYMENT_E001_001)
	}

	// Nếu trạng thái payment không phải mới tạo
	if paymentRecord.Status != constants.PAYMENT_RECORD_STATUS_CREATED {
		return _error.NewErrorByString(_error.PAYMENT_E001_002)
	}

	// Nếu loại thanh toán không phải chuyển khoản
	if paymentRecord.PaymentMethod.Id.Hex() != constants.PAYMENT_METHOD_TRANSFER_ID {
		return _error.NewErrorByString(_error.PAYMENT_E001_003)
	}

	// Gửi cập nhật premium user sang auth service
	updateUserPremiumRequest := serviceToServiceRequests.ServiceUpdateUserPremiumRequest{
		UserId:    paymentRecord.User.Id,
		PaymentId: paymentRecord.Id,
		Duration:  paymentRecord.Duration,
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

	return nil
}

func findOnePaymentRecordById(id primitive.ObjectID) (*models.PaymentRecord, error) {
	filter := bson.M{"_id": id}
	result, err := repositories.FindOnePaymentRecord(filter)
	return result, err
}

func adminGetAllPaymentRecordsPagination(page int64, size int64, sortParams bson.D, filterString string, paymentMethodId []primitive.ObjectID, subscriptionPackageId []primitive.ObjectID, status []int64) (int64, *mongo.Cursor, error) {
	filter := bson.M{}

	if len(paymentMethodId) > 0 {
		filter["paymentMethod.id"] = bson.M{"$in": paymentMethodId}
	}

	if len(subscriptionPackageId) > 0 {
		filter["subscriptionPackage.id"] = bson.M{"$in": subscriptionPackageId}
	}

	if len(status) > 0 {
		filter["status"] = bson.M{"$in": status}
	}

	if filterString != "" {
		filter = bson.M{
			"$and": []bson.M{
				filter,
				{
					"$or": []bson.M{
						{
							"bankInfo.bankName": bson.M{
								"$regex": primitive.Regex{Pattern: filterString, Options: "i"},
							},
						},
						{
							"bankInfo.bankAccount": bson.M{
								"$regex": primitive.Regex{Pattern: filterString, Options: "i"},
							},
						},
						{
							"bankInfo.message": bson.M{
								"$regex": primitive.Regex{Pattern: filterString, Options: "i"},
							},
						},
						{
							"user.name": filterString,
						},
					},
				},
			},
		}
	}

	matchFilter := bson.M{
		"$match": filter,
	}

	if len(sortParams) == 0 {
		sortParams = bson.D{
			{"recordInfo.createdAt", -1},
			{"_id", 1},
		}
	}

	sort := bson.M{
		"$sort": sortParams,
	}

	limit := bson.M{
		"$limit": size,
	}

	skip := bson.M{
		"$skip": utils.CalculatePaginatedSkip(page, size),
	}

	pipeline := []bson.M{matchFilter, sort, skip, limit}

	count, err := repositories.CountPaymentRecords(filter)
	if err != nil {
		return 0, nil, err
	}

	if count == 0 {
		return 0, nil, nil
	}

	result, err := repositories.FindListPaymentRecordsWithPipeline(pipeline)
	if err != nil {
		return 0, nil, err
	}

	return count, result, err
}
