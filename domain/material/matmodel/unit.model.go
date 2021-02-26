package matmodel

import (
	"github.com/syronz/limberr"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/dict"
)

// UnitTable is a global instance for working with unit
const (
	UnitTable = "mat_units"
)

// Unit model
type Unit struct {
	types.FixedCol
	Name        string  `gorm:"not null;unique;type:varchar(20)" json:"name,omitempty"`
	Description *string `gorm:"type:varchar(200)" json:"description,omitempty"`
}

// Validate check the type of fields
func (p *Unit) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Name) < 2 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 2)
		}

		if len(p.Name) > 20 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 20)
		}

		if p.Description != nil {
			if len(*p.Description) > 200 {
				err = limberr.AddInvalidParam(err, "description",
					corerr.MaximumAcceptedCharacterForVisV,
					dict.R(corterm.Description), 200)
			}
		}
	}

	return err
}
