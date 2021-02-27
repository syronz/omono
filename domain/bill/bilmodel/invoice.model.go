package bilmodel

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/bill/bilterm"
	"omono/domain/bill/enum/invoicestatus"
	"omono/domain/bill/enum/invoicetype"
	"omono/domain/bill/enum/pricemode"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/location/locmodel"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/helper"
)

// InvoiceTable is a global instance for working with invoice
const (
	InvoiceTable = "bil_invoices"
)

// Invoice model
type Invoice struct {
	types.FixedCol
	CreatedBy       types.RowID      `gorm:"not null" json:"created_by,omitempty"`
	StoreID         types.RowID      `gorm:"not null;uniqueIndex:uidx_counter_year_store;uniqueIndex:uidx_cumulative_store" json:"store_id,omitempty"`
	AccountID       types.RowID      `gorm:"not null" json:"account_id,omitempty"`
	CurrencyID      types.RowID      `gorm:"not null" json:"currency_id,omitempty"`
	CurrencyRate    float64          `json:"currency_rate"`
	Invoice         string           `gorm:"not null;unique;index:idx_invoice_number" json:"invoice,omitempty"`
	Year            int              `gorm:"not null;uniqueIndex:uidx_counter_year_store" json:"year,omitempty"`
	YearCounter     uint64           `gorm:"index:idx_year_counter;uniqueIndex:uidx_counter_year_store" json:"year_counter,omitempty"`
	YearCumulative  uint64           `gorm:"index:idx_year_cumulative;uniqueIndex:uidx_cumulative_store" json:"year_cumulative,omitempty"`
	Type            types.Enum       `gorm:"not null;default:'sale';type:enum('sale','purchase','transfer')" json:"type,omitempty"`
	Status          types.Enum       `gorm:"not null;default:'new';type:enum('new','pending','done','cancel','lock')" json:"status,omitempty,omitempty"`
	PriceMode       types.Enum       `gorm:"not null;default:'retail';type:enum('whole','vip','distributor','export','retail')" json:"price_mode,omitempty,omitempty"`
	TotalCurrency   float64          `json:"total_currency,omitempty"`
	Total           float64          `json:"total,omitempty"`
	ApplyInventory  bool             `json:"apply_inventory,omitempty"`
	ApplyAccounting bool             `json:"apply_accounting,omitempty"`
	Notes           string           `json:"notes,omitempty"`
	Store           locmodel.Store   `gorm:"-" json:"store,omitempty"`
	Rate            eacmodel.Rate    `gorm:"-" json:"rate,omitempty"`
	Products        []InvoiceProduct `gorm:"-" json:"products,omitempty"`
}

// Validate check the type of fields
func (p *Invoice) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Notes) > 255 {
			err = limberr.AddInvalidParam(err, "notes",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 255)
		}

		if ok, _ := helper.Includes(invoicestatus.List, p.Status); !ok {
			return limberr.AddInvalidParam(err, "status",
				corerr.AcceptedValueForVareV, dict.R(corterm.Status),
				invoicestatus.Join())
		}

		if ok, _ := helper.Includes(invoicetype.List, p.Type); !ok {
			return limberr.AddInvalidParam(err, "type",
				corerr.AcceptedValueForVareV, dict.R(corterm.Type),
				invoicetype.Join())
		}

		if ok, _ := helper.Includes(pricemode.List, p.PriceMode); !ok {
			return limberr.AddInvalidParam(err, "price_mode",
				corerr.AcceptedValueForVareV, dict.R(bilterm.PriceMode),
				pricemode.Join())
		}

	}

	return err
}
