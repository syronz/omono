package types

import (
	"github.com/syronz/dict"
	"strconv"
	"strings"
	"time"
)

// Envkey environment key
type Envkey string

// Envs holds all environments
type Envs map[Envkey]string

// ToBool cast string to boolean
func (p Envs) ToBool(key Envkey) bool {
	return strings.ToUpper(p[key]) == "TRUE"
}

// ToRowID casting string to RowID
func (p Envs) ToRowID(key Envkey) RowID {
	num, _ := StrToRowID(p[key])
	return num
}

// ToUint64 casting string to Uint64
func (p Envs) ToUint64(key Envkey) uint64 {
	num, _ := strconv.ParseUint(p[key], 10, 64)
	return num
}

// ToInt64 casting string to Int64
func (p Envs) ToInt64(key Envkey) int64 {
	num, _ := strconv.ParseInt(p[key], 10, 64)
	return num
}

// ToLang casting string to dict.Lang
func (p Envs) ToLang(key Envkey) dict.Lang {
	lang := dict.Lang(p[key])
	return lang
}

// ToByte casting string to []byte
func (p Envs) ToByte(key Envkey) []byte {
	return []byte(p[key])
}

// ToDuration casting string to time.Duration
func (p Envs) ToDuration(key Envkey) time.Duration {
	num := p.ToUint64(key)
	return time.Duration(num)
}

// ToInt casting string to int
func (p Envs) ToInt(key Envkey) int {
	num, _ := strconv.Atoi(p[key])
	return num
}
