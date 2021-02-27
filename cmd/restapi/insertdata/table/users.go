package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/base/enum/userstatus"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"

	"github.com/syronz/dict"
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
				ID: 79,
			},
			RoleID:   1,
			Name:     engine.Envs[base.AdminUsername],
			Username: engine.Envs[base.AdminUsername],
			Password: engine.Envs[base.AdminPassword],
			Lang:     dict.Ku,
			Status:   userstatus.Active,
		},
		// {
		// 	FixedCol: types.FixedCol{
		// 	},
		// 	RoleID:   2,
		// 	Name:     "cashier",
		// 	Username: "cashier",
		// 	Password: "cashier2020",
		// 	Lang:     dict.En,
		// },
		// {
		// 	FixedCol: types.FixedCol{
		// 	},
		// 	RoleID:   3,
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
