package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Code        string             `json:"code,omitempty" validate:"required"`
	Name        string             `json:"name,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Status      *int               `json:"status,omitempty" bson:"status"`
	Telegram    Telegram           `json:"telegram,omitempty" bson:"telegram"`
	ApiKey      string             `json:"apiKey,omitempty" bson:"apiKey" validate:"required"`
	Secretkey   string             `json:"secretkey,omitempty" bson:"secretkey" validate:"required"`
	CreatedAt   int64              `json:"createdAt,omitempty"  bson:"createdAt"`
	UpdateAt    int64              `json:"updateAt,omitempty" bson:"updateAt"`
	CreateUser  string             `json:"createUser,omitempty" bson:"createUser"`
	UpdateUser  string             `json:"updateUser,omitempty" bson:"updateUser"`
}
