package response

import (
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"strconv"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"
)

// NotBind use special custom_error for reduced it
func (r *Response) NotBind(err error, code, domain string, part string) {
	err = limberr.Take(err, code).Domain(domain).
		Message(corerr.ErrorInBindingV, dict.R(part)).
		Custom(corerr.BindingErr).Build()

	r.Error(err).JSON()
}

// Bind is used to make it more easear for binding items
func (r *Response) Bind(st interface{}, code, domain, part string) (err error) {
	if err = r.Context.ShouldBindJSON(&st); err != nil {
		r.NotBind(err, code, domain, part)
		return
	}

	return
}

// GetID returns the ID
func (r *Response) GetID(idIn, code, part string) (id uint, err error) {
	tmpID, err := strconv.ParseUint(idIn, 10, 32)
	err = limberr.Take(err, code).
		Message(corerr.InvalidVForV, dict.R(corterm.ID), dict.R(part)).
		Custom(corerr.ValidationFailedErr).Build()
	r.Error(err).JSON()
	id = uint(tmpID)
	return
}
