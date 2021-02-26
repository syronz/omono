package location

import "omono/internal/types"

// types for location domain
const (
	CreateStore types.Event = "store-create"
	UpdateStore types.Event = "store-update"
	DeleteStore types.Event = "store-delete"
	ListStore   types.Event = "store-list"
	ViewStore   types.Event = "store-view"
	ExcelStore  types.Event = "store-excel"
	AddUser     types.Event = "store-add-user"
	DelUser     types.Event = "store-delete-user"
)
