package models

type Telegram struct {
	Token       string `json:"token,omitempty" validate:"required"`
	ChatId      string `json:"chatId,omitempty"  bson:"chatId" validate:"required"`
	Description string `json:"description,omitempty"`
}
