package table

import (
	"omono/domain/location/enum/storestatus"
	"omono/domain/location/enum/storetype"
	"omono/domain/location/locmodel"
	"omono/domain/location/locrepo"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
)

// InsertStores for add required stores
func InsertStores(engine *core.Engine) {
	storeRepo := locrepo.ProvideStoreRepo(engine)
	storeService := service.ProvideLocStoreService(storeRepo)
	stores := []locmodel.Store{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			// TODO:  <23-02-21, yourname> use store in envs //
			CityID:            1,
			Name:              "HQ",
			Type:              storetype.ShowRoom,
			Code:              "108",
			FooterNote:        "Thank you",
			InvoiceThemeNew:   "default",
			InvoiceThemePrint: "default",
			Status:            storestatus.Active,
		},
	}

	for _, v := range stores {
		if _, err := storeService.Save(v); err != nil {
			glog.Fatal("error in saving stores", err)
		}
	}

	storeUsers := []locmodel.StoreUser{
		{
			FixedCol: types.FixedCol{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			StoreID: 1,
			UserID:  79,
		},
	}

	for _, v := range storeUsers {
		if _, err := storeService.AddUser(v); err != nil {
			glog.Error("error in saving stores", err)
		}
	}

}
