package eaccounting

import "omono/internal/types"

// list of resources for eaccounting domain
const (
	Domain string = "eaccounting"

	CurrencyWrite types.Resource = "currency:write"
	CurrencyRead  types.Resource = "currency:read"
	CurrencyExcel types.Resource = "currency:excel"

	RateWrite types.Resource = "rate:write"
	RateRead  types.Resource = "rate:read"
	RateExcel types.Resource = "rate:excel"

	TransactionRead   types.Resource = "transaction:read"
	TransactionManual types.Resource = "transaction:manual"
	TransactionUpdate types.Resource = "transaction:update"
	TransactionDelete types.Resource = "transaction:delete"
	TransactionExcel  types.Resource = "transaction:excel"
	JournalWrite      types.Resource = "journal:write"
	JournalPrint      types.Resource = "journal:print"
	LastYearCounter   types.Resource = "transaction:read"
	VoucherWrite      types.Resource = "voucher:write"
	PaymentEntry      types.Resource = "payment:entry"

	BalanceSheet types.Resource = "balancesheet:read"
)
