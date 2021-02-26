package matmodel

import (
	"omono/internal/types"
)

// ProductTagTable is a global instance for working with product_tag
const (
	ProductTagTable = "mat_product_tags"
)

// ProductTag model
// it is mostly controlled inside the product
type ProductTag struct {
	types.FixedCol
	ProductID types.RowID `gorm:"not null;uniqueIndex:uniqueidx_product_tag" json:"product_id"`
	TagID     types.RowID `gorm:"not null;uniqueIndex:uniqueidx_product_tag" json:"tag_id"`
}
