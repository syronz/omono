package locmodel

import (
	"github.com/syronz/limberr"
	"omono/domain/base/basmodel"
	"omono/domain/location/enum/storestatus"
	"omono/domain/location/enum/storetype"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
)

// StoreTable is a global instance for working with store
const (
	StoreTable = "loc_stores"
)

// Store model
type Store struct {
	types.FixedCol
	ParentID          *types.RowID    `json:"parent_id"`
	CityID            types.RowID     `json:"city_id,omitempty"`
	Name              string          `gorm:"not null;unique;type:varchar(40)" json:"name,omitempty"`
	Type              types.Enum      `gorm:"not null;default:'show-room';type:enum('show-room','office','head-quarter','branch','warehouse','representative')" json:"type,omitempty"`
	Code              string          `gorm:"not null;unique;type:varchar(10)" json:"code,omitempty"`
	FooterNote        string          `json:"footer_note,omitempty"`
	InvoiceThemeNew   string          `gorm:"type:varchar(100)" json:"invoice_theme_new,omitempty"`
	InvoiceThemePrint string          `gorm:"type:varchar(100)" json:"invoice_theme_print,omitempty"`
	Status            types.Enum      `gorm:"not null;default:'active';type:enum('active','inactive')" json:"status,omitempty,omitempty"`
	DiscountAccount   *types.RowID    `json:"discount_account,omitempty"`
	COGSAccount       *types.RowID    `json:"cogs_account,omitempty"`
	SaleAccount       *types.RowID    `json:"sale_account,omitempty"`
	Users             []basmodel.User `gorm:"-" json:"users,omitempty" table:"-"`
}

// Validate check the type of fields
func (p *Store) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Name) < 2 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 2)
		}

		if len(p.Name) > 40 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 40)
		}

		if len(p.Code) > 10 {
			err = limberr.AddInvalidParam(err, "code",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Code), 10)
		}

		if ok, _ := helper.Includes(storestatus.List, p.Status); !ok {
			return limberr.AddInvalidParam(err, "status",
				corerr.AcceptedValueForVareV, dict.R(corterm.Status),
				storestatus.Join())
		}

		if ok, _ := helper.Includes(storetype.List, p.Type); !ok {
			return limberr.AddInvalidParam(err, "type",
				corerr.AcceptedValueForVareV, dict.R(corterm.Type),
				storetype.Join())
		}

	}

	return err
}
