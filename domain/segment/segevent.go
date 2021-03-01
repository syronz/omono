package segment

import "omono/internal/types"

// types for subscribe domain
const (
	CreateCompany types.Event = "company-create"
	UpdateCompany types.Event = "company-update"
	DeleteCompany types.Event = "company-delete"
	ListCompany   types.Event = "company-list"
	ViewCompany   types.Event = "company-view"
	ExcelCompany  types.Event = "company-excel"
)
