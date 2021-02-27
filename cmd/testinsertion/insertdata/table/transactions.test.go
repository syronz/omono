package table

import (
	"omono/domain/base/basrepo"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/eaccounting/enum/transactiontype"
	"omono/domain/service"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
	"omono/pkg/helper"
	"time"
)

// InsertTransactions for add required accounts
func InsertTransactions(engine *core.Engine) {
	phoneServ := service.ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountRepo := basrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideSubAccountService(accountRepo, phoneServ)

	currencyRepo := eacrepo.ProvideCurrencyRepo(engine)
	currencyService := service.ProvideEacCurrencyService(currencyRepo)

	slotRepo := eacrepo.ProvideSlotRepo(engine)
	slotService := service.ProvideEacSlotService(slotRepo, currencyService, accountService)

	transactionRepo := eacrepo.ProvideTransactionRepo(engine)
	transactionService := service.ProvideEacTransactionService(transactionRepo, slotService)

	time1, _ := time.Parse(consts.TimeLayoutZone, "2020-10-19 15:10:00 +0300")

	transactions := []eacmodel.Transaction{
		{ // A -- 1000$ -- > B
			FixedCol: types.FixedCol{
				CompanyID: 1001,
				NodeID:    101,
			},
			Pioneer:     31,
			Follower:    32,
			CurrencyID:  1,
			Amount:      1000,
			PostDate:    time1,
			Description: helper.StrPointer("A -> B : 1000$"),
			Type:        transactiontype.Manual,
			CreatedBy:   11,
		},
		{ // A -- 800$ -- > B
			FixedCol: types.FixedCol{
				CompanyID: 1001,
				NodeID:    101,
			},
			Pioneer:     31,
			Follower:    32,
			CurrencyID:  1,
			Amount:      800,
			PostDate:    time1,
			Description: helper.StrPointer("A -> B : 800$"),
			Type:        transactiontype.Manual,
			CreatedBy:   11,
		},
		{ // C -- 200$ -- > B
			FixedCol: types.FixedCol{
				CompanyID: 1001,
				NodeID:    101,
			},
			Pioneer:     33,
			Follower:    32,
			CurrencyID:  1,
			Amount:      200,
			PostDate:    time1,
			Description: helper.StrPointer("C -> B : 200$"),
			Type:        transactiontype.Manual,
			CreatedBy:   11,
		},
		{ // D -- 300$ -- > A
			FixedCol: types.FixedCol{
				CompanyID: 1001,
				NodeID:    101,
			},
			Pioneer:     34,
			Follower:    31,
			CurrencyID:  1,
			Amount:      300,
			PostDate:    time1,
			Description: helper.StrPointer("D -> A : 300$"),
			Type:        transactiontype.Manual,
			CreatedBy:   11,
		},
	}

	for _, v := range transactions {
		if _, err := transactionService.Transfer(v); err != nil {
			glog.Fatal(err)
		}
	}

}

// InsertJournals is used for testing journal entry/delete/update
func InsertJournals(engine *core.Engine) {
	phoneServ := service.ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountRepo := basrepo.ProvideAccountRepo(engine)
	accountService := service.ProvideSubAccountService(accountRepo, phoneServ)

	currencyRepo := eacrepo.ProvideCurrencyRepo(engine)
	currencyService := service.ProvideEacCurrencyService(currencyRepo)

	slotRepo := eacrepo.ProvideSlotRepo(engine)
	slotService := service.ProvideEacSlotService(slotRepo, currencyService, accountService)

	transactionRepo := eacrepo.ProvideTransactionRepo(engine)
	transactionService := service.ProvideEacTransactionService(transactionRepo, slotService)

	time2, _ := time.Parse(consts.TimeLayoutZone, "2020-12-10 12:00:00 +0300")

	journals := []eacmodel.Transaction{
		{
			/*
				A:31	5	200	0	-1300
				B:32	5	800	0	2800
				C:33	5	0	1000	-1200
			*/
			FixedCol: types.FixedCol{
				CompanyID: 1001,
				NodeID:    101,
			},
			CurrencyID:  1,
			PostDate:    time2,
			Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
			Type:        transactiontype.JournalEntry,
			CreatedBy:   11,
			Slots: []eacmodel.Slot{
				{
					FixedCol: types.FixedCol{
						CompanyID: 1001,
						NodeID:    101,
					},
					CurrencyID:  1,
					AccountID:   31,
					Debit:       200,
					Credit:      0,
					Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					PostDate:    time2,
				},
				{
					FixedCol: types.FixedCol{
						CompanyID: 1001,
						NodeID:    101,
					},
					CurrencyID:  1,
					AccountID:   32,
					Debit:       800,
					Credit:      0,
					Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					PostDate:    time2,
				},
				{
					FixedCol: types.FixedCol{
						CompanyID: 1001,
						NodeID:    101,
					},
					CurrencyID:  1,
					AccountID:   33,
					Debit:       0,
					Credit:      1000,
					Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					PostDate:    time2,
				},
			},
		},
		{
			FixedCol: types.FixedCol{
				CompanyID: 1001,
				NodeID:    101,
			},
			CurrencyID:  1,
			PostDate:    time2,
			Description: helper.StrPointer("B cr300, D dr300"),
			Type:        transactiontype.JournalEntry,
			CreatedBy:   11,
			Slots: []eacmodel.Slot{
				{
					FixedCol: types.FixedCol{
						CompanyID: 1001,
						NodeID:    101,
					},
					CurrencyID:  1,
					AccountID:   32,
					Debit:       0,
					Credit:      300,
					Description: helper.StrPointer("B cr300, D dr300"),
					PostDate:    time2,
				},
				{
					FixedCol: types.FixedCol{
						CompanyID: 1001,
						NodeID:    101,
					},
					CurrencyID:  1,
					AccountID:   34,
					Debit:       300,
					Credit:      0,
					Description: helper.StrPointer("B cr300, D dr300"),
					PostDate:    time2,
				},
			},
		},
	}

	for _, v := range journals {
		if _, err := transactionService.JournalEntry(v); err != nil {
			glog.Fatal(err)
		}
	}

}
