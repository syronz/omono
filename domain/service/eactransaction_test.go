package service

import (
	"omono/domain/base/basrepo"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/eaccounting/enum/transactiontype"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/helper"
	"omono/test/kernel"
	"testing"
	"time"
)

func initTransactionTest() (engine *core.Engine, transactionServ EacTransactionServ) {
	logQuery, debugLevel := initServiceTest()
	engine = kernel.StartMotor(logQuery, debugLevel)

	phoneServ := ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountServ := ProvideBasAccountService(basrepo.ProvideAccountRepo(engine), phoneServ)
	currencyServ := ProvideEacCurrencyService(eacrepo.ProvideCurrencyRepo(engine))
	slotServ := ProvideEacSlotService(eacrepo.ProvideSlotRepo(engine), currencyServ, accountServ)
	transactionServ = ProvideEacTransactionService(eacrepo.ProvideTransactionRepo(engine), slotServ)

	return
}

func TestTransactionTransfer(t *testing.T) {
	_, transactionServ := initTransactionTest()
	// time1, err := time.Parse(consts.TimeLayout, "2020-10-20 15:10:00")
	time1, err := time.Parse(consts.TimeLayoutZone, "2020-10-20 15:10:00 +0300")
	if err != nil {
		t.Errorf("error in parsing date %v in layout %v", consts.DefaultLimit, "2020-10-21 21:10:35")
	}

	samples := []struct {
		in  eacmodel.Transaction
		err error
	}{
		{in: eacmodel.Transaction{
			FixedNode: types.FixedNode{
				CompanyID: 1001,
				NodeID:    101,
			},
			Type:       transactiontype.Manual,
			CreatedBy:  11,
			Pioneer:    31,
			Follower:   32,
			CurrencyID: 1,
			Amount:     1000,
			PostDate:   time1,
		},
			err: nil,
		},
	}

	for _, v := range samples {
		_, err := transactionServ.Transfer(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}

}

func TestTransactionEditManual(t *testing.T) {
	_, transactionServ := initTransactionTest()
	time1, err := time.Parse(consts.TimeLayoutZone, "2020-10-19 15:10:00 +0300")
	if err != nil {
		t.Errorf("error in parsing date %v in layout %v", consts.DefaultLimit, "2020-10-21 21:10:35")
	}

	samples := []struct {
		in  eacmodel.Transaction
		err error
	}{
		{
			in: eacmodel.Transaction{
				FixedNode: types.FixedNode{
					CompanyID: 1001,
					NodeID:    101,
					ID:        1,
				},
				Pioneer:     33,
				Follower:    32,
				CurrencyID:  1,
				Amount:      500,
				PostDate:    time1,
				Description: helper.StrPointer("changed!, A -> C & 1000$ -> 500$"),
				Type:        transactiontype.Manual,
				CreatedBy:   11,
			},
			err: nil,
		},
	}

	for _, v := range samples {
		_, err := transactionServ.EditTransfer(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}

}

func TestTransactionDelete(t *testing.T) {
	_, transactionServ := initTransactionTest()

	sample := types.FixedCol{
		CompanyID: 1001,
		NodeID:    101,
		ID:        2,
	}

	if _, err := transactionServ.Delete(sample); err != nil {
		t.Errorf("there is an error for deleting a transaction, %v", err.Error())
	}
}

func TestJournalUpdate(t *testing.T) {
	_, transactionServ := initTransactionTest()

	time1, err := time.Parse(consts.TimeLayoutZone, "2020-10-18 15:10:00 +0300")
	_ = time1
	time2, err := time.Parse(consts.TimeLayoutZone, "2020-12-10 12:00:00 +0300")
	if err != nil {
		t.Errorf("error in parsing date %v in layout %v", consts.DefaultLimit, "2020-12-10 12:00:00 +0300")
	}

	samples := []struct {
		in  eacmodel.Transaction
		err error
	}{
		{
			in: eacmodel.Transaction{
				FixedNode: types.FixedNode{
					CompanyID: 1001,
					NodeID:    101,
					ID:        5,
				},
				CurrencyID:  1,
				PostDate:    time2,
				Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
				Type:        transactiontype.JournalEntry,
				CreatedBy:   11,
				Slots: []eacmodel.Slot{
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
							ID:        9,
						},
						CurrencyID:  1,  // original is 1
						AccountID:   31, // original is 31
						Debit:       200,
						Credit:      0,
						Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
							ID:        10,
						},
						CurrencyID:  1,
						AccountID:   32,
						Debit:       1800,
						Credit:      0,
						Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
							ID:        11,
						},
						CurrencyID:  1,
						AccountID:   33,
						Debit:       0,
						Credit:      1001, //original is 1000
						Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					},
					// {
					// 	FixedNode: types.FixedNode{
					// 		CompanyID: 1001,
					// 		NodeID:    101,
					// 	},
					// 	CurrencyID:  1,
					// 	AccountID:   34,
					// 	Debit:       0,
					// 	Credit:      1000,
					// 	Description: helper.StrPointer("new"),
					// 	PostDate:    time2,
					// },
				},
			},
			err: nil,
		},
	}

	for _, v := range samples {
		_, _, err := transactionServ.JournalUpdate(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}
}

