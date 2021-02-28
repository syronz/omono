package types

type Setting string

type SettingMap struct {
	Value string
	Type  string
}

// Touint return uint for id
func (p SettingMap) Touint() uint {
	n, _ := StrToUint(p.Value)
	return n
}
