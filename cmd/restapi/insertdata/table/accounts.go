package table

import (
	"omono/domain/service"
	"omono/domain/subscriber/enum/accountstatus"
	"omono/domain/subscriber/enum/accounttype"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/pkg/glog"
	"omono/pkg/helper"

	"gorm.io/gorm"
)

// InsertAccounts for add required accounts
func InsertAccounts(engine *core.Engine) {
	engine.DB.Exec("UPDATE bas_accounts SET deleted_at = null WHERE id IN (1)")
	phoneServ := service.ProvideSubPhoneService(subrepo.ProvidePhoneRepo(engine))
	accountRepo := subrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideSubAccountService(accountRepo, phoneServ)
	accounts := []submodel.Account{
		{
			Model: gorm.Model{
				ID: 1,
			},
			NameEn: "Test Customer",
			NameKu: helper.StrPointer("test customer"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
	}

	for _, v := range accounts {
		if _, err := accountService.FindByID(v.ID); err == nil {
			if _, _, err := accountService.Save(v); err != nil {
				glog.Fatal("error in saving accounts", err)
			}
		} else {
			if _, err := accountService.Create(v); err != nil {
				glog.Fatal("error in creating accounts", err)
			}
		}
	}

}
