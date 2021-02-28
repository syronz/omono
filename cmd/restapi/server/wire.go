// +build wireinject

package server

import (
	"omono/domain/base/basapi"
	"omono/domain/base/basrepo"
	"omono/domain/notification/notapi"
	"omono/domain/notification/notrepo"
	"omono/domain/service"
	"omono/domain/subscriber/subapi"
	"omono/domain/subscriber/subrepo"

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

// Subscriber Domain
func initSubAccountAPI(e *core.Engine, phoneServ service.SubPhoneServ) subapi.AccountAPI {
	wire.Build(subrepo.ProvideAccountRepo, service.ProvideSubAccountService,
		subapi.ProvideAccountAPI)
	return subapi.AccountAPI{}
}

func initSubPhoneAPI(e *core.Engine) subapi.PhoneAPI {
	wire.Build(subrepo.ProvidePhoneRepo, service.ProvideSubPhoneService,
		subapi.ProvidePhoneAPI)
	return subapi.PhoneAPI{}
}
