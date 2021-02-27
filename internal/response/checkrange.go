package response

import (
	"fmt"
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/service"
	"omono/internal/consts"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
)

// CheckRange will checks the range for companyID and nodeID
func (r *Response) CheckRange(companyID uint64) bool {
	// accessServ := service.ProvideBasAccessService(basrepo.ProvideAccessRepo(r.Engine))

	var ok bool
	var err error
	var intUserID interface{}

	if intUserID, ok = r.Context.Get("USER_ID"); !ok {
		err = limberr.New("user_id not exist in the token", "E1058403").
			Message(corerr.VNotExist, "user_id").
			Custom(corerr.ForbiddenErr).Build()
		r.Error(err).JSON()
		return false
	}

	if service.IsSuperAdmin(intUserID.(types.RowID)) {
		return true
	}

	var intCompanyID interface{}
	if intCompanyID, ok = r.Context.Get("COMPANY_ID"); !ok {
		err = limberr.New("company_id not exist in the token", "E1068491").
			Message(corerr.VNotExist, "company_id").
			Custom(corerr.ForbiddenErr).Build()
		r.Error(err).JSON()
		return false
	}
	tokenCompanyID := intCompanyID.(uint64)

	result := true

	if tokenCompanyID != companyID {
		err := limberr.New("you don't have permission to this scope", "E1052722").
			Message(corerr.ForbiddenToVV, dict.R(corterm.Scope), fmt.Sprintf("%v %v", companyID, consts.DefaultNodeID)).
			Custom(corerr.ForbiddenErr).Build()
		r.Error(err).JSON()
		result = false
	}

	return result
}
