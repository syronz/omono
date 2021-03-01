package service

import (
	"errors"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/test/kernel"
	"testing"

	"gorm.io/gorm"
)

func initRoleTest() (engine *core.Engine, roleServ BasRoleServ) {
	logQuery, debugLevel := initServiceTest()
	engine = kernel.StartMotor(logQuery, debugLevel)
	roleServ = ProvideBasRoleService(basrepo.ProvideRoleRepo(engine))

	return
}

func TestRoleCreate(t *testing.T) {
	_, roleServ := initRoleTest()

	samples := []struct {
		in  basmodel.Role
		err error
	}{
		{
			in: basmodel.Role{
				Name:        "created 1",
				Resources:   string(base.SuperAccess),
				Description: "created 1",
			},
			err: nil,
		},
		{
			in: basmodel.Role{
				Name:        "created 1",
				Resources:   string(base.SuperAccess),
				Description: "created 2",
			},
			err: errors.New("duplicate"),
		},
		{
			in: basmodel.Role{
				Name:      "minimum fields",
				Resources: string(base.SuperAccess),
			},
			err: nil,
		},
		{
			in: basmodel.Role{
				Name:        "long name: big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name big name",
				Resources:   string(base.SuperAccess),
				Description: "created 2",
			},
			err: errors.New("data too long for name"),
		},
		{
			in: basmodel.Role{
				Resources:   string(base.SuperAccess),
				Description: "created 3",
			},
			err: errors.New("name is required"),
		},
	}

	for _, v := range samples {
		_, err := roleServ.Create(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}

}

func TestRoleUpdate(t *testing.T) {
	_, roleServ := initRoleTest()

	samples := []struct {
		in  basmodel.Role
		err error
	}{
		{
			in: basmodel.Role{
				Model: gorm.Model{
					ID: 5,
				},
				Name:        "num 1 update",
				Resources:   string(base.SuperAccess),
				Description: "num 1 update",
			},
			err: nil,
		},
		{
			in: basmodel.Role{
				Model: gorm.Model{
					ID: 6,
				},
				Name:        "num 2 update",
				Description: "num 2 update",
			},
			err: errors.New("resources are required"),
		},
	}

	for _, v := range samples {
		_, _, err := roleServ.Save(v.in)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("ERROR FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.in, err, v.err)
		}
	}

}

func TestRoleDelete(t *testing.T) {
	_, roleServ := initRoleTest()

	samples := []struct {
		id  uint
		err error
	}{
		{
			id:  7,
			err: nil,
		},
		{
			id:  99999999,
			err: errors.New("record not found"),
		},
	}

	for _, v := range samples {
		_, err := roleServ.Delete(v.id)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) {
			t.Errorf("ERROR FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.id, err, v.err)
		}
	}
}

func TestRoleList(t *testing.T) {
	_, roleServ := initRoleTest()
	regularParam := getRegularParam("bas_roles.id asc")
	regularParam.Filter = "description[like]'searchTerm1'"

	samples := []struct {
		params param.Param
		count  int64
		err    error
	}{
		{
			params: param.New(),
			err:    nil,
			count:  11,
		},
		{
			params: regularParam,
			err:    nil,
			count:  3,
		},
	}

	for _, v := range samples {
		_, count, err := roleServ.List(v.params)

		if (v.err == nil && err != nil) || (v.err != nil && err == nil) || count != v.count {
			t.Errorf("FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", v.params, count, v.count)
		}
	}
}

func TestRoleExcel(t *testing.T) {
	_, roleServ := initRoleTest()
	regularParam := getRegularParam("bas_roles.id asc")

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
		data, err := roleServ.Excel(v.params)
		if (v.err == nil && err != nil) || (v.err != nil && err == nil) || int64(len(data)) < v.count {
			t.Errorf("FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v::: \nErr :::%+v:::",
				v.params, int64(len(data)), v.count, err)
		}
	}
}
