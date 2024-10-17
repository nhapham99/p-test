package models

type BankInfo struct {
	BankName    string `json:"bankName,omitempty" bson:"bankName"`
	BankAccount string `json:"bankAccount,omitempty" bson:"bankAccount"`
	Message     string `json:"message,omitempty" bson:"message"`
}
