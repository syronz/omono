package base

import "omono/internal/types"

// list of resources for base domain
const (
	Domain string = "base"

	SuperAccess types.Resource = "supper:access"
	ReadDeleted types.Resource = "deleted:read"

	UserWrite types.Resource = "user:write"
	UserRead  types.Resource = "user:read"
	UserExcel types.Resource = "user:excel"

	RoleRead  types.Resource = "role:read"
	RoleWrite types.Resource = "role:write"
	RoleExcel types.Resource = "role:excel"

	SettingRead  types.Resource = "setting:read"
	SettingWrite types.Resource = "setting:write"
	SettingExcel types.Resource = "setting:excel"

	//for activity all companies we use SuperAdmin resource
	ActivitySelf types.Resource = "activity:self"

	CityWrite types.Resource = "city:write"
	CityRead  types.Resource = "city:read"
	CityExcel types.Resource = "city:excel"

	Ping types.Resource = "ping"
)
