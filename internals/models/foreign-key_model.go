package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ForeignKey struct {
	Id      primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name"`
	Version string             `json:"version,omitempty" bson:"version"`
}
