package matmodel

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/material/enum/productstatus"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/helper"
)

// GroupTable is a global instance for working with group
const (
	GroupTable = "mat_groups"
)

// Group model
type Group struct {
	types.FixedCol
	Name     string     `gorm:"not null;unique" json:"name,omitempty"`
	Status   types.Enum `gorm:"not null;default:'active';type:enum('active','inactive')" json:"status,omitempty"`
	Products []Product  `gorm:"-" json:"products,omitempty" table:"-"`
}

// Validate check the type of fields
func (p *Group) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Name) < 3 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 2)
		}

		if len(p.Name) > 255 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 255)
		}

		if ok, _ := helper.Includes(productstatus.List, p.Status); !ok {
			return limberr.AddInvalidParam(err, "status",
				corerr.AcceptedValueForVareV, dict.R(corterm.Status),
				productstatus.Join())
		}

	}

	return err
}
