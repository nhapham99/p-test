package models

type Service struct {
	Code          string `json:"code,omitempty" validate:"required"`
	Name          string `json:"name,omitempty" validate:"required"`
	URL           string `json:"url,omitempty" validate:"required"`
	Method        string `json:"method,omitempty" validate:"required"`
	HttpCode      int    `json:"httpCode,omitempty" bson:"httpCode" validate:"required"`
	Description   string `json:"description,omitempty" validate:"required"`
	Status        *int   `json:"status,omitempty" bson:"status"`
	Interval      int    `json:"interval,omitempty" bson:"interval" validate:"required"`
	MessageError  string `json:"messageError,omitempty" bson:"messageError" validate:"required"`
	WarningType   int    `json:"warningType,omitempty" bson:"warningType" validate:"required"`
	NotifySuccess int    `json:"notifySuccess,omitempty" bson:"notifySuccess" validate:"required"`
}
