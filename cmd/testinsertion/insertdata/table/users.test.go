package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/glog"
)

// InsertUsers for add required users
func InsertUsers(engine *core.Engine) {
	userRepo := basrepo.ProvideUserRepo(engine)
	userService := service.ProvideBasUserService(userRepo)

	users := []basmodel.User{
		{
			FixedCol: types.FixedCol{
				ID:        11,
				CompanyID: 1001,
				NodeID:    101,
			},
			RoleID:   1,
			Code:     "12001",
			Name:     engine.Envs[base.AdminUsername],
			Username: engine.Envs[base.AdminUsername],
			Password: engine.Envs[base.AdminPassword],
			Lang:     dict.Ku,
		},
		{
			FixedCol: types.FixedCol{
				ID:        12,
				CompanyID: 1001,
				NodeID:    101,
			},
			RoleID:   2,
			Code:     "12002",
			Name:     "cashier",
			Username: "cashier",
			Password: "cashier2020",
			Lang:     dict.En,
		},
		{
			FixedCol: types.FixedCol{
				ID:        13,
				CompanyID: 1001,
				NodeID:    101,
			},
			RoleID:   3,
			Code:     "12003",
			Name:     "reader",
			Username: "reader",
			Password: "reader2020",
			Lang:     dict.Ar,
		},
	}

	for _, v := range users {
		if _, err := userService.Save(v); err != nil {
			glog.Error(err)
		}
	}

}
