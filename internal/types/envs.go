package types

import (
	"omono/pkg/dict"
	"strconv"
	"strings"
	"time"
)

type Envkey string

// Envs holds all environments
type Envs map[Envkey]string

func (p Envs) ToBool(key Envkey) bool {
	return strings.ToUpper(p[key]) == "TRUE"
}

func (p Envs) ToRowID(key Envkey) RowID {
	num, _ := StrToRowID(p[key])
	return num
}

func (p Envs) ToUint64(key Envkey) uint64 {
	num, _ := strconv.ParseUint(p[key], 10, 64)
	return num
}

func (p Envs) ToInt64(key Envkey) int64 {
	num, _ := strconv.ParseInt(p[key], 10, 64)
	return num
}

func (p Envs) ToLang(key Envkey) dict.Lang {
	lang := dict.Lang(p[key])
	return lang
}

func (p Envs) ToByte(key Envkey) []byte {
	return []byte(p[key])
}

func (p Envs) ToDuration(key Envkey) time.Duration {
	num := p.ToUint64(key)
	return time.Duration(num)
}

func (p Envs) ToInt(key Envkey) int {
	num, _ := strconv.Atoi(p[key])
	return num
}
