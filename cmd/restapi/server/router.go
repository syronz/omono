package server

import (
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basmid"
	"omono/domain/notification"
	"omono/internal/core"

	"github.com/gin-gonic/gin"
)

// Route trigger router and api methods
func Route(rg gin.RouterGroup, engine *core.Engine) {

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

	// Html Domain
	rg.StaticFS("/public", http.Dir("public"))

	rg.POST("/login", basAuthAPI.Login)
	rg.POST("/register", basAuthAPI.Register)

	rg.Use(basmid.AuthGuard(engine))

	rg.GET("/profile", basAuthAPI.Profile)

	access := basmid.NewAccessMid(engine)

	rg.POST("/logout", basAuthAPI.Logout)

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

}
