package types

import (
	"fmt"
	"testing"
)

type Store struct {
	FixedCol
	ParentID          *RowID `json:"parent_id"`
	CityID            RowID  `json:"city_id,omitempty"`
	Name              string `gorm:"not null;unique" json:"name,omitempty"`
	Type              Enum   `gorm:"not null" json:"type,omitempty"`
	Code              string `gorm:"not null;unique" json:"code,omitempty"`
	FooterNote        string `json:"footer_note,omitempty"`
	InvoiceThemeNew   string `json:"invoice_theme_new,omitempty"`
	InvoiceThemePrint string `json:"invoice_theme_print,omitempty"`
	Status            Enum   `json:"status,omitempty,omitempty"`
	DiscountAccount   *RowID `json:"discount_account,omitempty"`
	COGSAccount       *RowID `json:"cogs_account,omitempty"`
	SaleAccount       *RowID `json:"sale_account,omitempty"`
}

func TestExtractFixedCol(t *testing.T) {
	store := Store{}
	store.CompanyID = 1001
	store.NodeID = 101
	store.ID = 1
	store.Name = "HQ"

	r := ExtractFixedCol(store)

	fmt.Println("result of extractor: ", r)

}
