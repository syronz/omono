package eaccounting

import "omono/internal/types"

// types for eaccounting domain
const (
	CreateCurrency types.Event = "currency-create"
	UpdateCurrency types.Event = "currency-update"
	DeleteCurrency types.Event = "currency-delete"
	ListCurrency   types.Event = "currency-list"
	ViewCurrency   types.Event = "currency-view"
	ExcelCurrency  types.Event = "currency-excel"

	CreateRate types.Event = "rate-create"
	UpdateRate types.Event = "rate-update"
	DeleteRate types.Event = "rate-delete"
	ListRate   types.Event = "rate-list"
	ViewRate   types.Event = "rate-view"
	ExcelRate  types.Event = "rate-excel"

	CreateTransaction types.Event = "transaction-create"
	UpdateTransaction types.Event = "transaction-update"
	DeleteTransaction types.Event = "transaction-delete"
	ListTransaction   types.Event = "transaction-list"
	ViewTransaction   types.Event = "transaction-view"
	ExcelTransaction  types.Event = "transaction-excel"
	ManualTransfer    types.Event = "transaction-manual"
	EnterJournal      types.Event = "journal-entered"
	UpdateJournal     types.Event = "journal-update"
	PrintJorunal      types.Event = "journal-print"
)
