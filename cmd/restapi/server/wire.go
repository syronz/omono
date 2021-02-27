// +build wireinject

package server

import (
	"omono/domain/base/basapi"
	"omono/domain/base/basrepo"
	"omono/domain/notification/notapi"
	"omono/domain/notification/notrepo"
	"omono/domain/service"

	"omono/internal/core"

	"github.com/google/wire"
)

// Base Domain
func initSettingAPI(e *core.Engine) basapi.SettingAPI {
	wire.Build(basrepo.ProvideSettingRepo, service.ProvideBasSettingService,
		basapi.ProvideSettingAPI)
	return basapi.SettingAPI{}
}

func initRoleAPI(e *core.Engine) basapi.RoleAPI {
	wire.Build(basrepo.ProvideRoleRepo, service.ProvideBasRoleService,
		basapi.ProvideRoleAPI)
	return basapi.RoleAPI{}
}

func initUserAPI(engine *core.Engine) basapi.UserAPI {
	wire.Build(basrepo.ProvideUserRepo, service.ProvideBasUserService, basapi.ProvideUserAPI)
	return basapi.UserAPI{}
}

func initAuthAPI(e *core.Engine) basapi.AuthAPI {
	wire.Build(service.ProvideBasAuthService, basapi.ProvideAuthAPI)
	return basapi.AuthAPI{}
}

func initActivityAPI(engine *core.Engine) basapi.ActivityAPI {
	wire.Build(basrepo.ProvideActivityRepo, service.ProvideBasActivityService, basapi.ProvideActivityAPI)
	return basapi.ActivityAPI{}
}

func initAccountAPI(e *core.Engine, phoneServ service.BasPhoneServ) basapi.AccountAPI {
	wire.Build(basrepo.ProvideAccountRepo, service.ProvideSubAccountService,
		basapi.ProvideAccountAPI)
	return basapi.AccountAPI{}
}

func initBasPhoneAPI(e *core.Engine) basapi.PhoneAPI {
	wire.Build(basrepo.ProvidePhoneRepo, service.ProvideBasPhoneService,
		basapi.ProvidePhoneAPI)
	return basapi.PhoneAPI{}
}

func initBasCityAPI(e *core.Engine) basapi.CityAPI {
	wire.Build(basrepo.ProvideCityRepo, service.ProvideBasCityService,
		basapi.ProvideCityAPI)
	return basapi.CityAPI{}
}

// Notification Domain
func initNotMessageAPI(e *core.Engine) notapi.MessageAPI {
	wire.Build(notrepo.ProvideMessageRepo, service.ProvideNotMessageService,
		notapi.ProvideMessageAPI)
	return notapi.MessageAPI{}
}
