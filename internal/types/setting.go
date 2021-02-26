package types

type Setting string

type SettingMap struct {
	Value string
	Type  string
}

// ToRowID return RowID for id
func (p SettingMap) ToRowID() RowID {
	n, _ := StrToRowID(p.Value)
	return n
}
