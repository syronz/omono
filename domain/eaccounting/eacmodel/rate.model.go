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

// RateTable is a global instance for working with rate
const (
	RateTable = "eac_rates"
)

// Rate model
type Rate struct {
	types.FixedCol
	CurrencyID types.RowID `gorm:"not null" json:"currency_id,omitempty"`
	CityID     types.RowID `gorm:"not null" json:"city_id,omitempty"`
	Rate       float64     `json:"rate,omitempty"`
	Notes      string      `json:"notes,omitempty"`
	Name       string      `gorm:"<-:false" json:"name,omitempty" table:"eac_currencies.name"`
	Code       string      `gorm:"<-:false" json:"code,omitempty" table:"eac_currencies.code"`
}

// Validate check the type of fields
func (p *Rate) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if p.CurrencyID == 0 {
			err = limberr.AddInvalidParam(err, "currency_id",
				corerr.VisRequired, dict.R(eacterm.Currency))
		}

		if len(p.Notes) > 255 {
			err = limberr.AddInvalidParam(err, "notes",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Notes), 255)
		}

		if p.Rate == 0 {
			err = limberr.AddInvalidParam(err, "rate",
				corerr.VisRequired, dict.R(eacterm.Rate))
		}

	}

	return err
}
