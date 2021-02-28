package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/base/enum/userstatus"
	"omono/domain/service"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/pkg/glog"

	"github.com/syronz/dict"
	"gorm.io/gorm"
)

// InsertUsers for add required users
func InsertUsers(engine *core.Engine) {
	// engine.DB.Exec("DELETE FROM bas_users WHERE username = ?", engine.Envs[base.AdminUsername])
	// engine.DB.Exec("DELETE FROM bas_accounts WHERE name_en = ?", engine.Envs[base.AdminUsername])
	userRepo := basrepo.ProvideUserRepo(engine)
	userService := service.ProvideBasUserService(userRepo)
	users := []basmodel.User{
		{
			Model: gorm.Model{
				ID: consts.UserSuperAdminID,
			},
			RoleID:   1,
			Name:     engine.Envs[base.AdminUsername],
			Username: engine.Envs[base.AdminUsername],
			Password: engine.Envs[base.AdminPassword],
			Lang:     dict.Ku,
			Status:   userstatus.Active,
		},
		// {
		// 	Model: gorm.Model{
		// 	},
		// 	RoleID:   2,
		// 	Name:     "cashier",
		// 	Username: "cashier",
		// 	Password: "cashier2020",
		// 	Lang:     dict.En,
		// },
		// {
		// 	Model: gorm.Model{
		// 	},
		// 	RoleID:   3,
		// 	Name:     "reader",
		// 	Username: "reader",
		// 	Password: "reader2020",
		// 	Lang:     dict.Ar,
		// },
	}

	for _, v := range users {
		if _, err := userService.FindByID(v.ID); err == nil {
			if _, _, err := userService.Save(v); err != nil {
				glog.Fatal("error in saving users", err)
			}
		} else {
			if _, err := userService.Create(v); err != nil {
				glog.Fatal("error in creating users", err)
			}
		}
	}

}
