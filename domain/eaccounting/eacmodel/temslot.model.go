package eacmodel

import (
	"github.com/syronz/limberr"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/dict"
	"time"
)

// TempSlotTable is a global instance for working with slot
const (
	TempSlotTable = "eac_temp_slots"
)

// TempSlot model
type TempSlot struct {
	types.FixedNode
	CurrencyID    types.RowID `gorm:"not null" json:"currency_id,omitempty"`
	AccountID     types.RowID `gorm:"not null" json:"account_id,omitempty"`
	TransactionID types.RowID `gorm:"not null" json:"transaction_id,omitempty"`
	Debit         float64     `json:"debit,omitempty"`
	Credit        float64     `json:"credit,omitempty"`
	Balance       float64     `json:"balance,omitempty"`
	Description   *string     `json:"description,omitempty"`
	PostDate      time.Time   `json:"post_date,omitempty" table:"eac_temp_slots.post_date"`
	Rows          int         `json:"rows,omitempty"`
}

// 	debit	credit	balance	description	post_date	row

// ValidateCheckType check the type of fields
func (p *TempSlot) ValidateCheckType(act coract.Action) (err error) {

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
