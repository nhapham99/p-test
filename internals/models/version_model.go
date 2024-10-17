package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Version struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Version    int64              `json:"version,omitempty" validate:"required"`
	Interval   int                `json:"interval,omitempty" validate:"required"`
	CreatedAt  int64              `json:"createdAt,omitempty"  bson:"createdAt"`
	UpdateAt   int64              `json:"updateAt,omitempty" bson:"updateAt"`
	CreateUser string             `json:"createUser,omitempty" bson:"createUser"`
	UpdateUser string             `json:"updateUser,omitempty" bson:"updateUser"`
}
