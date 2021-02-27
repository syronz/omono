package eacmodel

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/eaccounting/eacterm"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
)

// CurrencyTable is a global instance for working with currency
const (
	CurrencyTable = "eac_currencies"
)

// Currency model
type Currency struct {
	types.FixedCol
	Name   string `gorm:"not null;unique" json:"name,omitempty"`
	Symbol string `gorm:"not null;unique" json:"symbol,omitempty"`
	Code   string `gorm:"not null;unique" json:"code,omitempty"`
}

// Validate check the type of fields
func (p *Currency) Validate(act coract.Action) (err error) {

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

		if p.Symbol == "" {
			err = limberr.AddInvalidParam(err, "symbol",
				corerr.VisRequired, dict.R(eacterm.Symbol))
		}

		if p.Code == "" {
			err = limberr.AddInvalidParam(err, "code",
				corerr.VisRequired, dict.R(corterm.Code))
		}

		if len(p.Code) > 255 {
			err = limberr.AddInvalidParam(err, "description",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Code), 255)
		}
	}

	return err
}
