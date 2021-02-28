package server

import (
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basmid"
	"omono/domain/notification"
	"omono/domain/subscriber"
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
	basCityAPI := initBasCityAPI(engine)

	// Notification Domain
	notMessageAPI := initNotMessageAPI(engine)

	// Subscriber Domain
	basPhoneAPI := initSubPhoneAPI(engine)
	basAccountAPI := initSubAccountAPI(engine, basPhoneAPI.Service)

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

	rg.GET("/settings",
		access.Check(base.SettingRead), basSettingAPI.List)
	rg.GET("/settings/:settingID",
		access.Check(base.SettingRead), basSettingAPI.FindByID)
	rg.PUT("/settings/:settingID",
		access.Check(base.SettingWrite), basSettingAPI.Update)
	rg.GET("/excel/settings",
		access.Check(base.SettingExcel), basSettingAPI.Excel)

	rg.GET("/roles",
		access.Check(base.RoleRead), basRoleAPI.List)
	rg.GET("/roles/:roleID",
		access.Check(base.RoleRead), basRoleAPI.FindByID)
	rg.POST("/roles",
		access.Check(base.RoleWrite), basRoleAPI.Create)
	rg.PUT("/roles/:roleID",
		access.Check(base.RoleWrite), basRoleAPI.Update)
	rg.DELETE("roles/:roleID",
		access.Check(base.RoleWrite), basRoleAPI.Delete)
	rg.GET("/excel/roles",
		access.Check(base.RoleExcel), basRoleAPI.Excel)

	rg.GET("/username/:username",
		access.Check(base.UserRead), basUserAPI.FindByUsername)
	rg.GET("/users",
		access.Check(base.UserRead), basUserAPI.List)
	rg.GET("/users/:userID",
		access.Check(base.UserRead), basUserAPI.FindByID)
	rg.POST("/users",
		access.Check(base.UserWrite), basUserAPI.Create)
	rg.PUT("/users/:userID",
		access.Check(base.UserWrite), basUserAPI.Update)
	rg.DELETE("/users/:userID",
		access.Check(base.UserWrite), basUserAPI.Delete)
	rg.GET("/excel/users",
		access.Check(base.UserExcel), basUserAPI.Excel)

	rg.GET("/activities",
		access.Check(base.SuperAccess), basActivityAPI.ListAll)
	rg.GET("/activities/self",
		access.Check(base.ActivitySelf), basActivityAPI.ListSelf)

	rg.GET("/cities",
		access.Check(base.CityRead), basCityAPI.List)
	rg.GET("/cities/:cityID",
		access.Check(base.CityRead), basCityAPI.FindByID)
	rg.POST("/cities",
		access.Check(base.CityWrite), basCityAPI.Create)
	rg.PUT("/cities/:cityID",
		access.Check(base.CityWrite), basCityAPI.Update)
	rg.DELETE("/cities/:cityID",
		access.Check(base.CityWrite), basCityAPI.Delete)
	rg.GET("/excel/cities",
		access.Check(base.CityExcel), basCityAPI.Excel)

	// Notification Domain
	rg.GET("/messages",
		notMessageAPI.List)
	rg.GET("/messages/:cityID",
		access.Check(notification.MessageRead), notMessageAPI.FindByID)
	rg.GET("/hash/messages/:hash", notMessageAPI.ViewByHash)
	rg.POST("/messages",
		access.Check(notification.MessageWrite), notMessageAPI.Create)
	rg.PUT("/messages/:cityID",
		access.Check(notification.MessageWrite), notMessageAPI.Update)
	rg.DELETE("messages/:cityID",
		access.Check(notification.MessageWrite), notMessageAPI.Delete)
	rg.GET("/excel/messages",
		access.Check(notification.MessageExcel), notMessageAPI.Excel)

	// Subscriber Domain
	rg.GET("/accounts",
		access.Check(subscriber.AccountRead), basAccountAPI.List)
	rg.GET("/accounts/:accountID",
		access.Check(subscriber.AccountRead), basAccountAPI.FindByID)
	rg.POST("/accounts",
		access.Check(subscriber.AccountWrite), basAccountAPI.Create)
	rg.PUT("/accounts/:accountID",
		access.Check(subscriber.AccountWrite), basAccountAPI.Update)
	rg.DELETE("/accounts/:accountID",
		access.Check(subscriber.AccountWrite), basAccountAPI.Delete)
	rg.GET("/excel/accounts",
		access.Check(subscriber.AccountExcel), basAccountAPI.Excel)

	rg.GET("/phones",
		access.Check(base.SuperAccess), basPhoneAPI.List)
	rg.GET("/phones/:phoneID",
		access.Check(subscriber.PhoneRead), basPhoneAPI.FindByID)
	rg.POST("/phones",
		access.Check(subscriber.PhoneWrite), basPhoneAPI.Create)
	rg.PUT("/phones/:phoneID",
		access.Check(base.SuperAccess), basPhoneAPI.Update)
	rg.DELETE("/phones/:phoneID",
		access.Check(subscriber.PhoneWrite), basPhoneAPI.Delete)
	rg.GET("/excel/phones",
		access.Check(subscriber.PhoneExcel), basPhoneAPI.Excel)
	rg.DELETE("/separate/:accountPhoneID",
		access.Check(subscriber.PhoneWrite), basPhoneAPI.Separate)

}
