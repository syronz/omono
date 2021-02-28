package table

import (
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertCurrencies for add required currencies
func InsertCurrencies(engine *core.Engine) {
	engine.DB.Exec("UPDATE eac_currencies SET deleted_at = null WHERE id IN (1,2)")
	currencyRepo := eacrepo.ProvideCurrencyRepo(engine)
	currencyService := service.ProvideEacCurrencyService(currencyRepo)
	currencies := []eacmodel.Currency{
		{
			gorm.Model: gorm.Model{
				ID:        1,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:   "Dollar",
			Symbol: "$",
			Code:   "USD",
		},
		{
			gorm.Model: gorm.Model{
				ID:        2,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:   "Dinar",
			Symbol: "IQD",
			Code:   "IQD",
		},
	}

	for _, v := range currencies {
		if _, err := currencyService.Save(v); err != nil {
			glog.Fatal("error in saving currencies", err)
		}
	}

}
