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
	phoneServ := service.ProvideSubPhoneService(subrepo.ProvidePhoneRepo(engine))
	accountRepo := subrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideSubAccountService(accountRepo, phoneServ)
	accounts := []submodel.Account{
		{
			FixedCol: types.FixedCol{
				ID: 1,
			},
			NameEn: "Test Customer",
			NameKu: helper.StrPointer("test customer"),
			Code:   "1",
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
	}

	for _, v := range accounts {
		// fmt.Println("created account type:", v.Type)
		if _, err := accountService.Save(v); err != nil {
			glog.Error("error in saving accounts", err)
		}
	}

}
