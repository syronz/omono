package locmodel

import (
	"omono/internal/types"
)

// StoreUserTable is a global instance for working with store_user
const (
	StoreUserTable = "loc_store_users"
)

// StoreUser model
// it is mostly controlled inside the store
type StoreUser struct {
	types.FixedCol
	StoreID types.RowID `gorm:"not null;uniqueIndex:uniqueidx_user_store" json:"store_id"`
	UserID  types.RowID `gorm:"not null;uniqueIndex:uniqueidx_user_store" json:"user_id"`
}
