package eacmodel

// // TransactionTable is a global instance for working with transaction
// const (
// 	TransactionTable = "eac_transactions"
// )

// BalanceSheet model
type BalanceSheet struct {
	MainAccount   string         ` json:"main_account"`
	Balance       float64        ` json:"balance"`
	NameEn        string         ` json:"name_en"`
	NameKu        string         ` json:"name_ku"`
	NameAr        string         ` json:"name_ar"`
	ChildBalances []BalanceSheet ` json:"child_balances"`
}
