package basmodel

import (
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"gorm.io/gorm"
)

// SettingTable is used inside the repo layer for specify the table name
const (
	SettingTable = "bas_settings"
)

// Setting model
type Setting struct {
	gorm.Model
	Property    types.Setting `gorm:"not null" json:"property,omitempty"`
	Value       string        `gorm:"type:text" json:"value,omitempty"`
	Type        string        `json:"type,omitempty"`
	Description string        `json:"description,omitempty"`
}

// Validate check the type of fields
func (p *Setting) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:
		if p.Property == "" {
			err = limberr.AddInvalidParam(err, "property",
				corerr.VisRequired, dict.R(corterm.Property))
		}
		fallthrough
	case coract.Update:
		if p.Value == "" {
			err = limberr.AddInvalidParam(err, "value",
				corerr.VisRequired, dict.R(corterm.Value))
		}
	}

	return err
}
