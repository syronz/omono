package table

import (
	"github.com/syronz/dict"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/base/enum/accountstatus"
	"omono/domain/base/enum/accounttype"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
	"omono/pkg/helper"
)

// InsertAccounts for add required accounts
func InsertAccounts(engine *core.Engine) {
	phoneServ := service.ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountRepo := basrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideSubAccountService(accountRepo, phoneServ)

	// reset the accounts table
	// reset in the roles.test.go

	accounts := []basmodel.Account{
		{
			FixedNode: types.FixedNode{
				ID:        1,
				CompanyID: 1001,
				NodeID:    101,
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
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("Capital"),
			NameKu: helper.StrPointer("Capital"),
			Code:   "21",
			Type:   accounttype.Capital,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        3,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("Cash"),
			NameKu: helper.StrPointer("Cash"),
			Code:   "181",
			Type:   accounttype.Cash,
			Status: accountstatus.Inactive,
		},
		{
			FixedNode: types.FixedNode{
				ID:        4,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for foreign 1"),
			NameKu: helper.StrPointer("for foreign 1"),
			Code:   "181001",
			Type:   accounttype.Equity,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        5,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for update 1"),
			NameKu: helper.StrPointer("for update 1"),
			Code:   "181002",
			Type:   accounttype.Partner,
			Status: accountstatus.Inactive,
		},
		{
			FixedNode: types.FixedNode{
				ID:        6,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for update 2"),
			NameKu: helper.StrPointer("for update 2"),
			Code:   "181003",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        7,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for delete 1"),
			NameKu: helper.StrPointer("for delete 1"),
			Code:   "181004",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        8,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for search 1, searchTerm1"),
			NameKu: helper.StrPointer("for search 1, searchTerm1"),
			Code:   "181005",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        9,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for search 2, searchTerm1"),
			NameKu: helper.StrPointer("for search 2, searchTerm1"),
			Code:   "181006",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        10,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for search 3, searchTerm1"),
			NameKu: helper.StrPointer("for search 3, searchTerm1"),
			Code:   "181007",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        21,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("for delete 2"),
			NameKu: helper.StrPointer("for delete 2"),
			Code:   "181008",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        30,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("active provider"),
			NameKu: helper.StrPointer("active provider"),
			Code:   "181009",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        31,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("A"),
			NameKu: helper.StrPointer("A"),
			Code:   "181010",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        32,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("B"),
			NameKu: helper.StrPointer("B"),
			Code:   "181011",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        33,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("C"),
			NameKu: helper.StrPointer("C"),
			Code:   "181012",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        34,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("D"),
			NameKu: helper.StrPointer("D"),
			Code:   "181013",
			Type:   accounttype.Partner,
			Status: accountstatus.Active,
		},
	}

	for _, v := range accounts {
		if _, err := accountService.Save(v); err != nil {
			glog.Fatal(err)
		}
	}

}
