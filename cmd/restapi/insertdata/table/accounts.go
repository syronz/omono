package table

import (
	"omono/domain/service"
	"omono/domain/subscriber/enum/accountstatus"
	"omono/domain/subscriber/enum/accounttype"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
	"omono/pkg/helper"
)

// InsertAccounts for add required accounts
func InsertAccounts(engine *core.Engine) {
	engine.DB.Exec("UPDATE bas_accounts SET deleted_at = null WHERE id IN (1,2,3,4,5)")
	phoneServ := service.ProvideBasPhoneService(subrepo.ProvidePhoneRepo(engine))
	accountRepo := subrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideBasAccountService(accountRepo, phoneServ)
	accounts := []submodel.Account{
		{
			FixedNode: types.FixedNode{
				ID:        1,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("Asset"),
			NameKu: helper.StrPointer("asset"),
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
			NameEn:   helper.StrPointer("Cash"),
			NameKu:   helper.StrPointer("cash"),
			ParentID: types.RowIDPointer(1),
			Code:     "11",
			Type:     accounttype.Cash,
			Status:   accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        79,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn:   helper.StrPointer("Users"),
			NameKu:   helper.StrPointer("users"),
			ParentID: types.RowIDPointer(1),
			Code:     "12",
			Type:     accounttype.User,
			Status:   accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        4,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn: helper.StrPointer("Equity"),
			NameKu: helper.StrPointer("equity"),
			Code:   "2",
			Type:   accounttype.Equity,
			Status: accountstatus.Active,
		},
		{
			FixedNode: types.FixedNode{
				ID:        5,
				CompanyID: 1001,
				NodeID:    101,
			},
			NameEn:   helper.StrPointer("Capital"),
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
