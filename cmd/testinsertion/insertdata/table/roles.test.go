package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// InsertRoles for add required roles
func InsertRoles(engine *core.Engine) {
	roleRepo := basrepo.ProvideRoleRepo(engine)
	roleService := service.ProvideBasRoleService(roleRepo)

	// reset the tables: roles, slots, transactions, accounts and users
	roleRepo.Engine.DB.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE bas_users;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE sub_accounts;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE sub_settings;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE sub_account_phones;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE sub_phones;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE bas_roles;")
	roleRepo.Engine.DB.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	roles := []basmodel.Role{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Super-Admin",
			Resources: types.ResourceJoin([]types.Resource{
				base.SettingRead, base.SettingWrite, base.SettingExcel,
				base.UserWrite, base.UserRead, base.UserExcel,
				base.ActivitySelf,
				base.RoleRead, base.RoleWrite, base.RoleExcel,
			}),
			Description: "super-admin has all privileges - do not edit",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Name: "Admin",
			Resources: types.ResourceJoin([]types.Resource{
				base.SettingRead, base.SettingWrite, base.SettingExcel,
				base.UserWrite, base.UserRead, base.UserExcel,
				base.ActivitySelf,
				base.RoleRead, base.RoleWrite, base.RoleExcel,
			}),
			Description: "admin has all privileges - do not edit",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Name:        "Cashier",
			Resources:   types.ResourceJoin([]types.Resource{base.ActivitySelf}),
			Description: "cashier has all privileges - after migration reset",
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			Name:        "for foreign 1",
			Resources:   string(base.SettingRead),
			Description: "for foreign 1",
		},
		{
			Model: gorm.Model{
				ID: 5,
			},
			Name:        "for update 1",
			Resources:   string(base.SettingRead),
			Description: "for update 1",
		},
		{
			Model: gorm.Model{
				ID: 6,
			},
			Name:        "for update 2",
			Resources:   string(base.SettingRead),
			Description: "for update 2",
		},
		{
			Model: gorm.Model{
				ID: 7,
			},
			Name:        "for delete 1",
			Resources:   string(base.SettingRead),
			Description: "for delete 1",
		},
		{
			Model: gorm.Model{
				ID: 8,
			},
			Name:        "for search 1",
			Resources:   string(base.SettingRead),
			Description: "searchTerm1",
		},
		{
			Model: gorm.Model{
				ID: 9,
			},
			Name:        "for search 2",
			Resources:   string(base.SettingRead),
			Description: "searchTerm1",
		},
		{
			Model: gorm.Model{
				ID: 10,
			},
			Name:        "for search 3",
			Resources:   string(base.SettingRead),
			Description: "searchTerm1",
		},
		{
			Model: gorm.Model{
				ID: 11,
			},
			Name:        "for delete 2",
			Resources:   string(base.SettingRead),
			Description: "for delete 2",
		},
	}

	for _, v := range roles {
		if _, err := roleService.Create(v); err != nil {
			glog.Fatal("error in creating roles", err)
		}
	}

}
