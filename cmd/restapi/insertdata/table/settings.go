package table

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/pkg/glog"

	"github.com/syronz/dict"
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
				ID: 2,
			},
			Property:    base.DefaultLang,
			Value:       string(dict.En),
			Type:        "string",
			Description: "in case of user JWT not specified this value has been used",
		},
	}

	for _, v := range settings {
		if _, err := settingService.FindByID(v.ID); err == nil {
			if _, _, err := settingService.Save(v); err != nil {
				glog.Fatal("error in saving settings", err)
			}
		} else {
			if _, _, err := settingService.Save(v); err != nil {
				glog.Fatal("error in creating settings", err)
			}
		}
	}

}
