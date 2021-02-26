package bilmodel

import (
	"omono/internal/core/coract"
	"omono/internal/types"
	"time"
)

// InvoiceProductTable is a global instance for working with invoice_products
const (
	InvoiceProductTable = "bil_invoice_products"
)

// InvoiceProduct model
type InvoiceProduct struct {
	types.FixedCol
	InvoiceID   types.RowID `gorm:"not null" json:"invoice_id,omitempty"`
	SourceID    types.RowID `gorm:"not null" json:"source_id,omitempty"`
	DestID      types.RowID `gorm:"not null" json:"dest_id,omitempty"`
	ProductID   types.RowID `gorm:"not null" json:"product_id,omitempty"`
	Description *string     `json:"description,omitempty"`
	Notes       *string     `json:"notes,omitempty"`
	Price       float64     `json:"price,omitempty"`
	Start       *string     `json:"start,omitempty"`
	End         *string     `json:"end,omitempty"`
	QTY         *float64    `json:"qty,omitempty"`
	Discount    *float64    `json:"discount,omitempty"`
	Expiration  *time.Time  `json:"expiration,omitempty"`
	Unit        *string     `json:"unit,omitempty"`
	Barcode     *string     `json:"barcode,omitempty"`
}

// Validate check the type of fields
func (p *InvoiceProduct) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

	}

	return err
}
