package submodel

import (
	"gorm.io/gorm"
)

// AccountPhoneTable is used inside the repo layer
const (
	AccountPhoneTable = "sub_account_phones"
)

// AccountPhone model
type AccountPhone struct {
	gorm.Model
	AccountID uint `gorm:"not null;uniqueIndex:uniqueidx_account_phone" json:"account_id"`
	PhoneID   uint `gorm:"not null;uniqueIndex:uniqueidx_account_phone" json:"phone_id"`
	Default   byte `gorm:"default:0" json:"default"`
}
