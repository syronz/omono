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

// ColorTable is a global instance for working with color
const (
	ColorTable = "mat_colors"
)

// Color model
type Color struct {
	types.FixedCol
	Name string `gorm:"not null;unique" json:"name,omitempty"`
	Code string `gorm:"not null;unique" json:"code,omitempty"`
}

// Validate check the type of fields
func (p *Color) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if p.Code == "" {
			err = limberr.AddInvalidParam(err, "code",
				corerr.VisRequired,
				dict.R(corterm.Code))
		}

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
