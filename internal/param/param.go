package param

import (
	"fmt"
	"omono/internal/consts"

	"github.com/syronz/dict"
)

// Param for describing request's parameter
type Param struct {
	Pagination
	Search          string
	Filter          string
	PreCondition    string
	UserID          uint
	CompanyID       uint64
	NodeID          uint64
	Lang            dict.Lang
	ErrPanel        string
	ShowDeletedRows bool
}

// Pagination is a struct, contains the fields which affected the front-end pagination
type Pagination struct {
	Select string
	Order  string
	Limit  int
	Offset int
}

// New return an intiate of the param with default limit
func New() Param {
	var param Param
	param.Limit = consts.DefaultLimit
	param.ShowDeletedRows = consts.ShowDeletedRows
	param.Order = "id"

	return param
}

// NewForDelete is used for checking delete an element
func NewForDelete(table string, col string, id interface{}) Param {
	var param Param
	param.Limit = 1
	param.Select = "*"
	param.Order = fmt.Sprintf("%v.id asc", table)
	param.PreCondition = fmt.Sprintf("%v.%v = %v", table, col, id)
	param.PreCondition += fmt.Sprintf(" AND %v.deleted_at IS NULL ", table)

	return param
}
