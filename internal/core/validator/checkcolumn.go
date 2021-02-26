package validator

import (
	"github.com/syronz/limberr"
	"omono/internal/core/corerr"
	"omono/pkg/helper"
	"strings"
)

// CheckColumns will check columns for security
func CheckColumns(cols []string, requestedCols string) (string, error) {
	var err error

	if requestedCols == "*" || requestedCols == "" {
		return strings.Join(cols, ","), nil
	}

	variates := strings.Split(requestedCols, ",")
	for _, v := range variates {
		if ok, _ := helper.Includes(cols, v); !ok {
			err = limberr.AddInvalidParam(err, v, corerr.VisNotValid, v)
			err = limberr.SetCustom(err, corerr.ValidationFailedErr)
		}
	}

	return requestedCols, err

}
