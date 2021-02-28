package service

import (
	"encoding/json"
	"fmt"
	"omono/domain/subscriber/enum/accountstatus"
	"omono/domain/subscriber/enum/accounttype"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/pkg/helper"
	"omono/test/kernel"
	"testing"

	"github.com/syronz/dict"
	"gorm.io/gorm"
)

func initAccountTest() (engine *core.Engine, accountServ SubAccountServ) {
	logQuery, debugLevel := initServiceTest()
	engine = kernel.StartMotor(logQuery, debugLevel)

	phoneServ := ProvideSubPhoneService(subrepo.ProvidePhoneRepo(engine))
	accountServ = ProvideSubAccountService(subrepo.ProvideAccountRepo(engine), phoneServ)

	return
}
func TestTreeChartOfAccounts(t *testing.T) {
	accounts := []submodel.Account{
		{
			gorm.Model: gorm.Model{
				ID: 1,
			},
			Code:   "1",
			NameEn: "regular",
			Type:   accounttype.VIP,
		},
		{
			gorm.Model: gorm.Model{
				ID: 2,
			},
			ParentID: uintPointer(1),
			Code:     "11",
			NameEn:   "Regular USD",
			Type:     accounttype.Regular,
		},
		{
			gorm.Model: gorm.Model{
				ID: 3,
			},
			ParentID: uintPointer(1),
			Code:     "12",
			NameEn:   "Regular IQD",
			Type:     accounttype.Regular,
		},
		{
			gorm.Model: gorm.Model{
				ID: 4,
			},
			Code:   "3",
			NameEn: "Business",
			Type:   accounttype.Business,
		},
		{
			gorm.Model: gorm.Model{
				ID: 5,
			},
			ParentID: uintPointer(4),
			Code:     "31",
			NameEn:   "Building",
			Type:     accounttype.Business,
		},
		{
			gorm.Model: gorm.Model{
				ID: 6,
			},
			ParentID: uintPointer(1),
			Code:     "311",
			NameEn:   "HQ",
			Type:     accounttype.Business,
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
		in  submodel.Account
		err error
	}{
		{
			in: submodel.Account{
				Code:     "12135",
				NameEn:   "child 1 of asset ",
				NameKu:   helper.StrPointer("3"),
				Type:     accounttype.VIP,
				Status:   accountstatus.Active,
				ParentID: uintPointer(1),
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

	sample := gorm.Model{
		ID: 21,
	}

	if _, err := testAccountServ.Delete(sample); err != nil {
		t.Errorf("there is an error for deleting the account, %v", err.Error())
	}
}

func TestAccountList(t *testing.T) {
	_, testAccountServ := initAccountTest()

	regularParam := getRegularParam("bas_accounts.id")
	regularParam.Filter = "name_en[like]'regular'"

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
		in  submodel.Account
		err error
	}{
		{
			in: submodel.Account{
				gorm.Model: gorm.Model{
					ID: 31,
				},
				Code:    "1231",
				NameEn:  "updated Partner Account",
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

	//data := make(map[string]interface{})

	var root submodel.Tree
	root, err = testAccountServ.ChartOfAccount(params)

	if err != nil {
		t.Errorf("there is an error for fetching chart of accounts, %v", err.Error())

	}

	fmt.Println(root)

	// for _, v := range root {
	// 	fmt.Println(v)
	// }
}
