package service

import (
	"errors"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/internal/types"
	"omono/test/kernel"
	"testing"
)

func initUserTest() (engine *core.Engine, userService BasUserServ) {
	queryLog, debugLevel := initServiceTest()

	engine = kernel.StartMotor(queryLog, debugLevel)

	userService = ProvideBasUserService(basrepo.ProvideUserRepo(engine))

	return

}

func TestUserCreate(test *testing.T) {
	//the engine is skipped
	_, userService := initUserTest()

	collector := []struct {
		user basmodel.User
		err  error
	}{
		{
			user: basmodel.User{
				gorm.Model: gorm.Model{
					CompanyID: 1001,
					NodeID:    101,
				},
				RoleID:   2,
				Username: "tester",
				Password: "21312349807709",
				Name:     "Alan wake",
				Lang:     "en",
				Email:    "",
				Phone:    "",
			},
			err: nil,
		},
		{
			user: basmodel.User{
				gorm.Model: gorm.Model{
					CompanyID: 1001,
					NodeID:    101,
				},
				RoleID:   3,
				Username: "tester",
				Password: "21312",
				Lang:     "en",
				Email:    "",
				Phone:    "",
			},
			err: errors.New("property is less than 8 characters"),
		},

		{
			user: basmodel.User{
				gorm.Model: gorm.Model{
					CompanyID: 1001,
					NodeID:    101,
				},
				RoleID:   4,
				Username: "",
				Password: "1111111111111111",
				Lang:     "en",
				Email:    "",
				Phone:    "",
			},
			err: errors.New("username is empty"),
		},

		{
			user: basmodel.User{
				gorm.Model: gorm.Model{
					CompanyID: 1001,
					NodeID:    101,
				},
				RoleID:   0,
				Username: "tester",
				Password: "1111111111111111",
				Lang:     "en",
				Email:    "",
				Phone:    "",
			},
			err: errors.New("Role is invalid"),
		},
		{
			user: basmodel.User{
				gorm.Model: gorm.Model{
					CompanyID: 1001,
					NodeID:    101,
				},
				RoleID:   2,
				Username: "tester",
				Password: "1111111111111111",
				Lang:     "fa",
				Email:    "",
				Phone:    "",
			},
			err: errors.New("Language is not accepted"),
		},
		/*
			{
				user: basmodel.User{
					gorm.Model: gorm.Model{
						CompanyID: 1001,
						NodeID:    101,
					},
					RoleID:   2,
					Username: "tester",
					Password: "1111111111111111",
					Lang:     "en",
					Email:    "aran@aran.com",
					Phone:    "",
				},
				err: errors.New("email is not verified"),
			},
		*/
	}

	for _, value := range collector {
		_, err := userService.Create(value.user)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.user, err, value.err)
		}

	}
}

func TestUserUpdate(test *testing.T) {
	//the engine is skipped
	_, userService := initUserTest()

	type err error
	collector := []struct {
		user basmodel.User
		err  error
	}{
		{
			user: basmodel.User{
				gorm.Model: gorm.Model{
					ID:        11,
					CompanyID: 1001,
					NodeID:    101,
				},
				RoleID:   1,
				Username: "updated",
				Password: "32131323132",
				Lang:     "ku",
				Email:    "test@test.com",
				Phone:    "updated",
			},
			err: nil,
		},
		{
			user: basmodel.User{
				gorm.Model: gorm.Model{
					ID:        11,
					CompanyID: 1001,
					NodeID:    101,
				},
				RoleID:   3,
				Username: "updated ",
				Password: "32131323132",
				Email:    "Updated",
				Phone:    "",
			},
			err: errors.New("language is required"),
		},
	}

	for _, value := range collector {

		_, err := userService.Save(value.user)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.user, err, value.err)
		}
	}
}

//Test for delete
//notice for deletion we just take the ided columns
//the service/user.Delete() func only accepts the ided columnss
func TestUserDelete(test *testing.T) {
	//the engine is skipped
	_, userService := initUserTest()
	type err error
	collector := []struct {
		user gorm.Model
		err  error
	}{
		{
			user: gorm.Model{
				ID:        12,
				CompanyID: 1001,
				NodeID:    101,
			},
			err: nil,
		},
		{
			user: gorm.Model{
				ID:        2525252,
				CompanyID: 1001,
				NodeID:    101,
			},
			err: errors.New("Record was not found for deletion"),
		},
	}

	for _, value := range collector {
		_, err := userService.Delete(value.user)
		test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.user.ID, err, value.err)

	}
}

func TestUserExcel(test *testing.T) {
	//the engine is skipped
	_, userService := initUserTest()
	regularParam := getRegularParam("bas_users.id asc")

	collector := []struct {
		params param.Param
		count  int64
		err    error
	}{
		{
			params: regularParam,
			err:    nil,
			count:  3,
		},
	}

	for _, value := range collector {
		users, err := userService.Excel(value.params)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) || int64(len(users)) < value.count {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v::: \nErr :::%+v:::", value.params, int64(len(users)), value.count, err)
		}
	}
}

func TestUserList(t *testing.T) {
	_, userService := initUserTest()
	regularParam := getRegularParam("bas_users.id asc")
	regularParam.Filter = "username[like]'searchTerm1'"

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
		_, count, err := userService.List(value.params)

		if (value.err == nil && err != nil) || (value.err != nil && err == nil) || count != value.count {
			t.Errorf("FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.params, count, value.count)
		}
	}
}
func TestUserFindByID(test *testing.T) {
	//the engine is skipped
	_, userService := initUserTest()
	type err error
	collector := []struct {
		user gorm.Model
		err  error
	}{
		{
			user: gorm.Model{
				ID:        2,
				CompanyID: 1001,
				NodeID:    101,
			},
			err: nil,
		},
		{
			user: gorm.Model{
				ID:        32131312,
				CompanyID: 1001,
				NodeID:    101,
			},
			err: errors.New("User was not found"),
		},
	}

	for _, value := range collector {
		user, err := userService.FindByID(value.user)
		if value.err == nil && err != nil {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.user, user.ID, value.user.ID)
		}

	}
}

func TestUserFindByUsername(test *testing.T) {
	//the engine is skipped
	_, userService := initUserTest()
	type err error
	collector := []struct {
		id       gorm.Model
		username string
		err      error
	}{
		{
			id: gorm.Model{
				ID:        11,
				CompanyID: 1001,
				NodeID:    101,
			},
			username: "admin",
			err:      nil,
		},
		{
			id: gorm.Model{
				ID:        0,
				CompanyID: 1001,
				NodeID:    101,
			},
			username: "unknownUser",
			err:      nil,
		},
	}

	for _, value := range collector {
		user, err := userService.FindByUsername(value.username)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.username, user.ID, value.id)
		}

	}
}
