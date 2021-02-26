package notification

import "omono/internal/types"

// types for notification domain
const (
	CreateMessage types.Event = "message-create"
	UpdateMessage types.Event = "message-update"
	DeleteMessage types.Event = "message-delete"
	ListMessage   types.Event = "message-list"
	ViewMessage   types.Event = "message-view"
	ExcelMessage  types.Event = "message-excel"
)
