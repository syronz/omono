package settingfields

import (
	"omono/internal/types"
	"strings"
)

//setting fields
const (
	CompanyName     types.Setting = "company_name"
	ReceiptHeader   types.Setting = "receipt_header"
	DefaultLang     types.Setting = "default_language"
	CompanyLogo     types.Setting = "company_logo"
	InvoiceLogo     types.Setting = "invoice_logo"
	ReceiptPhone    types.Setting = "receipt_phone"
	ReceiptAddress  types.Setting = "receipt_address"
	CashAccountID   types.Setting = "cash_account_id"
	CompanyEmail    types.Setting = "company_email"
	CompanyPhone    types.Setting = "company_phone"
	CompanyAddress  types.Setting = "company_address"
	MainAssetID     types.Setting = "main_asset_id"
	MainLiabilityID types.Setting = "main_liability_id"
	MainEquityID    types.Setting = "main_equity_id"
	MainRevenueID   types.Setting = "main_revenue_id"
	MainExpenseID   types.Setting = "main_expense_id"
	MainSalesID     types.Setting = "main_sales_id"
	MainCgsID       types.Setting = "main_cgs_id"
	DefaultCurrency types.Setting = "default_currency"
)

//List of setting types
var List = []types.Setting{
	CompanyName,
	ReceiptHeader,
	DefaultLang,
	CompanyLogo,
	InvoiceLogo,
	ReceiptPhone,
	ReceiptAddress,
	CashAccountID,
	CompanyEmail,
	CompanyPhone,
	MainAssetID,
	MainLiabilityID,
	MainEquityID,
	MainRevenueID,
	MainExpenseID,
	MainSalesID,
	MainCgsID,
	DefaultCurrency,
}

// Join make a string for showing in the api
func Join() string {
	var strArr []string

	for _, v := range List {
		strArr = append(strArr, string(v))
	}

	return strings.Join(strArr, ", ")
}
