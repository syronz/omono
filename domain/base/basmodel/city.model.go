package basmodel

import (
	"omono/domain/base/message/basterm"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"gorm.io/gorm"
)

// CityTable is used inside the repo layer
const (
	CityTable = "bas_cities"
)

// City model
type City struct {
	gorm.Model
	City  string `gorm:"not null;unique" json:"city,omitempty"`
	Notes string `json:"notes,omitempty"`
}

// Validate check the type of fields
func (p *City) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.City) < 5 {
			err = limberr.AddInvalidParam(err, "city",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(basterm.City), 5)
		}

		if len(p.City) > 255 {
			err = limberr.AddInvalidParam(err, "city",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(basterm.City), 255)
		}

		if len(p.Notes) > 255 {
			err = limberr.AddInvalidParam(err, "notes",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Notes), 255)
		}
	}

	return err
}
