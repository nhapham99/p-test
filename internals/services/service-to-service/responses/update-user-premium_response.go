package serviceToServiceResponses

import "time"

type ServiceUpdateUserPremiumResponse struct {
	OldPremiumFrom *time.Time `json:"oldPremiumFrom"`
	OldPremiumTo   *time.Time `json:"oldPremiumTo"`
	NewPremiumFrom *time.Time `json:"newPremiumFrom"`
	NewPremiumTo   *time.Time `json:"newPremiumTo"`
}
