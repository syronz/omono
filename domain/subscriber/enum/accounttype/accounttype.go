package accounttype

import (
	"omono/internal/types"
)

//the type of accounts
const (
	Asset     types.Enum = "asset"
	Capital   types.Enum = "capital"
	Cash      types.Enum = "cash"
	Equity    types.Enum = "Equity"
	Expense   types.Enum = "expense"
	Income    types.Enum = "income"
	Liability types.Enum = "liability"
	Partner   types.Enum = "partner"
	User      types.Enum = "user"
	Inventory types.Enum = "inventory"
	Customer  types.Enum = "customer"
)

var List = []types.Enum{
	Asset,
	Capital,
	Cash,
	Equity,
	Expense,
	Income,
	Income,
	Liability,
	Partner,
	User,
	Inventory,
	Customer,
}

var ForbiddenNegative = []types.Enum{
	Cash,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
