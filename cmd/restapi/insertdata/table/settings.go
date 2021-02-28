package table

import (
	"omono/cmd/restapi/enum/settingfields"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// InsertSettings for add required settings
func InsertSettings(engine *core.Engine) {
	// engine.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %v", basmodel.SettingTable))
	settingRepo := basrepo.ProvideSettingRepo(engine)
	settingService := service.ProvideBasSettingService(settingRepo)
	settings := []basmodel.Setting{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Property:    settingfields.CompanyName,
			Value:       "item",
			Type:        "string",
			Description: "company's name in the header of invoices",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Property:    settingfields.DefaultLang,
			Value:       "ku",
			Type:        "string",
			Description: "in case of user JWT not specified this value has been used",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Property:    settingfields.CompanyLogo,
			Value:       "invoice",
			Type:        "string",
			Description: "logo for showed on the application and not invoices",
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			Property:    settingfields.InvoiceLogo,
			Value:       "public/logo.png",
			Type:        "string",
			Description: "path of logo, if branch logo won’t defined use this logo for invoices",
		},
		{
			Model: gorm.Model{
				ID: 6,
			},
			Property:    settingfields.CashAccountID,
			Value:       "2",
			Type:        "rowid",
			Description: "cash_account_id is used for returning the default account's id which is set to the default cash account",
		},
		{
			Model: gorm.Model{
				ID: 7,
			},
			Property:    settingfields.CompanyEmail,
			Value:       "XYZ@mail.com",
			Type:        "string",
			Description: "email in the header of invoice",
		},
		{
			Model: gorm.Model{
				ID: 8,
			},
			Property:    settingfields.CompanyPhone,
			Value:       "+96477000000",
			Type:        "string",
			Description: "Phone in the header of invoice",
		},
		{
			Model: gorm.Model{
				ID: 9,
			},
			Property:    settingfields.CompanyAddress,
			Value:       "Iraq Sulaimani 203452",
			Type:        "string",
			Description: "Phone in the header of invoice",
		},
		{
			Model: gorm.Model{
				ID: 10,
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