func TestOpeningEntryCreate(t *testing.T) {
	_, transactionServ := initTransactionTest()

	time1, err := time.Parse(consts.TimeLayoutZone, "2021-10-18 15:10:00 +0300")

	if err != nil {
		t.Errorf("error in parsing date %v in layout %v", consts.DefaultLimit, "2021-10-18 15:10:00 +0300")
	}

	samples := []struct {
		in  eacmodel.Transaction
		err error
	}{
		{
			in: eacmodel.Transaction{
				FixedNode: types.FixedNode{
					CompanyID: 1001,
					NodeID:    101,
				},
				CurrencyID:  1,
				PostDate:    time1,
				Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
				Type:        transactiontype.OpeningEntry,
				CreatedBy:   11,
				Slots: []eacmodel.Slot{
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   1, // original is 31
						Debit:       1500,
						Credit:      0,
						Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   2,
						Debit:       0,
						Credit:      500,
						Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   3,
						Debit:       0,
						Credit:      1000, //original is 1000
						Description: helper.StrPointer("A dr200, B dr800, C cr1000"),
					},
				},
			},
			err: nil,
		},
	}

	for _, v := range samples {
		_, err := transactionServ.JournalEntry(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}
}
func TestPaymentEntry(t *testing.T) {
	_, transactionServ := initTransactionTest()

	time1, err := time.Parse(consts.TimeLayoutZone, "2019-01-02 15:10:00 +0300")

	if err != nil {
		t.Errorf("error in parsing date %v in layout %v", consts.DefaultLimit, "2020-01-02 15:10:00 +0300")
	}

	samples := []struct {
		in  eacmodel.Transaction
		err error
	}{
		{
			in: eacmodel.Transaction{
				FixedNode: types.FixedNode{
					CompanyID: 1001,
					NodeID:    101,
				},
				CurrencyID:  1,
				PostDate:    time1,
				Description: helper.StrPointer("A cr2000, B dr1000, C dr1000"),
				Type:        transactiontype.PaymentEntry,
				CreatedBy:   11,
				Slots: []eacmodel.Slot{
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   3,
						Debit:       0,
						Credit:      2000,
						Description: helper.StrPointer("A cr2000, B dr1000, C dr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   1,
						Debit:       1000,
						Credit:      0,
						Description: helper.StrPointer("A cr2000, B dr1000, C dr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   2,
						Debit:       1000,
						Credit:      0,
						Description: helper.StrPointer("A cr2000, B dr1000, C dr1000"),
					},
				},
			},
			err: nil,
		},
	}

	for _, v := range samples {
		_, err := transactionServ.JournalEntry(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}
}
func TestReceiptEntry(t *testing.T) {
	_, transactionServ := initTransactionTest()

	time1, err := time.Parse(consts.TimeLayoutZone, "2019-01-02 15:10:00 +0300")

	if err != nil {
		t.Errorf("error in parsing date %v in layout %v", consts.DefaultLimit, "2020-01-02 15:10:00 +0300")
	}

	samples := []struct {
		in  eacmodel.Transaction
		err error
	}{
		{
			in: eacmodel.Transaction{
				FixedNode: types.FixedNode{
					CompanyID: 1001,
					NodeID:    101,
				},
				CurrencyID:  1,
				PostDate:    time1,
				Description: helper.StrPointer("A dr2000, B dr1000, C cr1000"),
				Type:        transactiontype.ReceiptEntry,
				CreatedBy:   11,
				Slots: []eacmodel.Slot{
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   2,
						Debit:       2000,
						Credit:      0,
						Description: helper.StrPointer("A dr2000, B cr1000, C cr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   1,
						Debit:       0,
						Credit:      1000,
						Description: helper.StrPointer("A dr2000, B cr1000, C cr1000"),
					},
					{
						FixedNode: types.FixedNode{
							CompanyID: 1001,
							NodeID:    101,
						},
						AccountID:   3,
						Debit:       0,
						Credit:      1000,
						Description: helper.StrPointer("A dr2000, B cr1000, C cr1000"),
					},
				},
			},
			err: nil,
		},
	}

	for _, v := range samples {
		_, err := transactionServ.JournalEntry(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}
}
