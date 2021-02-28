package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// InsertSettings for add required settings
func InsertSettings(engine *core.Engine) {
	settingRepo := basrepo.ProvideSettingRepo(engine)
	settingService := service.ProvideBasSettingService(settingRepo)

	// reset the table by deleting everything
	engine.DB.Exec("TRUNCATE TABLE bas_settings;")

	settings := []basmodel.Setting{
		{
			Model: gorm.Model{
				ID: 2,
			},
			Property:    base.DefaultLang,
			Value:       "ku",
			Type:        "string",
			Description: "in case of user JWT not specified this value has been used",
		},
	}

	for _, v := range settings {
		if _, _, err := settingService.Save(v); err != nil {
			glog.Fatal("error in inserting settings", err)
		}
	}

}
