// +build wireinject

package server

import (
	"omono/domain/base/basapi"
	"omono/domain/base/basrepo"
	"omono/domain/bill/bilapi"
	"omono/domain/bill/bilrepo"
	"omono/domain/eaccounting/eacapi"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/location/locapi"
	"omono/domain/location/locrepo"
	"omono/domain/material/matapi"
	"omono/domain/material/matrepo"
	"omono/domain/notification/notapi"
	"omono/domain/notification/notrepo"
	"omono/domain/service"
	"omono/domain/sync/synapi"
	"omono/domain/sync/synrepo"

	"omono/internal/core"

	"github.com/google/wire"
)

// Sync Domain
func initSynCompanyAPI(e *core.Engine) synapi.CompanyAPI {
	wire.Build(synrepo.ProvideCompanyRepo, service.ProvideSynCompanyService,
		synapi.ProvideCompanyAPI)
	return synapi.CompanyAPI{}
}

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
	wire.Build(basrepo.ProvideAccountRepo, service.ProvideBasAccountService,
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

// EAccountig Domain
func initCurrencyAPI(e *core.Engine) eacapi.CurrencyAPI {
	wire.Build(eacrepo.ProvideCurrencyRepo, service.ProvideEacCurrencyService,
		eacapi.ProvideCurrencyAPI)
	return eacapi.CurrencyAPI{}
}

func initTransactionAPI(e *core.Engine, slotServ service.EacSlotServ) eacapi.TransactionAPI {
	wire.Build(eacrepo.ProvideTransactionRepo, service.ProvideEacTransactionService,
		eacapi.ProvideTransactionAPI)
	return eacapi.TransactionAPI{}
}

func initVoucherAPI(e *core.Engine, tempSlotServ service.EacTempSlotServ) eacapi.VoucherAPI {
	wire.Build(eacrepo.ProvideVoucherRepo, service.ProvideEacVoucherService,
		eacapi.ProvideVoucherAPI)
	return eacapi.VoucherAPI{}
}

func initBalanceSheetAPI(e *core.Engine) eacapi.BalanceSheetAPI {
	wire.Build(eacrepo.ProvideBalanceSheetRepo, service.ProvideEacBalanceSheetService,
		eacapi.ProvideBalanceSheetAPI)
	return eacapi.BalanceSheetAPI{}
}

func initSlotAPI(e *core.Engine, currencyServ service.EacCurrencyServ,
	accountServ service.BasAccountServ) eacapi.SlotAPI {
	wire.Build(eacrepo.ProvideSlotRepo, service.ProvideEacSlotService,
		eacapi.ProvideSlotAPI)
	return eacapi.SlotAPI{}
}

func initTempSlotAPI(e *core.Engine, currencyServ service.EacCurrencyServ,
	accountServ service.BasAccountServ) eacapi.TempSlotAPI {
	wire.Build(eacrepo.ProvideTempSlotRepo, service.ProvideEacTempSlotService,
		eacapi.ProvideTempSlotAPI)
	return eacapi.TempSlotAPI{}
}

func initEacRateAPI(e *core.Engine) eacapi.RateAPI {
	wire.Build(eacrepo.ProvideRateRepo, service.ProvideEacRateService, eacapi.ProvideRateAPI)
	return eacapi.RateAPI{}
}

// Material Domain
func initMatCompanyAPI(e *core.Engine) matapi.CompanyAPI {
	wire.Build(matrepo.ProvideCompanyRepo, service.ProvideMatCompanyService,
		matapi.ProvideCompanyAPI)
	return matapi.CompanyAPI{}
}

func initMatColorAPI(e *core.Engine) matapi.ColorAPI {
	wire.Build(matrepo.ProvideColorRepo, service.ProvideMatColorService,
		matapi.ProvideColorAPI)
	return matapi.ColorAPI{}
}

func initMatGroupAPI(e *core.Engine) matapi.GroupAPI {
	wire.Build(matrepo.ProvideGroupRepo, service.ProvideMatGroupService,
		matapi.ProvideGroupAPI)
	return matapi.GroupAPI{}
}

func initMatUnitAPI(e *core.Engine) matapi.UnitAPI {
	wire.Build(matrepo.ProvideUnitRepo, service.ProvideMatUnitService,
		matapi.ProvideUnitAPI)
	return matapi.UnitAPI{}
}

func initMatTagAPI(e *core.Engine) matapi.TagAPI {
	wire.Build(matrepo.ProvideTagRepo, service.ProvideMatTagService,
		matapi.ProvideTagAPI)
	return matapi.TagAPI{}
}

func initMatProductAPI(e *core.Engine) matapi.ProductAPI {
	wire.Build(matrepo.ProvideProductRepo, service.ProvideMatProductService,
		matapi.ProvideProductAPI)
	return matapi.ProductAPI{}
}

// Location Domain
func initLocStoreAPI(e *core.Engine) locapi.StoreAPI {
	wire.Build(locrepo.ProvideStoreRepo, service.ProvideLocStoreService,
		locapi.ProvideStoreAPI)
	return locapi.StoreAPI{}
}

// Bill Domain
func initBilInvoiceAPI(e *core.Engine) bilapi.InvoiceAPI {
	wire.Build(bilrepo.ProvideInvoiceRepo, service.ProvideBilInvoiceService,
		bilapi.ProvideInvoiceAPI)
	return bilapi.InvoiceAPI{}
}
