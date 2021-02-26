package corstartoff

import (
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
)

// LoadSetting read settings from database and assign them to the engine.Setting
func LoadSetting(engine *core.Engine) {

	params := param.Param{
		Pagination: param.Pagination{
			Select: "*",
			Order:  "id asc",
			Limit:  1000,
			Offset: 0,
		},
	}

	settingRepo := basrepo.ProvideSettingRepo(engine)
	var settings []basmodel.Setting
	var err error
	if settings, err = settingRepo.List(params); err != nil {
		// engine.ServerLog.Fatal(err, "failed in loading settings")
		glog.Fatal(err, "failed in loading settings")
	}

	engine.Setting = make(map[types.Setting]types.SettingMap, len(settings))

	for _, v := range settings {
		settingVal := types.SettingMap{
			Value: v.Value,
			Type:  v.Type,
		}
		engine.Setting[types.Setting(v.Property)] = settingVal
	}

}
