package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/pkg/glog"

	"github.com/syronz/dict"
	"gorm.io/gorm"
)

// InsertUsers for add required users
func InsertUsers(engine *core.Engine) {
	userRepo := basrepo.ProvideUserRepo(engine)
	userService := service.ProvideBasUserService(userRepo)

	users := []basmodel.User{
		{
			Model: gorm.Model{
				ID: 11,
			},
			RoleID:   1,
			Name:     engine.Envs[base.AdminUsername],
			Username: engine.Envs[base.AdminUsername],
			Password: engine.Envs[base.AdminPassword],
			Lang:     dict.Ku,
		},
		{
			Model: gorm.Model{
				ID: 12,
			},
			RoleID:   2,
			Name:     "cashier",
			Username: "cashier",
			Password: "cashier2020",
			Lang:     dict.En,
		},
		{
			Model: gorm.Model{
				ID: 13,
			},
			RoleID:   3,
			Name:     "reader",
			Username: "reader",
			Password: "reader2020",
			Lang:     dict.Ar,
		},
	}

	for _, v := range users {
		if _, err := userService.Create(v); err != nil {
			glog.Fatal("error in creating users", err)
		}
	}

}
