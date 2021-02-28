package service

import (
	"errors"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/test/kernel"
	"testing"

	"gorm.io/gorm"
)

func initSettingTest() (engine *core.Engine, settingServ BasSettingServ) {
	logQuery, debugLevel := initServiceTest()
	engine = kernel.StartMotor(logQuery, debugLevel)
	settingServ = ProvideBasSettingService(basrepo.ProvideSettingRepo(engine))

	return
}

func TestUpdateSetting(t *testing.T) {
	_, settingServ := initSettingTest()

	samples := []struct {
		in  basmodel.Setting
		err error
	}{
		{
			in: basmodel.Setting{
				Model: gorm.Model{
					ID: 20,
				},
				Property:    "num 1 updated",
				Value:       "num 1 updated",
				Type:        "num 1 updated",
				Description: "num 1 updated",
			},
			err: nil,
		},
		{
			in: basmodel.Setting{
				Model: gorm.Model{
					ID: 21,
				},
				Value:       "num 2 updated",
				Type:        "num 2 updated",
				Description: "num 2 updated",
			},
			err: errors.New("property is required"),
		},
	}

	for _, v := range samples {
		_, _, err := settingServ.Save(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("ERROR FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}

}

func TestListSetting(t *testing.T) {
	_, settingServ := initSettingTest()
	regularParam := getRegularParam("bas_settings.id asc")
	// regularParam.Search = "searchTerm1"
	regularParam.Filter = "description[like]'%searchTerm1%'"

	samples := []struct {
		params param.Param
		count  int64
		err    error
	}{
		{
			params: param.Param{},
			err:    nil,
			count:  6,
		},
		{
			params: regularParam,
			err:    nil,
			count:  3,
		},
	}

	for _, v := range samples {
		_, count, err := settingServ.List(v.params)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) || count != v.count {
			t.Errorf("FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.params, count, v.count)
		}
	}
}

func TestSettingExcel(t *testing.T) {
	_, settingServ := initSettingTest()
	regularParam := getRegularParam("bas_settings.id asc")

	samples := []struct {
		params param.Param
		count  int64
		err    error
	}{
		{
			params: regularParam,
			err:    nil,
			count:  6,
		},
	}

	for _, v := range samples {
		data, err := settingServ.Excel(v.params)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) || int64(len(data)) < v.count {
			t.Errorf("FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v::: \nErr :::%+v:::",
				v.params, int64(len(data)), v.count, err)
		}
	}
}
