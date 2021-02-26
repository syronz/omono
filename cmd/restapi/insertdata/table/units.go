package table

import (
	"omono/domain/material/matmodel"
	"omono/domain/material/matrepo"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
	"omono/pkg/helper"
)

// InsertUnits for add required units
func InsertUnits(engine *core.Engine) {
	unitRepo := matrepo.ProvideUnitRepo(engine)
	unitService := service.ProvideMatUnitService(unitRepo)
	units := []matmodel.Unit{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Name:        "index",
			Description: helper.StrPointer("simple number"),
		},
		{
			FixedCol: types.FixedCol{
				ID:        2,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Name:        "Kg",
			Description: helper.StrPointer("simple number"),
		},
		{
			FixedCol: types.FixedCol{
				ID:        3,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Name:        "box",
			Description: helper.StrPointer("simple number"),
		},
	}

	for _, v := range units {
		if _, err := unitService.Save(v); err != nil {
			glog.Fatal("error in saving units", err)
		}
	}

}
