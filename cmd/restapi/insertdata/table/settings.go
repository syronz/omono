package table

import (
	"fmt"
	"omono/cmd/restapi/enum/settingfields"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/bill"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertSettings for add required settings
func InsertSettings(engine *core.Engine) {
	// engine.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %v", basmodel.SettingTable))
	settingRepo := basrepo.ProvideSettingRepo(engine)
	settingService := service.ProvideBasSettingService(settingRepo)
	settings := []basmodel.Setting{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.CompanyName,
			Value:       "item",
			Type:        "string",
			Description: "company's name in the header of invoices",
		},
		{
			FixedCol: types.FixedCol{
				ID:        2,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.DefaultLang,
			Value:       "ku",
			Type:        "string",
			Description: "in case of user JWT not specified this value has been used",
		},
		{
			FixedCol: types.FixedCol{
				ID:        3,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.CompanyLogo,
			Value:       "invoice",
			Type:        "string",
			Description: "logo for showed on the application and not invoices",
		},
		{
			FixedCol: types.FixedCol{
				ID:        4,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.InvoiceLogo,
			Value:       "public/logo.png",
			Type:        "string",
			Description: "path of logo, if branch logo won’t defined use this logo for invoices",
		},
		{
			FixedCol: types.FixedCol{
				ID:        5,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property: bill.InvoiceNumberPattern,
			Value: fmt.Sprintf("%v-%v-%v", consts.InvoicePatternYear, consts.InvoicePatternStoreCode,
				consts.InvoicePatternYearCounter),
			Type:        "string",
			Description: "YEAR, STORE_CDODE, YEAR_COUNTER, YEAR_CUMULATIVE",
		},
		{
			FixedCol: types.FixedCol{
				ID:        6,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.CashAccountID,
			Value:       "2",
			Type:        "rowid",
			Description: "cash_account_id is used for returning the default account's id which is set to the default cash account",
		},
		{
			FixedCol: types.FixedCol{
				ID:        7,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.CompanyEmail,
			Value:       "XYZ@mail.com",
			Type:        "string",
			Description: "email in the header of invoice",
		},
		{
			FixedCol: types.FixedCol{
				ID:        8,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.CompanyPhone,
			Value:       "+96477000000",
			Type:        "string",
			Description: "Phone in the header of invoice",
		},
		{
			FixedCol: types.FixedCol{
				ID:        9,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.CompanyAddress,
			Value:       "Iraq Sulaimani 203452",
			Type:        "string",
			Description: "Phone in the header of invoice",
		},
		{
			FixedCol: types.FixedCol{
				ID:        10,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			Property:    settingfields.DefaultCurrency,
			Value:       "USD",
			Type:        "string",
			Description: "default currency",
		},
	}

	for _, v := range settings {
		if _, err := settingService.Save(v); err != nil {
			glog.Fatal(err)
		}
	}

}
