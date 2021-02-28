package subscriber

import "omono/internal/types"

// types for subscribe domain
const (
	CreateAccount  types.Event = "account-create"
	UpdateAccount  types.Event = "account-update"
	DeleteAccount  types.Event = "account-delete"
	ListAccount    types.Event = "account-list"
	ChartOfAccount types.Event = "chart-of-accounts"
	ViewAccount    types.Event = "account-view"
	ExcelAccount   types.Event = "account-excel"

	CreatePhone types.Event = "phone-create"
	UpdatePhone types.Event = "phone-update"
	DeletePhone types.Event = "phone-delete"
	ListPhone   types.Event = "phone-list"
	ViewPhone   types.Event = "phone-view"
	ExcelPhone  types.Event = "phone-excel"
)
