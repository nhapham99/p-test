package models

type Message struct {
	Message  string `json:"message,omitempty" validate:"required"`
	Time     int64  `json:"time,omitempty"  validate:"required"`
	Checksum string `json:"checksum,omitempty" validate:"required"`
}
