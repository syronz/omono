package basmodel

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
)

// Auth model
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate check the type of fields for auth
func (p *Auth) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Login:
		if p.Username == "" {
			err = limberr.AddInvalidParam(err, "username",
				corerr.VisRequired, dict.R(corterm.Username))
		}

		if p.Password == "" {
			err = limberr.AddInvalidParam(err, "password",
				corerr.VisRequired, dict.R(corterm.Password))
		}
	}

	return err
}
