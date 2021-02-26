package table

import (
	"omono/domain/material/enum/productstatus"
	"omono/domain/material/enum/producttype"
	"omono/domain/material/matmodel"
	"omono/domain/material/matrepo"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
	"omono/pkg/helper"
)

// InsertProducts for add required products
func InsertProducts(engine *core.Engine) {
	productRepo := matrepo.ProvideProductRepo(engine)
	productService := service.ProvideMatProductService(productRepo)
	products := []matmodel.Product{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			UnitID:      1,
			Barcode:     "r1",
			Name:        "Thinkpad X1 Extreme Gen5",
			NameEn:      helper.StrPointer("english Thinkpad X1 Extreme Gen5"),
			Status:      productstatus.Active,
			Type:        producttype.Stocking,
			HasExpiry:   false,
			Description: helper.StrPointer("sample added in insertdata"),
			PriceRetail: 400,
		},
		{
			FixedCol: types.FixedCol{
				ID:        2,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			UnitID:      1,
			Barcode:     "83242",
			Name:        "Corner Alfemo, greendize 32x171",
			NameEn:      helper.StrPointer("english Corner Alfemo, greendize 34x170"),
			Status:      productstatus.Active,
			Type:        producttype.Stocking,
			HasExpiry:   false,
			Description: helper.StrPointer("sample added in insertdata"),
			PriceRetail: 300,
		},
		{
			FixedCol: types.FixedCol{
				ID:        3,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			UnitID:      2,
			Barcode:     "23415231",
			Name:        "cocoa ice-cream",
			NameEn:      helper.StrPointer("english cocoa ice-cream"),
			Status:      productstatus.Active,
			Type:        producttype.Stocking,
			HasExpiry:   true,
			Description: helper.StrPointer("sample added in insertdata"),
			PriceRetail: 10,
		},
	}

	for _, v := range products {
		if _, err := productService.Save(v); err != nil {
			glog.Fatal("error in saving products", err)
		}
	}

}
