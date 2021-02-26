package table

import (
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertCities for add required cities
func InsertCities(engine *core.Engine) {
	engine.DB.Exec("UPDATE bas_cities SET deleted_at = null WHERE id IN (1)")
	cityRepo := basrepo.ProvideCityRepo(engine)
	cityService := service.ProvideBasCityService(cityRepo)
	cities := []basmodel.City{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			// TODO:  <23-02-21, yourname> use city in envs //
			City:  "Sulaimaniyah",
			Notes: "",
		},
		{
			FixedCol: types.FixedCol{
				ID:        2,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			City:  "Hawler",
			Notes: "please delete",
		},
		{
			FixedCol: types.FixedCol{
				ID:        3,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			City:  "Kirkuk",
			Notes: "please delete",
		},
		{
			FixedCol: types.FixedCol{
				ID:        4,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			City:  "Duhok",
			Notes: "please delete",
		},
	}

	for _, v := range cities {
		if _, err := cityService.Save(v); err != nil {
			glog.Fatal("error in saving cities", err)
		}
	}

}
