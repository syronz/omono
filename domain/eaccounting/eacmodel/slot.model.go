package eacmodel

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"time"
)

// SlotTable is a global instance for working with slot
const (
	SlotTable = "eac_slots"
)

// Slot model
type Slot struct {
	types.FixedNode
	CurrencyID    types.RowID `gorm:"not null" json:"currency_id,omitempty"`
	AccountID     types.RowID `gorm:"not null" json:"account_id,omitempty"`
	TransactionID types.RowID `gorm:"not null" json:"transaction_id,omitempty"`
	Debit         float64     `json:"debit,omitempty"`
	Credit        float64     `json:"credit,omitempty"`
	Balance       float64     `json:"balance,omitempty"`
	Description   *string     `json:"description,omitempty"`
	PostDate      time.Time   `json:"post_date,omitempty" table:"eac_slots.post_date"`
	Rows          int         `json:"rows,omitempty"`
}

//DetailedSlots include account name, code, and currency in detail
type DetailedSlots struct {
	types.FixedNode
	CurrencyID    types.RowID `json:"currency_id,omitempty"`
	CurrencyName  string      `json:"currency_name,omitempty"`
	AccountID     types.RowID `json:"account_id,omitempty"`
	AccountName   string      `json:"account_name,omitempty"`
	AccountCode   string      `json:"account_code,omitempty"`
	TransactionID types.RowID `json:"transaction_id,omitempty"`
	Debit         float64     `json:"debit,omitempty"`
	Credit        float64     `json:"credit,omitempty"`
	Balance       float64     `json:"balance,omitempty"`
	Description   *string     `json:"description,omitempty"`
	PostDate      time.Time   `json:"post_date,omitempty"`
	Rows          int         `json:"rows,omitempty"`
}

// 	debit	credit	balance	description	post_date	row

// Validate check the type of fields
func (p *Slot) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:
		if p.Description != nil {
			if len(*p.Description) > 255 {
				err = limberr.AddInvalidParam(err, "description",
					corerr.MaximumAcceptedCharacterForVisV,
					dict.R(corterm.Description), 255)
			}
		}
	}

	return err
}
