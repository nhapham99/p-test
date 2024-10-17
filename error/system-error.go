package errors

import (
	"errors"
	paymentError "payment-module/vnpay/errors"
)

type SystemError struct {
	errorCode int
	err       error
}

func (e *SystemError) Error() error {
	return e.err
}

func (e *SystemError) ErrorMessage() string {
	if e == nil || e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *SystemError) ErrorCode() int {
	return e.errorCode
}

func NewOrderNotFound() *SystemError {
	return &SystemError{
		errorCode: paymentError.ORDER_NOT_FOUND,
		err:       errors.New(PAYMENT_E004_001),
	}
}

func NewOrderalreadyConfirmed() *SystemError {
	return &SystemError{
		errorCode: paymentError.ORDER_ALREADY_CONFIRMED,
		err:       errors.New(PAYMENT_E004_002),
	}
}

func NewInvalidAmount() *SystemError {
	return &SystemError{
		errorCode: paymentError.INVALID_AMOUNT,
		err:       errors.New(PAYMENT_E004_003),
	}
}

func NewChecksumInvalid() *SystemError {
	return &SystemError{
		errorCode: paymentError.INVALID_CHECKSUM,
		err:       errors.New(PAYMENT_E099_002),
	}
}

func NewInvalidIP(ip string) *SystemError {
	return &SystemError{
		errorCode: paymentError.INVALID_IP,
		err:       errors.New(PAYMENT_E004_004 + ip),
	}
}

func NewHandlerErrorWithMessage(message string) *SystemError {
	return &SystemError{
		errorCode: HANDLER_ERROR_WITH_MESSAGE,
		err:       errors.New(message),
	}
}

func New(err error) *SystemError {
	return &SystemError{
		errorCode: 0,
		err:       err,
	}
}

func NewErrorByString(err string) *SystemError {
	return &SystemError{
		errorCode: 0,
		err:       errors.New(err),
	}
}
