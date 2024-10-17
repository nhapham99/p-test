package serviceToService

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"payment-module/configs"
	_error "payment-module/error"
	"payment-module/internals/constants"
	"payment-module/internals/responses"
	serviceToServiceRequests "payment-module/internals/services/service-to-service/requests"
	serviceToServiceResponses "payment-module/internals/services/service-to-service/responses"
	"payment-module/logger"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAPaymentMethod(paymentMethodId primitive.ObjectID) (*serviceToServiceResponses.ServiceGetPaymentMethodResponse, error) {
	url := configs.GetGetAPaymentMethodUrl() + "/" + paymentMethodId.Hex()
	logger.Infof("GetAPaymentMethod with url:%s", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		errMessage := fmt.Sprintf(_error.PAYMENT_E002_001, _error.PAYMENT_E002_CATEGORY_SERVICE_NAME, strconv.Itoa(response.StatusCode))
		return nil, errors.New(errMessage)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData responses.BaseResponse
	json.Unmarshal(data, &responseData)

	if response.StatusCode == http.StatusBadRequest {
		return nil, errors.New(_error.PAYMENT_E002_CATEGORY_SERVICE_NAME + ": " + responseData.Message)
	}

	inputBytes, _ := json.Marshal(responseData.Data)
	var result serviceToServiceResponses.ServiceGetPaymentMethodResponse
	if err := json.Unmarshal(inputBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetASubscriptionPackage(subscriptionPackageId primitive.ObjectID) (*serviceToServiceResponses.ServiceGetSubscriptionPackageResponse, error) {
	url := configs.GetGetASubscriptionPackageUrl() + "/" + subscriptionPackageId.Hex()
	logger.Infof("GetASubscriptionPackage with url:%s", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		errMessage := fmt.Sprintf(_error.PAYMENT_E002_001, _error.PAYMENT_E002_CATEGORY_SERVICE_NAME, strconv.Itoa(response.StatusCode))
		return nil, errors.New(errMessage)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData responses.BaseResponse
	json.Unmarshal(data, &responseData)

	if response.StatusCode == http.StatusBadRequest {
		return nil, errors.New(_error.PAYMENT_E002_CATEGORY_SERVICE_NAME + ": " + responseData.Message)
	}

	inputBytes, _ := json.Marshal(responseData.Data)
	var result serviceToServiceResponses.ServiceGetSubscriptionPackageResponse
	if err := json.Unmarshal(inputBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateUserPremium(request serviceToServiceRequests.ServiceUpdateUserPremiumRequest) (*serviceToServiceResponses.ServiceUpdateUserPremiumResponse, error) {
	url := configs.GetUpdateUserPremiumUrl()
	logger.Infof("UpdateUserPremium with url:%s", url)

	requestDataJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(url, constants.APPLICATION_JSON, bytes.NewBuffer(requestDataJson))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		errMessage := fmt.Sprintf(_error.PAYMENT_E002_001, _error.PAYMENT_E002_CATEGORY_SERVICE_NAME, strconv.Itoa(response.StatusCode))
		return nil, errors.New(errMessage)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData responses.BaseResponse
	json.Unmarshal(data, &responseData)

	if response.StatusCode == http.StatusBadRequest {
		return nil, errors.New(_error.PAYMENT_E002_CATEGORY_SERVICE_NAME + ": " + responseData.Message)
	}

	inputBytes, _ := json.Marshal(responseData.Data)
	var result serviceToServiceResponses.ServiceUpdateUserPremiumResponse
	if err := json.Unmarshal(inputBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func AddPaymentToTracking(request serviceToServiceRequests.AddPaymentToTrackingRequest) error {
	url := configs.GetAddPaymentToTrackingUrl()
	logger.Infof("AddPaymentToTracking with url:%s", url)

	requestDataJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := http.Post(url, constants.APPLICATION_JSON, bytes.NewBuffer(requestDataJson))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		errMessage := fmt.Sprintf(_error.PAYMENT_E002_001, _error.PAYMENT_E002_CATEGORY_SERVICE_NAME, strconv.Itoa(response.StatusCode))
		return errors.New(errMessage)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var responseData responses.BaseResponse
	json.Unmarshal(data, &responseData)

	if response.StatusCode == http.StatusBadRequest {
		return errors.New(_error.PAYMENT_E002_CATEGORY_SERVICE_NAME + ": " + responseData.Message)
	}

	return nil
}

func GetUserInfo(userId primitive.ObjectID) (*serviceToServiceResponses.ServiceGetUserInfoResponse, error) {
	url := configs.GetGetUserInfoUrl() + "/" + userId.Hex()
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		errMessage := fmt.Sprintf(_error.PAYMENT_E002_001, _error.PAYMENT_E001_AUTH_SERVICE_NAME, strconv.Itoa(response.StatusCode))
		return nil, errors.New(errMessage)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData responses.BaseResponse
	json.Unmarshal(data, &responseData)

	if response.StatusCode == http.StatusBadRequest {
		return nil, errors.New(_error.PAYMENT_E001_AUTH_SERVICE_NAME + ": " + responseData.Message)
	}

	inputBytes, _ := json.Marshal(responseData.Data)
	var result serviceToServiceResponses.ServiceGetUserInfoResponse
	if err := json.Unmarshal(inputBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateUserSubscription(userId primitive.ObjectID, premiumFrom, premiumTo time.Time, lastPurchaseToken *string, platform int) (*serviceToServiceResponses.ServiceUpdateUserPremiumResponse, error) {
	url := configs.GetUpdateUserSubscriptionUrl()
	logger.Infof("UpdateUserSubscription with url:%s", url, userId, premiumFrom, premiumTo)

	// create body request
	requestData := serviceToServiceRequests.ServiceUpdateUserSubscriptionRequest{
		UserId:            userId,
		PremiumFrom:       premiumFrom,
		PremiumTo:         premiumTo,
		LastPurchaseToken: lastPurchaseToken,
		Platform:          platform,
	}
	requestDataJson, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	// send request
	response, err := http.Post(url, constants.APPLICATION_JSON, bytes.NewBuffer(requestDataJson))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		errMessage := fmt.Sprintf(_error.PAYMENT_E002_001, _error.PAYMENT_E002_CATEGORY_SERVICE_NAME, strconv.Itoa(response.StatusCode))
		return nil, errors.New(errMessage)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData responses.BaseResponse
	json.Unmarshal(data, &responseData)

	if response.StatusCode == http.StatusBadRequest {
		return nil, errors.New(_error.PAYMENT_E002_CATEGORY_SERVICE_NAME + ": " + responseData.Message)
	}

	inputBytes, _ := json.Marshal(responseData.Data)
	var result serviceToServiceResponses.ServiceUpdateUserPremiumResponse
	if err := json.Unmarshal(inputBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
