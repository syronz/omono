package table

import (
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// InsertCities for add required cities
func InsertCities(engine *core.Engine) {
	engine.DB.Exec("UPDATE bas_cities SET deleted_at = null WHERE id IN (1)")
	cityRepo := basrepo.ProvideCityRepo(engine)
	cityService := service.ProvideBasCityService(cityRepo)
	cities := []basmodel.City{
		{
			Model: gorm.Model{
				ID: 1,
			},
			// TODO:  <23-02-21, yourname> use city in envs //
			City:  "Sulaimaniyah",
			Notes: "",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			City:  "Hawler",
			Notes: "please delete",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			City:  "Kirkuk",
			Notes: "please delete",
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			City:  "Duhok",
			Notes: "please delete",
		},
	}

	for _, v := range cities {
		if _, err := cityService.FindByID(v.ID); err == nil {
			if _, _, err := cityService.Save(v); err != nil {
				glog.Fatal("error in saving cities", err)
			}
		} else {
			if _, err := cityService.Create(v); err != nil {
				glog.Fatal("error in creating cities", err)
			}
		}
	}

}
