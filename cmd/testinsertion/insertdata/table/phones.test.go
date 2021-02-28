package table

import (
	"omono/domain/base/basmodel"
	"omono/domain/service"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// InsertPhones for add required users
func InsertPhones(engine *core.Engine) {
	phoneRepo := subrepo.ProvidePhoneRepo(engine)
	phoneService := service.ProvideSubPhoneService(phoneRepo)

	phones := []basmodel.Phone{
		{
			gorm.Model: gorm.Model{
				ID: 1,
			},
			//Default:   []byte("default"),
			AccountID: 1,
			Phone:     "07701001111",
			Notes:     "original",
		},

		{
			gorm.Model: gorm.Model{
				ID: 2,
			},
			//Default:   []byte("default"),
			AccountID: 2,
			Phone:     "07701002222",
			Notes:     "original",
		},
		{
			gorm.Model: gorm.Model{
				ID: 3,
			},
			//Default:   []byte("default"),
			AccountID: 3,
			Phone:     "07701003333",
			Notes:     "original",
		},
		{
			gorm.Model: gorm.Model{
				ID: 4,
			},
			//Default:   []byte("default"),
			AccountID: 4,
			Phone:     "07701004444",
			Notes:     "original",
		},
	}

	for _, v := range phones {
		if _, err := phoneService.Save(v); err != nil {
			glog.Fatal(err)
		}
	}

}
