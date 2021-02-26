/*
Deprecated should be delete after database goes to version 16
*/
package matmodel

import (
	"github.com/syronz/limberr"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/dict"
)

// CompanyTable is a global instance for working with company
const (
	CompanyTable = "mat_companies"
)

// Company model
type Company struct {
	types.FixedCol
	Name    string `gorm:"not null;unique" json:"name,omitempty"`
	Website string `json:"website,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Address string `json:"address,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

// Validate check the type of fields
func (p *Company) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Name) < 2 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 2)
		}

		if len(p.Name) > 255 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 255)
		}

	}

	return err
}
