package eacmodel

import (
	"errors"
	"github.com/syronz/limberr"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/eaccounting/enum/transactionstatus"
	"omono/domain/eaccounting/enum/transactiontype"
	"omono/internal/consts"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
	"time"
)

// TransactionTable is a global instance for working with transaction
const (
	TransactionTable = "eac_transactions"
)

// Transaction model
type Transaction struct {
	types.FixedNode
	CurrencyID     types.RowID  `gorm:"not null" json:"currency_id,omitempty"`
	Rate           float64      `gorm:"not null" json:"rate,omitempty"`
	CreatedBy      types.RowID  `gorm:"not null" json:"created_by,omitempty"`
	Hash           string       `gorm:"not null;unique;type:varchar(100)" json:"hash,omitempty"`
	Type           types.Enum   `json:"type,omitempty"`
	Description    *string      `json:"description,omitempty"`
	Amount         float64      `json:"amount,omitempty"`
	YearCounter    uint64       `gorm:"index:idx_year_counter" json:"year_counter,omitempty"`
	YearCumulative uint64       `gorm:"index:idx_year_cumulative" json:"year_cumulative,omitempty"`
	Invoice        string       `json:"invoice,omitempty"`
	Pioneer        types.RowID  `gorm:"-" json:"pioneer,omitempty" table:"-"`
	Follower       types.RowID  `gorm:"-" json:"follower,omitempty" table:"-"`
	PostDate       time.Time    `gorm:"index:idx_post_date" json:"post_date,omitempty"`
	Status         types.Enum   `json:"status,omitempty"`
	Slots          []Slot       `gorm:"-" json:"slots,omitempty" table:"-"`
	Err            error        `gorm:"-"`
	Before         *Transaction `gorm:"-"`
	UserID         types.RowID  `json:"user_id,omitempty" table:"-"`
	Message        string       `json:"message,omitempty" table:"-"`
}

// TransactionCh is used for send all transactions request via one worker
type TransactionCh struct {
	Transaction Transaction
	Type        types.Enum
	Respond     chan Transaction
}

// Validate check the type of fields
func (p *Transaction) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:
		{

			err = validatePostDate(err, p.PostDate)
			err = validateSlots(err, p.Slots)
			err = validateCurrencyID(err, p.CurrencyID)
			err = validateRate(err, p.Rate)

			if p.Description != nil {
				if len(*p.Description) > 255 {
					err = limberr.AddInvalidParam(err, "description",
						corerr.MaximumAcceptedCharacterForVisV,
						dict.R(corterm.Description), 255)
				}
			}

			if ok, _ := helper.Includes(transactiontype.List, p.Type); !ok {
				return limberr.AddInvalidParam(err, "type",
					corerr.AcceptedValueForVareV, dict.R(corterm.Type),
					transactiontype.Join())
			}

			if p.Type != transactiontype.Manual {

				//checking for status of transaction
				if ok, _ := helper.Includes(transactionstatus.List, p.Status); !ok {
					return limberr.AddInvalidParam(err, "status",
						corerr.AcceptedValueForVareV, dict.R(corterm.Status),
						transactionstatus.Join())
				}
			}
		}

	case coract.Update:
		{
			if p.Description != nil {
				if len(*p.Description) > 255 {
					err = limberr.AddInvalidParam(err, "description",
						corerr.MaximumAcceptedCharacterForVisV,
						dict.R(corterm.Description), 255)
				}
			}
			//checking fot duplicate accounts
			var uniqueAcc []types.RowID
			IsUnique := true
			for _, v := range p.Slots {
				for _, k := range uniqueAcc {
					if k == v.AccountID {
						IsUnique = false
						break
					}

				}
				if IsUnique {
					uniqueAcc = append(uniqueAcc, v.AccountID)
				}
				IsUnique = true
			}

			if len(uniqueAcc) != len(p.Slots) {
				err = errors.New("duplicate accounts exists")
				return
			}
		}
	case coract.Fetch:
		{
			if ok, _ := helper.Includes(transactiontype.List, p.Type); !ok {
				return limberr.AddInvalidParam(err, "type",
					corerr.AcceptedValueForVareV, dict.R(corterm.Type),
					transactionstatus.Join())
			}
		}
	}

	return err
}

func validatePostDate(err error, postDate time.Time) error {
	defaultTime := time.Time{}
	if postDate.Format(consts.TimeLayout) == defaultTime.Format(consts.TimeLayout) {
		return limberr.AddInvalidParam(err, "postdate",
			corerr.VisRequired, dict.R(eacterm.PostDate))
	}
	return err
}

func validateSlots(err error, slots []Slot) error {
	if slots == nil || len(slots) == 0 {
		return limberr.AddInvalidParam(err, "slots",
			corerr.VisRequired, dict.R(eacterm.Slots))
	}

	return err
}

func validateCurrencyID(err error, currencyID types.RowID) error {
	if currencyID == 0 {
		return limberr.AddInvalidParam(err, "currency ID",
			corerr.VisRequired, dict.R(eacterm.CurrencyID))
	}

	return err
}

func validateRate(err error, rate float64) error {
	if rate == 0 {
		return limberr.AddInvalidParam(err, "currency ID",
			corerr.VisRequired, dict.R(eacterm.Rate))
	}

	return err
}
