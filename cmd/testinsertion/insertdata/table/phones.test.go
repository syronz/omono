package table

import (
	"omono/domain/service"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// InsertPhones for add required users
func InsertPhones(engine *core.Engine) {
	phoneRepo := subrepo.ProvidePhoneRepo(engine)
	phoneService := service.ProvideSubPhoneService(phoneRepo)

	phones := []submodel.Phone{
		{
			Model: gorm.Model{
				ID: 1,
			},
			//Default:   []byte("default"),
			AccountID: 1,
			Phone:     "07701001111",
			Notes:     "original",
		},

		{
			Model: gorm.Model{
				ID: 2,
			},
			//Default:   []byte("default"),
			AccountID: 2,
			Phone:     "07701002222",
			Notes:     "original",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			//Default:   []byte("default"),
			AccountID: 3,
			Phone:     "07701003333",
			Notes:     "original",
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			//Default:   []byte("default"),
			AccountID: 4,
			Phone:     "07701004444",
			Notes:     "original",
		},
	}

	for _, v := range phones {
		if _, err := phoneService.Create(v); err != nil {
			glog.Fatal(err)
		}
	}

}
