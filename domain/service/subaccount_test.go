package service

import (
	"omono/domain/subscriber/enum/accountstatus"
	"omono/domain/subscriber/enum/accounttype"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/pkg/helper"
	"omono/test/kernel"
	"testing"

	"gorm.io/gorm"
)

func initAccountTest() (engine *core.Engine, accountServ SubAccountServ) {
	logQuery, debugLevel := initServiceTest()
	engine = kernel.StartMotor(logQuery, debugLevel)

	phoneServ := ProvideSubPhoneService(subrepo.ProvidePhoneRepo(engine))
	accountServ = ProvideSubAccountService(subrepo.ProvideAccountRepo(engine), phoneServ)

	return
}

func TestAccountCreate(t *testing.T) {

	_, accountServ := initAccountTest()

	samples := []struct {
		in  submodel.Account
		err error
	}{
		{
			in: submodel.Account{
				NameEn: "child 1 of asset ",
				NameKu: helper.StrPointer("3"),
				Type:   accounttype.VIP,
				Status: accountstatus.Active,
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

	var id uint = 21

	if _, err := testAccountServ.Delete(id); err != nil {
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
				Model: gorm.Model{
					ID: 31,
				},
				NameEn: "updated VIP Account",
				Type:   accounttype.VIP,
				Credit: 500000,
			},
			err: nil,
		},
	}
	for _, v := range samples {
		_, _, err := accountServ.Save(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}
}
