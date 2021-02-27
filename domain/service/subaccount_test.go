package service

import (
	"encoding/json"
	"fmt"
	"github.com/syronz/dict"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/base/enum/accountstatus"
	"omono/domain/base/enum/accounttype"
	"omono/internal/core"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/helper"
	"omono/test/kernel"
	"testing"
)

func initAccountTest() (engine *core.Engine, accountServ SubAccountServ) {
	logQuery, debugLevel := initServiceTest()
	engine = kernel.StartMotor(logQuery, debugLevel)

	phoneServ := ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountServ = ProvideSubAccountService(basrepo.ProvideAccountRepo(engine), phoneServ)

	return
}
func TestTreeChartOfAccounts(t *testing.T) {
	accounts := []basmodel.Account{
		{
			FixedCol: types.FixedCol{
				ID: 1,
			},
			Code:   "1",
			NameEn: helper.StrPointer("Asset"),
			Type:   accounttype.Asset,
		},
		{
			FixedCol: types.FixedCol{
				ID: 2,
			},
			ParentID: types.RowIDPointer(1),
			Code:     "11",
			NameEn:   helper.StrPointer("Cash USD"),
			Type:     accounttype.Cash,
		},
		{
			FixedCol: types.FixedCol{
				ID: 3,
			},
			ParentID: types.RowIDPointer(1),
			Code:     "12",
			NameEn:   helper.StrPointer("Cash IQD"),
			Type:     accounttype.Cash,
		},
		{
			FixedCol: types.FixedCol{
				ID: 4,
			},
			Code:   "3",
			NameEn: helper.StrPointer("Expense"),
			Type:   accounttype.Expense,
		},
		{
			FixedCol: types.FixedCol{
				ID: 5,
			},
			ParentID: types.RowIDPointer(4),
			Code:     "31",
			NameEn:   helper.StrPointer("Building"),
			Type:     accounttype.Expense,
		},
		{
			FixedCol: types.FixedCol{
				ID: 6,
			},
			ParentID: types.RowIDPointer(1),
			Code:     "311",
			NameEn:   helper.StrPointer("HQ"),
			Type:     accounttype.Expense,
		},
	}

	for _, v := range accounts {
		fmt.Println(v.NameEn, v.ID)
	}

	root := treeChartOfAccounts(accounts)

	b, err := json.MarshalIndent(root, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

}

func TestAccountCreate(t *testing.T) {

	_, accountServ := initAccountTest()

	samples := []struct {
		in  basmodel.Account
		err error
	}{
		{
			in: basmodel.Account{
				FixedCol: types.FixedCol{
					CompanyID: 1001,
					NodeID:    101,
				},
				Code:     "12135",
				NameEn:   helper.StrPointer("child 1 of asset "),
				NameKu:   helper.StrPointer("3"),
				Type:     accounttype.Asset,
				Status:   accountstatus.Active,
				ParentID: types.RowIDPointer(1),
			},
			err: nil,
		},
	}
	for _, v := range samples {
		_, err := accountServ.Create(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}
}

func TestAccountDelete(t *testing.T) {
	_, testAccountServ := initAccountTest()

	sample := types.FixedCol{
		ID:        21,
		CompanyID: 1001,
		NodeID:    101,
	}

	if _, err := testAccountServ.Delete(sample); err != nil {
		t.Errorf("there is an error for deleting the account, %v", err.Error())
	}
}

func TestAccountList(t *testing.T) {
	_, testAccountServ := initAccountTest()

	regularParam := getRegularParam("bas_accounts.id")
	regularParam.Filter = "name_en[like]'Asset'"

	collection := []struct {
		params param.Param
		count  int64
		err    error
	}{
		{
			params: param.Param{},
			err:    nil,
			count:  3,
		},
		{
			params: regularParam,
			err:    nil,
			count:  0,
		},
	}

	for _, value := range collection {
		_, count, err := testAccountServ.List(value.params)

		if (value.err == nil && err != nil) || (value.err != nil && err == nil) || count != value.count {
			t.Errorf("FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.params, count, value.count)
		}
	}

}

func TestAccountUpdate(t *testing.T) {

	_, accountServ := initAccountTest()
	samples := []struct {
		in  basmodel.Account
		err error
	}{
		{
			in: basmodel.Account{
				FixedCol: types.FixedCol{
					CompanyID: 1001,
					NodeID:    101,
					ID:        31,
				},
				Code:    "1231",
				NameEn:  helper.StrPointer("updated Partner Account"),
				NameAr:  helper.StrPointer("بە سەرکەوتوی ئەپدەیت کرا"),
				Type:    accounttype.Partner,
				Balance: 500000,
			},
			err: nil,
		},
	}
	for _, v := range samples {
		_, err := accountServ.Save(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}
}

func TestSearchLeafs(t *testing.T) {
	_, testAccountServ := initAccountTest()

	search := "181"

	value, err := testAccountServ.SearchLeafs(search, dict.En)
	if err != nil {
		t.Errorf("there is an error for searching the account(s), %v", err.Error())

	}

	if len(value) == 0 {
		fmt.Println("no results")
	}
	fmt.Println(value)
	for _, v := range value {
		fmt.Println(v)
	}
}

func TestChartOfAccounts(t *testing.T) {
	_, testAccountServ := initAccountTest()

	params := param.New()

	var err error
	params.Select = "bas_accounts.id,bas_accounts.parent_id,bas_accounts.code,bas_accounts.name_ar,bas_accounts.name_en,bas_accounts.name_ku,bas_accounts.type"

	params.PreCondition = fmt.Sprintf("bas_accounts.type != '%v'", accounttype.Customer)

	params.CompanyID = 1001

	//data := make(map[string]interface{})

	var root basmodel.Tree
	root, err = testAccountServ.ChartOfAccount(params)

	if err != nil {
		t.Errorf("there is an error for fetching chart of accounts, %v", err.Error())

	}

	fmt.Println(root)

	// for _, v := range root {
	// 	fmt.Println(v)
	// }
}
