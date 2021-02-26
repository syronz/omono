package eacmodel

import (
	"omono/internal/types"
)

// BalanceTable is used inside the repo layer
const (
	BalanceTable = "eac_balances"
)

// Balance model
type Balance struct {
	types.FixedCol
	AccountID  types.RowID `gorm:"not null;uniqueIndex:uniqueidx_account_balance" json:"account_id"`
	CurrencyID types.RowID `gorm:"not null;uniqueIndex:uniqueidx_account_balance" json:"currency_id"`
	Balance    float64     `gorm:"not null" json:"balance"`
}
