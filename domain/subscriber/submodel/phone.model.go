package submodel

import (
	"omono/domain/base/basterm"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"gorm.io/gorm"
)

// PhoneTable is used inside the repo layer
const (
	PhoneTable = "sub_phones"
)

// Phone model
type Phone struct {
	gorm.Model
	Phone     string `gorm:"not null;unique" json:"phone,omitempty"`
	Notes     string `json:"notes"`
	AccountID uint   `gorm:"-" json:"account_id" table:"-"`
	Default   byte   `gorm:"-" json:"default" table:"-"`
}

// Validate check the type of fields
func (p *Phone) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Phone) < 5 {
			err = limberr.AddInvalidParam(err, "phone",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(basterm.Phone), 5)
		}

		if len(p.Phone) > 13 {
			err = limberr.AddInvalidParam(err, "phone",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(basterm.Phone), 255)
		}

		if len(p.Notes) > 255 {
			err = limberr.AddInvalidParam(err, "notes",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Notes), 255)
		}
	}

	return err
}
