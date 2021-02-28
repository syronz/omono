package table

import (
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertPhones for add required users
func InsertPhones(engine *core.Engine) {
	phoneRepo := basrepo.ProvidePhoneRepo(engine)
	phoneService := service.ProvideBasPhoneService(phoneRepo)

	phones := []basmodel.Phone{

		/*
			gorm.Model: gorm.Model{
				ID:        11,
				CompanyID: 1001,
				NodeID:    101,
			},
			RoleID:   1,
			Name:     engine.Envs[base.AdminUsername],
			Username: engine.Envs[base.AdminUsername],
			Password: engine.Envs[base.AdminPassword],
			Lang:     dict.Ku,
		*/

		{
			gorm.Model: gorm.Model{
				ID:        1,
				CompanyID: 1001,
				NodeID:    101,
			},
			//Default:   []byte("default"),
			AccountID: 1,
			Phone:     "07701001111",
			Notes:     "original",
		},

		{
			gorm.Model: gorm.Model{
				ID:        2,
				CompanyID: 1001,
				NodeID:    101,
			},
			//Default:   []byte("default"),
			AccountID: 2,
			Phone:     "07701002222",
			Notes:     "original",
		},
		{
			gorm.Model: gorm.Model{
				ID:        3,
				CompanyID: 1001,
				NodeID:    101,
			},
			//Default:   []byte("default"),
			AccountID: 3,
			Phone:     "07701003333",
			Notes:     "original",
		},
		{
			gorm.Model: gorm.Model{
				ID:        4,
				CompanyID: 1001,
				NodeID:    101,
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
