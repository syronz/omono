package base

import (
	"omono/internal/types"
	"strings"
)

// settings key for base domain
const (
	DefaultRegisteredRole types.Setting = "default_registered_role"
)

// SettingList is used for validation
var SettingList = []types.Setting{
	DefaultRegisteredRole,
}

// SettingJoin make a string for showing in the api
func SettingJoin() string {
	var strArr []string

	for _, v := range SettingList {
		strArr = append(strArr, string(v))
	}

	return strings.Join(strArr, ", ")
}
