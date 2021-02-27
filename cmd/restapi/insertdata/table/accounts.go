package table

import (
	"github.com/syronz/dict"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/base/enum/accountstatus"
	"omono/domain/base/enum/accounttype"
	"omono/domain/base/message/basterm"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/service"
	"omono/domain/sync"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
	"omono/pkg/helper"
)

// InsertAccounts for add required accounts
func InsertAccounts(engine *core.Engine) {
	engine.DB.Exec("UPDATE bas_accounts SET deleted_at = null WHERE id IN (1,2,3,4,5)")
	phoneServ := service.ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountRepo := basrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideBasAccountService(accountRepo, phoneServ)
	accounts := []basmodel.Account{
		{
			FixedNode: types.FixedNode{
				ID:        1,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			NameEn: helper.StrPointer("Asset"),
			NameKu: helper.StrPointer(dict.T(eacterm.Asset, engine.Envs.ToLang(core.DefaultLang))),
			Code:   "1",
			Type:   accounttype.Asset,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        2,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			NameEn:   helper.StrPointer("Cash"),
			NameKu:   helper.StrPointer(dict.T(eacterm.Cash, engine.Envs.ToLang(core.DefaultLang))),
			ParentID: types.RowIDPointer(1),
			Code:     "11",
			Type:     accounttype.Cash,
			Status:   accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        79,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			NameEn:   helper.StrPointer("Users"),
			NameKu:   helper.StrPointer(dict.T(basterm.User, engine.Envs.ToLang(core.DefaultLang))),
			ParentID: types.RowIDPointer(1),
			Code:     "12",
			Type:     accounttype.User,
			Status:   accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        4,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			NameEn: helper.StrPointer("Equity"),
			NameKu: helper.StrPointer(dict.T(eacterm.Equity, engine.Envs.ToLang(core.DefaultLang))),
			Code:   "2",
			Type:   accounttype.Equity,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        5,
				CompanyID: engine.Envs.ToUint64(sync.CompanyID),
				NodeID:    engine.Envs.ToUint64(sync.NodeID),
			},
			NameEn:   helper.StrPointer("Capital"),
			NameKu:   helper.StrPointer(dict.T(eacterm.Capital, engine.Envs.ToLang(core.DefaultLang))),
			ParentID: types.RowIDPointer(4),
			Code:     "21",
			Type:     accounttype.Capital,
			Status:   accountstatus.Active,
		},
	}

	for _, v := range accounts {
		// fmt.Println("created account type:", v.Type)
		if _, err := accountService.Save(v); err != nil {
			glog.Error("error in saving accounts", err)
		}
	}

}
