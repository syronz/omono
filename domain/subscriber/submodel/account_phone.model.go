package submodel

import (
	"omono/internal/types"
)

// AccountPhoneTable is used inside the repo layer
const (
	AccountPhoneTable = "sub_account_phones"
)

// AccountPhone model
type AccountPhone struct {
	types.FixedCol
	AccountID types.RowID `gorm:"not null;uniqueIndex:uniqueidx_account_phone" json:"account_id"`
	PhoneID   types.RowID `gorm:"not null;uniqueIndex:uniqueidx_account_phone" json:"phone_id"`
	Default   byte        `gorm:"default:0" json:"default"`
}
