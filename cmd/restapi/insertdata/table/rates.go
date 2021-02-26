package table

import (
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertRates for add required rates
func InsertRates(engine *core.Engine) {
	engine.DB.Exec("UPDATE eac_rates SET deleted_at = null WHERE id IN (1,2)")
	rateRepo := eacrepo.ProvideRateRepo(engine)
	rateService := service.ProvideEacRateService(rateRepo)
	rates := []eacmodel.Rate{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     1,
			Rate:       11001,
		},
		{
			FixedCol: types.FixedCol{
				ID:        2,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     2,
			Rate:       12001,
		},
		{
			FixedCol: types.FixedCol{
				ID:        3,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     3,
			Rate:       13001,
		},
		{
			FixedCol: types.FixedCol{
				ID:        4,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     4,
			Rate:       14001,
		},
		{
			FixedCol: types.FixedCol{
				ID:        5,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 2,
			CityID:     1,
			Rate:       21011,
		},
		{
			FixedCol: types.FixedCol{
				ID:        6,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 2,
			CityID:     2,
			Rate:       22011,
		},
		{
			FixedCol: types.FixedCol{
				ID:        7,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 2,
			CityID:     3,
			Rate:       23011,
		},
		{
			FixedCol: types.FixedCol{
				ID:        8,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 2,
			CityID:     4,
			Rate:       24011,
		},
		//---
		{
			FixedCol: types.FixedCol{
				ID:        9,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     1,
			Rate:       11501,
		},
		{
			FixedCol: types.FixedCol{
				ID:        10,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     2,
			Rate:       12501,
		},
		{
			FixedCol: types.FixedCol{
				ID:        11,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     3,
			Rate:       13501,
		},
		{
			FixedCol: types.FixedCol{
				ID:        12,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 1,
			CityID:     4,
			Rate:       14501,
		},
		{
			FixedCol: types.FixedCol{
				ID:        13,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 2,
			CityID:     1,
			Rate:       21511,
		},
		{
			FixedCol: types.FixedCol{
				ID:        14,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 2,
			CityID:     2,
			Rate:       22511,
		},
		// {
		// 	FixedCol: types.FixedCol{
		// 		ID:        15,
		// 		CompanyID: engine.Envs.ToUint64(sync.CompanyID),
		// 		NodeID:    engine.Envs.ToUint64(sync.NodeID),
		// 	},
		// 	CurrencyID: 2,
		// 	CityID:     3,
		// 	Rate:       23511,
		// },
		{
			FixedCol: types.FixedCol{
				ID:        16,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			CurrencyID: 2,
			CityID:     4,
			Rate:       24511,
		},
	}

	for _, v := range rates {
		if _, err := rateService.Save(v); err != nil {
			glog.Fatal("error in saving rates", err)
		}
	}

}
