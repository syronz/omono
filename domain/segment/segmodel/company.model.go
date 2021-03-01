package segmodel

import (
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"gorm.io/gorm"
)

// CompanyTable is used inside the repo layer
const (
	CompanyTable = "seg_companies"
)

// Company model
type Company struct {
	gorm.Model
	Name  string  `gorm:"unique" json:"name,omitempty"`
	Phone string  `json:"phone"`
	Notes float64 `json:"notes,omitempty"`
}

// Validate check the type of fields
func (p *Company) Validate(act coract.Action) (err error) {

	// switch act {
	// case coract.Save:

	if len(p.Name) < 2 {
		err = limberr.AddInvalidParam(err, "name",
			corerr.MinimumAcceptedCharacterForVisV,
			dict.R(corterm.Name), 2)
	}

	if len(p.Name) > 200 {
		err = limberr.AddInvalidParam(err, "name",
			corerr.MaximumAcceptedCharacterForVisV,
			dict.R(corterm.Name), 255)
	}

	return err
}
