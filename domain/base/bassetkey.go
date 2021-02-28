package base

import (
	"omono/internal/types"
	"strings"
)

// settings key for base domain
const (
	DefaultLang           types.Setting = "default_language"
	DefaultRegisteredRole types.Setting = "default_registered_role"
)

// List is used for validation
var List = []types.Setting{
	DefaultLang,
	DefaultRegisteredRole,
}

// Join make a string for showing in the api
func Join() string {
	var strArr []string

	for _, v := range List {
		strArr = append(strArr, string(v))
	}

	return strings.Join(strArr, ", ")
}
