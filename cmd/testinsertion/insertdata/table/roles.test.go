package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertRoles for add required roles
func InsertRoles(engine *core.Engine) {
	roleRepo := basrepo.ProvideRoleRepo(engine)
	roleService := service.ProvideBasRoleService(roleRepo)

	// reset the tables: roles, slots, transactions, accounts and users
	roleRepo.Engine.DB.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE eac_slots;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE eac_transactions;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE bas_users;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE bas_accounts;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE bas_account_phones;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE bas_phones;")
	roleRepo.Engine.DB.Exec("TRUNCATE TABLE bas_roles;")
	roleRepo.Engine.DB.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	roles := []basmodel.Role{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name: "Super-Admin",
			Resources: types.ResourceJoin([]types.Resource{
				base.SettingRead, base.SettingWrite, base.SettingExcel,
				base.UserWrite, base.UserRead, base.UserExcel,
				base.ActivitySelf, base.ActivityCompany,
				base.RoleRead, base.RoleWrite, base.RoleExcel,
			}),
			Description: "super-admin has all privileges - do not edit",
		},
		{
			FixedCol: types.FixedCol{
				ID:        2,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name: "Admin",
			Resources: types.ResourceJoin([]types.Resource{
				base.SettingRead, base.SettingWrite, base.SettingExcel,
				base.UserWrite, base.UserRead, base.UserExcel,
				base.ActivitySelf, base.ActivityCompany,
				base.RoleRead, base.RoleWrite, base.RoleExcel,
			}),
			Description: "admin has all privileges - do not edit",
		},
		{
			FixedCol: types.FixedCol{
				ID:        3,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "Cashier",
			Resources:   types.ResourceJoin([]types.Resource{base.ActivitySelf}),
			Description: "cashier has all privileges - after migration reset",
		},
		{
			FixedCol: types.FixedCol{
				ID:        4,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for foreign 1",
			Resources:   string(base.SettingRead),
			Description: "for foreign 1",
		},
		{
			FixedCol: types.FixedCol{
				ID:        5,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for update 1",
			Resources:   string(base.SettingRead),
			Description: "for update 1",
		},
		{
			FixedCol: types.FixedCol{
				ID:        6,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for update 2",
			Resources:   string(base.SettingRead),
			Description: "for update 2",
		},
		{
			FixedCol: types.FixedCol{
				ID:        7,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for delete 1",
			Resources:   string(base.SettingRead),
			Description: "for delete 1",
		},
		{
			FixedCol: types.FixedCol{
				ID:        8,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for search 1",
			Resources:   string(base.SettingRead),
			Description: "searchTerm1",
		},
		{
			FixedCol: types.FixedCol{
				ID:        9,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for search 2",
			Resources:   string(base.SettingRead),
			Description: "searchTerm1",
		},
		{
			FixedCol: types.FixedCol{
				ID:        10,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for search 3",
			Resources:   string(base.SettingRead),
			Description: "searchTerm1",
		},
		{
			FixedCol: types.FixedCol{
				ID:        11,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name:        "for delete 2",
			Resources:   string(base.SettingRead),
			Description: "for delete 2",
		},
	}

	for _, v := range roles {
		if _, err := roleService.Save(v); err != nil {
			glog.Fatal(err)
		}
	}

}
