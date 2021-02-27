package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/notification"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertRoles for add required roles
func InsertRoles(engine *core.Engine) {
	engine.DB.Exec("UPDATE bas_roles SET deleted_at = null WHERE id IN (1,2,3)")
	roleRepo := basrepo.ProvideRoleRepo(engine)
	roleService := service.ProvideBasRoleService(roleRepo)
	roles := []basmodel.Role{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name: "Admin",
			Resources: types.ResourceJoin([]types.Resource{
				base.SuperAccess, base.ReadDeleted,
				base.SettingRead, base.SettingWrite, base.SettingExcel,
				base.UserWrite, base.UserRead, base.UserExcel,
				base.ActivitySelf,
				base.RoleRead, base.RoleWrite, base.RoleExcel,
				base.AccountRead, base.AccountWrite, base.AccountExcel,
				base.PhoneRead, base.PhoneWrite, base.PhoneExcel,
				base.CityRead, base.CityWrite, base.CityExcel,
				notification.MessageWrite, notification.MessageExcel,
			}),
			Description: "admin has all privileges - do not edit",
		},
		{
			FixedCol: types.FixedCol{
				ID:        2,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name: "Cashier",
			Resources: types.ResourceJoin([]types.Resource{
				base.ActivitySelf,
				base.AccountRead, base.AccountWrite, base.AccountExcel,
			}),
			Description: "cashier has privileges for adding transactions - after migration reset",
		},
		{
			FixedCol: types.FixedCol{
				ID:        3,
				CompanyID: 1001,
				NodeID:    101,
			},
			Name: "Reader",
			Resources: types.ResourceJoin([]types.Resource{
				base.SettingRead, base.SettingExcel,
				base.UserRead, base.UserExcel,
				base.RoleRead, base.RoleExcel,
			}),
			Description: "Reader can see all part without changes",
		},
		{
			FixedCol: types.FixedCol{
				ID:        4,
				CompanyID: 1002,
				NodeID:    101,
			},
			Name: "should_be_deleted",
			Resources: types.ResourceJoin([]types.Resource{
				base.SettingRead, base.SettingExcel,
				base.UserRead, base.UserExcel,
				base.RoleRead, base.RoleExcel,
			}),
			Description: "Reader can see all part without changes",
		},
	}

	for _, v := range roles {
		if _, err := roleService.Save(v); err != nil {
			glog.Fatal("error in saving roles", err)
		}

	}

}
