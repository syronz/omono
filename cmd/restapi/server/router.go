package server

import (
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basmid"
	"omono/domain/bill"
	"omono/domain/eaccounting"
	"omono/domain/location"
	"omono/domain/material"
	"omono/domain/notification"
	"omono/internal/core"

	"github.com/gin-gonic/gin"
)

// Route trigger router and api methods
func Route(rg gin.RouterGroup, engine *core.Engine) {
	// Sync Domain
	synCompanyAPI := initSynCompanyAPI(engine)

	// Base Domain
	basAuthAPI := initAuthAPI(engine)
	basUserAPI := initUserAPI(engine)
	basRoleAPI := initRoleAPI(engine)
	basSettingAPI := initSettingAPI(engine)
	basActivityAPI := initActivityAPI(engine)
	basPhoneAPI := initBasPhoneAPI(engine)
	basAccountAPI := initAccountAPI(engine, basPhoneAPI.Service)
	basCityAPI := initBasCityAPI(engine)

	// Notification Domain
	notMessageAPI := initNotMessageAPI(engine)

	// EAccountig Domain
	eacCurrencyAPI := initCurrencyAPI(engine)
	eacRateAPI := initEacRateAPI(engine)
	eacSlotAPI := initSlotAPI(engine, eacCurrencyAPI.Service, basAccountAPI.Service)
	eacTempSlotAPI := initTempSlotAPI(engine, eacCurrencyAPI.Service, basAccountAPI.Service)
	eacTransactionAPI := initTransactionAPI(engine, eacSlotAPI.Service)
	eacVoucherAPI := initVoucherAPI(engine, eacTempSlotAPI.Service)
	eacBalanceSheetAPI := initBalanceSheetAPI(engine)

	// Material Domain
	matCompanyAPI := initMatCompanyAPI(engine)
	matColorAPI := initMatColorAPI(engine)
	matGroupAPI := initMatGroupAPI(engine)
	matUnitAPI := initMatUnitAPI(engine)
	matTagAPI := initMatTagAPI(engine)
	matProductAPI := initMatProductAPI(engine)

	// Location Domain
	locStoreAPI := initLocStoreAPI(engine)

	// Bill Domain
	bilInvoiceAPI := initBilInvoiceAPI(engine)

	// Html Domain
	rg.StaticFS("/public", http.Dir("public"))

	rg.POST("/login", basAuthAPI.Login)
	rg.POST("/register", basAuthAPI.Register)

	rg.Use(basmid.AuthGuard(engine))

	rg.GET("/profile", basAuthAPI.Profile)

	access := basmid.NewAccessMid(engine)

	rg.POST("/logout", basAuthAPI.Logout)

	// Sync Domain
	rg.GET("/sync/companies",
		access.Check(base.SuperAccess), synCompanyAPI.List)
	rg.GET("/sync/companies/:companyID",
		access.Check(base.SuperAccess), synCompanyAPI.FindByID)
	rg.POST("/sync/companies",
		access.Check(base.SuperAccess), synCompanyAPI.Create)
	rg.PUT("/sync/companies/:companyID",
		access.Check(base.SuperAccess), synCompanyAPI.Update)
	rg.GET("/excel/sync/companies",
		access.Check(base.SuperAccess), synCompanyAPI.Excel)

	// Base Domain
	rg.GET("/temporary/token", basAuthAPI.TemporaryToken)

	rg.GET("/companies/:companyID/settings",
		access.Check(base.SettingRead), basSettingAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/settings/:settingID",
		access.Check(base.SettingRead), basSettingAPI.FindByID)
	rg.PUT("/companies/:companyID/nodes/:nodeID/settings/:settingID",
		access.Check(base.SettingWrite), basSettingAPI.Update)
	rg.GET("/excel/companies/:companyID/settings",
		access.Check(base.SettingExcel), basSettingAPI.Excel)

	rg.GET("/companies/:companyID/roles",
		access.Check(base.RoleRead), basRoleAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/roles/:roleID",
		access.Check(base.RoleRead), basRoleAPI.FindByID)
	rg.POST("/companies/:companyID/roles",
		access.Check(base.RoleWrite), basRoleAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/roles/:roleID",
		access.Check(base.RoleWrite), basRoleAPI.Update)
	rg.DELETE("companies/:companyID/nodes/:nodeID/roles/:roleID",
		access.Check(base.RoleWrite), basRoleAPI.Delete)
	rg.GET("/excel/companies/:companyID/roles",
		access.Check(base.RoleExcel), basRoleAPI.Excel)

	rg.GET("/companies/:companyID/accounts",
		access.Check(base.AccountRead), basAccountAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/accounts/:accountID",
		access.Check(base.AccountRead), basAccountAPI.FindByID)
	rg.POST("/companies/:companyID/accounts",
		access.Check(base.AccountWrite), basAccountAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/accounts/:accountID",
		access.Check(base.AccountWrite), basAccountAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/accounts/:accountID",
		access.Check(base.AccountWrite), basAccountAPI.Delete)
	rg.GET("/excel/companies/:companyID/accounts",
		access.Check(base.AccountExcel), basAccountAPI.Excel)
	rg.GET("/companies/:companyID/charts/accounts",
		access.Check(base.AccountRead), basAccountAPI.ChartOfAccount)
	rg.GET("/companies/:companyID/cash/account",
		access.Check(base.AccountRead), basAccountAPI.GetCashAccount)
	rg.GET("/companies/:companyID/accounts/leafs",
		access.Check(base.AccountRead), basAccountAPI.SearchLeafs)

	rg.GET("/phones",
		access.Check(base.SuperAccess), basPhoneAPI.List)
	rg.GET("/phones/:phoneID",
		access.Check(base.PhoneRead), basPhoneAPI.FindByID)
	rg.POST("/companies/:companyID/phones",
		access.Check(base.PhoneWrite), basPhoneAPI.Create)
	rg.PUT("/phones/:phoneID",
		access.Check(base.SuperAccess), basPhoneAPI.Update)
	rg.DELETE("/phones/:phoneID",
		access.Check(base.PhoneWrite), basPhoneAPI.Delete)
	rg.GET("/excel/companies/:companyID/phones",
		access.Check(base.PhoneExcel), basPhoneAPI.Excel)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/separate/:accountPhoneID",
		access.Check(base.PhoneWrite), basPhoneAPI.Separate)

	rg.GET("/username/:username",
		access.Check(base.UserRead), basUserAPI.FindByUsername)
	rg.GET("/companies/:companyID/users",
		access.Check(base.UserRead), basUserAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/users/:userID",
		access.Check(base.UserRead), basUserAPI.FindByID)
	rg.POST("/companies/:companyID/users",
		access.Check(base.UserWrite), basUserAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/users/:userID",
		access.Check(base.UserWrite), basUserAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/users/:userID",
		access.Check(base.UserWrite), basUserAPI.Delete)
	rg.GET("/excel/companies/:companyID/users",
		access.Check(base.UserExcel), basUserAPI.Excel)

	rg.GET("/activities",
		access.Check(base.SuperAccess), basActivityAPI.ListAll)
	rg.GET("/activities/companies/:companyID",
		access.Check(base.ActivityCompany), basActivityAPI.ListCompany)
	rg.GET("/activities/self",
		access.Check(base.ActivitySelf), basActivityAPI.ListSelf)

	rg.GET("/companies/:companyID/cities",
		access.Check(base.CityRead), basCityAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/cities/:cityID",
		access.Check(base.CityRead), basCityAPI.FindByID)
	rg.POST("/companies/:companyID/cities",
		access.Check(base.CityWrite), basCityAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/cities/:cityID",
		access.Check(base.CityWrite), basCityAPI.Update)
	rg.DELETE("companies/:companyID/nodes/:nodeID/cities/:cityID",
		access.Check(base.CityWrite), basCityAPI.Delete)
	rg.GET("/excel/companies/:companyID/cities",
		access.Check(base.CityExcel), basCityAPI.Excel)

	// Notification Domain
	rg.GET("/companies/:companyID/messages",
		notMessageAPI.List)
	rg.GET("/companies/:companyID/messages/:hash", notMessageAPI.ViewByHash)
	rg.GET("/companies/:companyID/nodes/:nodeID/messages/:cityID",
		access.Check(notification.MessageRead), notMessageAPI.FindByID)
	rg.POST("/companies/:companyID/messages",
		access.Check(notification.MessageWrite), notMessageAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/messages/:cityID",
		access.Check(notification.MessageWrite), notMessageAPI.Update)
	rg.DELETE("companies/:companyID/nodes/:nodeID/messages/:cityID",
		access.Check(notification.MessageWrite), notMessageAPI.Delete)
	rg.GET("/excel/companies/:companyID/messages",
		access.Check(notification.MessageExcel), notMessageAPI.Excel)

	// EAccountig Domain
	rg.GET("/companies/:companyID/currencies",
		access.Check(eaccounting.CurrencyRead), eacCurrencyAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/currencies/:currencyID",
		access.Check(eaccounting.CurrencyRead), eacCurrencyAPI.FindByID)
	rg.POST("/companies/:companyID/currencies",
		access.Check(eaccounting.CurrencyWrite), eacCurrencyAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/currencies/:currencyID",
		access.Check(eaccounting.CurrencyWrite), eacCurrencyAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/currencies/:currencyID",
		access.Check(eaccounting.CurrencyWrite), eacCurrencyAPI.Delete)
	rg.GET("/excel/companies/:companyID/currencies",
		access.Check(eaccounting.CurrencyExcel), eacCurrencyAPI.Excel)

	rg.GET("/companies/:companyID/rates",
		access.Check(eaccounting.RateRead), eacRateAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/rates/:rateID",
		access.Check(eaccounting.RateRead), eacRateAPI.FindByID)
	rg.POST("/companies/:companyID/rates",
		access.Check(eaccounting.RateWrite), eacRateAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/rates/:rateID",
		access.Check(eaccounting.RateWrite), eacRateAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/rates/:rateID",
		access.Check(eaccounting.RateWrite), eacRateAPI.Delete)
	rg.GET("/excel/companies/:companyID/rates",
		access.Check(eaccounting.RateExcel), eacRateAPI.Excel)
	rg.GET("/companies/:companyID/cities/:cityID/rates",
		access.Check(eaccounting.RateRead), eacRateAPI.RatesInCity)

	rg.GET("/companies/:companyID/transactions",
		access.Check(eaccounting.TransactionRead), eacTransactionAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/transactions/:transactionID",
		access.Check(eaccounting.TransactionRead), eacTransactionAPI.FindByID)
	rg.POST("/companies/:companyID/transactions",
		access.Check(eaccounting.TransactionManual), eacTransactionAPI.ManualTransfer)
	rg.PUT("/companies/:companyID/nodes/:nodeID/transactions/:transactionID",
		access.Check(eaccounting.TransactionUpdate), eacTransactionAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/transactions/:transactionID",
		access.Check(eaccounting.TransactionDelete), eacTransactionAPI.Delete)
	rg.GET("/excel/companies/:companyID/transactions",
		access.Check(eaccounting.TransactionExcel), eacTransactionAPI.Excel)
	rg.POST("/companies/:companyID/journals",
		access.Check(eaccounting.JournalWrite), eacTransactionAPI.JournalEntry)
	rg.PUT("/companies/:companyID/nodes/:nodeID/journals/:transactionID",
		access.Check(eaccounting.JournalWrite), eacTransactionAPI.JournalUpdate)
	rg.GET("/companies/:companyID/nodes/:nodeID/journals/print/:transactionID",
		access.Check(eaccounting.JournalPrint), eacTransactionAPI.JournalPrint)
	rg.GET("/companies/:companyID/journals/year/:year/type/:type",
		access.Check(eaccounting.LastYearCounter), eacTransactionAPI.LastYearCounterByType)

	rg.GET("/companies/:companyID/journals/counter/:counter/year/:year/type/:type",
		access.Check(eaccounting.TransactionRead), eacVoucherAPI.FindByYearCounter)

	rg.POST("/companies/:companyID/vouchers",
		access.Check(eaccounting.VoucherWrite), eacVoucherAPI.JournalVoucher)
	rg.PATCH("/companies/:companyID/nodes/:nodeID/vouchers/approve/:transactionID",
		access.Check(eaccounting.PaymentEntry), eacVoucherAPI.ApproveVoucher)
	rg.PUT("/companies/:companyID/nodes/:nodeID/vouchers/update/:voucherID",
		access.Check(eaccounting.VoucherWrite), eacVoucherAPI.VoucherUpdate)

	rg.GET("/companies/:companyID/balancesheet/:level",
		access.Check(eaccounting.BalanceSheet), eacBalanceSheetAPI.BalanceSheet)

	// Material Domain
	rg.GET("/companies/:companyID/companies",
		access.Check(material.CompanyRead), matCompanyAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/companies/:compID",
		access.Check(material.CompanyRead), matCompanyAPI.FindByID)
	rg.POST("/companies/:companyID/companies",
		access.Check(material.CompanyWrite), matCompanyAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/companies/:compID",
		access.Check(material.CompanyWrite), matCompanyAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/companies/:compID",
		access.Check(material.CompanyWrite), matCompanyAPI.Delete)
	rg.GET("/excel/companies/:companyID/companies",
		access.Check(material.CompanyExcel), matCompanyAPI.Excel)

	rg.GET("/companies/:companyID/colors",
		access.Check(material.ColorRead), matColorAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/colors/:colorID",
		access.Check(material.ColorRead), matColorAPI.FindByID)
	rg.POST("/companies/:companyID/colors",
		access.Check(material.ColorWrite), matColorAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/colors/:colorID",
		access.Check(material.ColorWrite), matColorAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/colors/:colorID",
		access.Check(material.ColorWrite), matColorAPI.Delete)
	rg.GET("/excel/companies/:companyID/colors",
		access.Check(material.ColorExcel), matColorAPI.Excel)

	rg.GET("/companies/:companyID/groups",
		access.Check(material.GroupRead), matGroupAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/groups/:groupID",
		access.Check(material.GroupRead), matGroupAPI.FindByID)
	rg.POST("/companies/:companyID/groups",
		access.Check(material.GroupWrite), matGroupAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/groups/:groupID",
		access.Check(material.GroupWrite), matGroupAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/groups/:groupID",
		access.Check(material.GroupWrite), matGroupAPI.Delete)
	rg.GET("/excel/companies/:companyID/groups",
		access.Check(material.GroupExcel), matGroupAPI.Excel)
	rg.POST("/companies/:companyID/groups/:groupID/products",
		access.Check(material.GroupWrite), matGroupAPI.AddProduct)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/group/products/:groupProductID",
		access.Check(material.GroupWrite), matGroupAPI.DelProduct)

	rg.GET("/companies/:companyID/units",
		access.Check(material.UnitRead), matUnitAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/units/:unitID",
		access.Check(material.UnitRead), matUnitAPI.FindByID)
	rg.POST("/companies/:companyID/units",
		access.Check(material.UnitWrite), matUnitAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/units/:unitID",
		access.Check(material.UnitWrite), matUnitAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/units/:unitID",
		access.Check(material.UnitWrite), matUnitAPI.Delete)
	rg.GET("/excel/companies/:companyID/units",
		access.Check(material.UnitExcel), matUnitAPI.Excel)

	rg.GET("/companies/:companyID/tags",
		access.Check(material.TagRead), matTagAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/tags/:tagID",
		access.Check(material.TagRead), matTagAPI.FindByID)
	rg.POST("/companies/:companyID/tags",
		access.Check(material.TagWrite), matTagAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/tags/:tagID",
		access.Check(material.TagWrite), matTagAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/tags/:tagID",
		access.Check(material.TagWrite), matTagAPI.Delete)
	rg.GET("/excel/companies/:companyID/tags",
		access.Check(material.TagExcel), matTagAPI.Excel)

	rg.GET("/companies/:companyID/products",
		access.Check(material.ProductRead), matProductAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/products/:productID",
		access.Check(material.ProductRead), matProductAPI.FindByID)
	rg.POST("/companies/:companyID/products",
		access.Check(material.ProductWrite), matProductAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/products/:productID",
		access.Check(material.ProductWrite), matProductAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/products/:productID",
		access.Check(material.ProductWrite), matProductAPI.Delete)
	rg.GET("/excel/companies/:companyID/products",
		access.Check(material.ProductExcel), matProductAPI.Excel)
	rg.POST("/companies/:companyID/products/:productID/tags",
		access.Check(material.ProductWrite), matProductAPI.AddTag)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/product/tags/:productTagID",
		access.Check(material.ProductWrite), matProductAPI.DelTag)

	// Location Domain
	rg.GET("/companies/:companyID/stores",
		access.Check(location.StoreRead), locStoreAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/stores/:storeID",
		access.Check(location.StoreRead), locStoreAPI.FindByID)
	rg.POST("/companies/:companyID/stores",
		access.Check(location.StoreWrite), locStoreAPI.Create)
	rg.PUT("/companies/:companyID/nodes/:nodeID/stores/:storeID",
		access.Check(location.StoreWrite), locStoreAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/stores/:storeID",
		access.Check(location.StoreWrite), locStoreAPI.Delete)
	rg.GET("/excel/companies/:companyID/stores",
		access.Check(location.StoreExcel), locStoreAPI.Excel)
	rg.POST("/companies/:companyID/stores/:storeID/users",
		access.Check(location.StoreWrite), locStoreAPI.AddUser)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/store/users/:storeUserID",
		access.Check(location.StoreWrite), locStoreAPI.DelUser)

	// Bill Domain
	rg.GET("/companies/:companyID/invoices",
		access.Check(bill.InvoiceRead), bilInvoiceAPI.List)
	rg.GET("/companies/:companyID/nodes/:nodeID/invoices/:invoiceID",
		access.Check(bill.InvoiceRead), bilInvoiceAPI.FindByID)
	rg.POST("/companies/:companyID/invoices",
		access.Check(bill.InvoiceWrite), bilInvoiceAPI.Create)
	// rg.PUT("/companies/:companyID/nodes/:nodeID/invoices/:invoiceID",
	// 	access.Check(bill.InvoiceWrite), bilInvoiceAPI.Update)
	rg.DELETE("/companies/:companyID/nodes/:nodeID/invoices/:invoiceID",
		access.Check(bill.InvoiceWrite), bilInvoiceAPI.Delete)
	rg.GET("/excel/companies/:companyID/invoices",
		access.Check(bill.InvoiceExcel), bilInvoiceAPI.Excel)
}
