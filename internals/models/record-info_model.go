package models

import "time"

type RecordInfo struct {
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
	CreatedBy ForeignKey `json:"createdBy,omitempty" bson:"createdBy"`
	UpdatedBy ForeignKey `json:"updatedBy,omitempty" bson:"updatedBy"`
}
