package matmodel

import (
	"github.com/syronz/limberr"
	"omono/domain/material/enum/productstatus"
	"omono/domain/material/enum/producttype"
	"omono/domain/material/matterm"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
)

// ProductTable is a global instance for working with product
const (
	ProductTable = "mat_products"
)

// Product model
type Product struct {
	types.FixedCol
	UnitID           types.RowID `json:"unit_id"`
	Barcode          string      `gorm:"not null;unique;type:varchar(40)" json:"barcode,omitempty"`
	Name             string      `gorm:"not null;unique;type:varchar(50)" json:"name,omitempty"`
	NameEn           *string     `gorm:"type:varchar(50)" json:"name_en,omitempty"`
	NameKu           *string     `gorm:"type:varchar(50)" json:"name_ku,omitempty"`
	NameAr           *string     `gorm:"type:varchar(50)" json:"name_ar,omitempty"`
	Status           types.Enum  `gorm:"not null;default:'active';type:enum('active','inactive')" json:"status,omitempty"`
	Type             types.Enum  `gorm:"not null;default:'stocking';type:enum('serial','stocking','service')" json:"type,omitempty"`
	HasExpiry        bool        `gorm:"default:false" json:"has_expiry,omitempty"`
	Description      *string     `json:"description,omitempty"`
	Tags             []Tag       `gorm:"-" json:"tags,omitempty" table:"-"`
	PriceRetail      float64     `json:"price_retail,omitempty"`
	PriceWhole       *float64    `json:"price_whole,omitempty"`
	PriceVIP         *float64    `json:"price_vip,omitempty"`
	PriceDistributor *float64    `json:"price_distributor,omitempty"`
	PriceExport      *float64    `json:"price_export,omitempty"`
	PriceOld         *float64    `json:"price_old,omitempty"`
}

// Validate check the type of fields
func (p *Product) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Name) < 2 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 2)
		}

		if len(p.Name) > 50 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 50)
		}

		if len(p.Barcode) > 40 {
			err = limberr.AddInvalidParam(err, "barcode",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(matterm.Barcode), 40)
		}

		if p.Description != nil {
			if len(*p.Description) > 255 {
				err = limberr.AddInvalidParam(err, "description",
					corerr.MaximumAcceptedCharacterForVisV,
					dict.R(corterm.Description), 255)
			}
		}

		if p.PriceRetail <= 0 {
			err = limberr.AddInvalidParam(err, "price_retail",
				corerr.VShouldntMoreThanV,
				dict.R(matterm.PriceRetail), 0)
		}

		if ok, _ := helper.Includes(productstatus.List, p.Status); !ok {
			return limberr.AddInvalidParam(err, "status",
				corerr.AcceptedValueForVareV, dict.R(corterm.Status),
				productstatus.Join())
		}

		if ok, _ := helper.Includes(producttype.List, p.Type); !ok {
			return limberr.AddInvalidParam(err, "type",
				corerr.AcceptedValueForVareV, dict.R(corterm.Type),
				producttype.Join())
		}
	}

	return err
}
