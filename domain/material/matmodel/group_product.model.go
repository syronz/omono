package matmodel

import (
	"omono/internal/types"
)

// GroupProductTable is a global instance for working with group_product
const (
	GroupProductTable = "mat_group_products"
)

// GroupProduct model
// it is mostly controlled inside the group
type GroupProduct struct {
	types.FixedCol
	GroupID   types.RowID `gorm:"not null;uniqueIndex:uniqueidx_product_group" json:"group_id"`
	ProductID types.RowID `gorm:"not null;uniqueIndex:uniqueidx_product_group" json:"product_id"`
}
