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
	phoneServ := service.ProvideSubPhoneService(subrepo.ProvidePhoneRepo(engine))
	accountRepo := subrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideSubAccountService(accountRepo, phoneServ)

	// reset the accounts table
	// reset in the roles.test.go

	accounts := []submodel.Account{
		{
			Model: gorm.Model{
				ID: 1,
			},
			NameEn: "VIP",
			NameKu: helper.StrPointer("VIP"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			NameEn: "Regular",
			NameKu: helper.StrPointer("Regular"),
			Type:   accounttype.Regular,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			NameEn: "Business",
			NameKu: helper.StrPointer("Business"),
			Type:   accounttype.Business,
			Status: accountstatus.Inactive,
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			NameEn: "for foreign 1",
			NameKu: helper.StrPointer("for foreign 1"),
			Type:   accounttype.Employee,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 5,
			},
			NameEn: "for update 1",
			NameKu: helper.StrPointer("for update 1"),
			Type:   accounttype.VIP,
			Status: accountstatus.Inactive,
		},
		{
			Model: gorm.Model{
				ID: 6,
			},
			NameEn: "for update 2",
			NameKu: helper.StrPointer("for update 2"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 7,
			},
			NameEn: "for delete 1",
			NameKu: helper.StrPointer("for delete 1"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 8,
			},
			NameEn: "for search 1, searchTerm1",
			NameKu: helper.StrPointer("for search 1, searchTerm1"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 9,
			},
			NameEn: "for search 2, searchTerm1",
			NameKu: helper.StrPointer("for search 2, searchTerm1"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 10,
			},
			NameEn: "for search 3, searchTerm1",
			NameKu: helper.StrPointer("for search 3, searchTerm1"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 21,
			},
			NameEn: "for delete 2",
			NameKu: helper.StrPointer("for delete 2"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 30,
			},
			NameEn: "active provider",
			NameKu: helper.StrPointer("active provider"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 31,
			},
			NameEn: "A",
			NameKu: helper.StrPointer("A"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 32,
			},
			NameEn: "B",
			NameKu: helper.StrPointer("B"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 33,
			},
			NameEn: "C",
			NameKu: helper.StrPointer("C"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
		{
			Model: gorm.Model{
				ID: 34,
			},
			NameEn: "D",
			NameKu: helper.StrPointer("D"),
			Type:   accounttype.VIP,
			Status: accountstatus.Active,
		},
	}

	for _, v := range accounts {
		if _, err := accountService.Create(v); err != nil {
			glog.Fatal("error in creating accounts", err)
		}
	}

}
