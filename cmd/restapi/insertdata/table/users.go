package table

import (
	"github.com/syronz/dict"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertUsers for add required users
func InsertUsers(engine *core.Engine) {
	// engine.DB.Exec("DELETE FROM bas_users WHERE username = ?", engine.Envs[base.AdminUsername])
	// engine.DB.Exec("DELETE FROM bas_accounts WHERE name_en = ?", engine.Envs[base.AdminUsername])
	userRepo := basrepo.ProvideUserRepo(engine)
	userService := service.ProvideBasUserService(userRepo)
	users := []basmodel.User{
		{
			FixedCol: types.FixedCol{
				ID:        79,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			RoleID:   1,
			Code:     "112001",
			Name:     engine.Envs[base.AdminUsername],
			Username: engine.Envs[base.AdminUsername],
			Password: engine.Envs[base.AdminPassword],
			Lang:     dict.Ku,
		},
		// {
		// 	FixedCol: types.FixedCol{
		// 		CompanyID: engine.Envs.ToUint64(sync.CompanyID),
		// 		NodeID:    engine.Envs.ToUint64(sync.NodeID),
		// 	},
		// 	RoleID:   2,
		// 	Code:     "112002",
		// 	Name:     "cashier",
		// 	Username: "cashier",
		// 	Password: "cashier2020",
		// 	Lang:     dict.En,
		// },
		// {
		// 	FixedCol: types.FixedCol{
		// 		CompanyID: engine.Envs.ToUint64(sync.CompanyID),
		// 		NodeID:    engine.Envs.ToUint64(sync.NodeID),
		// 	},
		// 	RoleID:   3,
		// 	Code:     "112003",
		// 	Name:     "reader",
		// 	Username: "reader",
		// 	Password: "reader2020",
		// 	Lang:     dict.Ar,
		// },
	}

	for _, v := range users {
		if _, err := userService.Save(v); err != nil {
			glog.Fatal("error in saving users", err)
		}
	}

}
