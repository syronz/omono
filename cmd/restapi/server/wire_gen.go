// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package server

import (
	"omono/domain/base/basapi"
	"omono/domain/base/basrepo"
	"omono/domain/notification/notapi"
	"omono/domain/notification/notrepo"
	"omono/domain/segment/segapi"
	"omono/domain/segment/segrepo"
	"omono/domain/service"
	"omono/domain/subscriber/subapi"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
)

// Injectors from wire.go:

// Base Domain
func initSettingAPI(e *core.Engine) basapi.SettingAPI {
	settingRepo := basrepo.ProvideSettingRepo(e)
	basSettingServ := service.ProvideBasSettingService(settingRepo)
	settingAPI := basapi.ProvideSettingAPI(basSettingServ)
	return settingAPI
}

func initRoleAPI(e *core.Engine) basapi.RoleAPI {
	roleRepo := basrepo.ProvideRoleRepo(e)
	basRoleServ := service.ProvideBasRoleService(roleRepo)
	roleAPI := basapi.ProvideRoleAPI(basRoleServ)
	return roleAPI
}

func initUserAPI(engine *core.Engine) basapi.UserAPI {
	userRepo := basrepo.ProvideUserRepo(engine)
	basUserServ := service.ProvideBasUserService(userRepo)
	userAPI := basapi.ProvideUserAPI(basUserServ)
	return userAPI
}

func initAuthAPI(e *core.Engine) basapi.AuthAPI {
	basAuthServ := service.ProvideBasAuthService(e)
	authAPI := basapi.ProvideAuthAPI(basAuthServ)
	return authAPI
}

func initActivityAPI(engine *core.Engine) basapi.ActivityAPI {
	activityRepo := basrepo.ProvideActivityRepo(engine)
	basActivityServ := service.ProvideBasActivityService(activityRepo)
	activityAPI := basapi.ProvideActivityAPI(basActivityServ)
	return activityAPI
}

func initBasCityAPI(e *core.Engine) basapi.CityAPI {
	cityRepo := basrepo.ProvideCityRepo(e)
	basCityServ := service.ProvideBasCityService(cityRepo)
	cityAPI := basapi.ProvideCityAPI(basCityServ)
	return cityAPI
}

// Notification Domain
func initNotMessageAPI(e *core.Engine) notapi.MessageAPI {
	messageRepo := notrepo.ProvideMessageRepo(e)
	notMessageServ := service.ProvideNotMessageService(messageRepo)
	messageAPI := notapi.ProvideMessageAPI(notMessageServ)
	return messageAPI
}

// Subscriber Domain
func initSubAccountAPI(e *core.Engine, phoneServ service.SubPhoneServ) subapi.AccountAPI {
	accountRepo := subrepo.ProvideAccountRepo(e)
	subAccountServ := service.ProvideSubAccountService(accountRepo, phoneServ)
	accountAPI := subapi.ProvideAccountAPI(subAccountServ)
	return accountAPI
}

func initSubPhoneAPI(e *core.Engine) subapi.PhoneAPI {
	phoneRepo := subrepo.ProvidePhoneRepo(e)
	subPhoneServ := service.ProvideSubPhoneService(phoneRepo)
	phoneAPI := subapi.ProvidePhoneAPI(subPhoneServ)
	return phoneAPI
}

// Segment Domain
func initSegCompanyAPI(e *core.Engine) segapi.CompanyAPI {
	companyRepo := segrepo.ProvideCompanyRepo(e)
	segCompanyServ := service.ProvideSegCompanyService(companyRepo)
	companyAPI := segapi.ProvideCompanyAPI(segCompanyServ)
	return companyAPI
}
