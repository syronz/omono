package service

import (
	"errors"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/test/kernel"
	"testing"

	"gorm.io/gorm"
)

func initPhoneTest() (engine *core.Engine, phoneService SubPhoneServ) {
	queryLog, debugLevel := initServiceTest()

	engine = kernel.StartMotor(queryLog, debugLevel)

	phoneService = ProvideSubPhoneService(subrepo.ProvidePhoneRepo(engine))

	return

}

func TestPhoneCreate(test *testing.T) {
	//we will call the initPhoneTest for starting the generating the engine special for TDD
	//then we fetch the phone service which included the phone repo
	//the engine is skipped

	_, phoneService := initPhoneTest()

	// we create a struct of phone model along with the error
	//then we treat each element of the struct as a test and pass it to the system for test.

	//First test element has no issue and should return NO ERRORS at all.
	//2nd test element has error because the input for phone is more than 8 digits
	//3rd test: ERROR: b/c input for phone is less than 5
	//4th test: ERROR: input for Notes is greater than 255 characters
	testCollector := []struct {
		phone submodel.Phone
		err   error
	}{
		{
			phone: submodel.Phone{
				Phone:     "077022222",
				Notes:     "This phone number has been created",
				AccountID: 1,
			},
			err: nil,
		},
		{
			phone: submodel.Phone{
				Phone:     "07702232133123213213",
				Notes:     "This phone number has been created",
				AccountID: 1,
			},
			err: errors.New("this phone has length more than 8 digits"),
		},

		{
			phone: submodel.Phone{
				Phone:     "077",
				Notes:     "this phone  number has been created",
				AccountID: 1,
			},
			err: errors.New("phone has less than 5 digits"),
		},

		{
			phone: submodel.Phone{
				Phone: "321332131",
				Notes: "This phone has been created, This phone has been created, This phone has been created, This phone has been created,This phone has been created, This phone has been created, This phone has been created, This phone has been created, This phone has been created,",
			},
			err: errors.New("The length of notes is greater than 255"),
		},
	}

	for _, value := range testCollector {
		_, err := phoneService.Create(value.phone)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.phone, err, value.err)
		}

	}
}

func TestPhoneSave(test *testing.T) {
	//the engine is skipped
	_, phoneService := initPhoneTest()

	type err error
	collector := []struct {
		phone submodel.Phone
		err   error
	}{
		{
			phone: submodel.Phone{
				Model: gorm.Model{
					ID: 1,
				},
				Phone: "23134142",
				Notes: "phone has been updated",
			},
			err: nil,
		},
		{
			phone: submodel.Phone{
				Model: gorm.Model{
					ID: 1314421,
				},
				Phone: "3131233",
				Notes: "phone has been updated",
			},
			err: errors.New("Phone doesn't exist"),
		},
		{
			phone: submodel.Phone{
				Model: gorm.Model{
					ID: 1,
				},
				Notes: "phone has been updated",
			},
			err: errors.New("Phone is required"),
		},
	}

	for _, value := range collector {

		_, err := phoneService.Save(value.phone)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.phone, err, value.err)
		}
	}
}

//TestPhoneUpdate() Commented Out(why?)
/*
	because service/update() accepts idNode arg
/*
func TestPhoneUpdate(test *testing.T) {
	//the engine is skipped
	_, phoneService := initPhoneTest()

	type err error
	collector := []struct {
		phone submodel.Phone
		err   error
	}{
		{
			phone: submodel.Phone{

				ID:        1,
				Phone:     "23134142",
				Notes:     "phone has been updated",
			},
			err: nil,
		},
		{
			phone: submodel.Phone{
				ID:        1314421,
				Phone:     "3131233",
				Notes:     "phone has been updated",
			},
			err: errors.New("Phone doesn't exist"),
		},
		{
			phone: submodel.Phone{
				ID:        1,
				Notes:     "phone has been updated",
			},
			err: errors.New("Phone is required"),
		},
	}

	for _, value := range collector {

		_, err := phoneService.Save(value.phone)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			test.Errorf("\nERROR FOR :::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.phone, err, value.err)
		}
	}
}
*/

//TestPhoneDelete() Commented Out(why?)
/*
	because service/Delete() accepts idNode arg
/*
/*
func TestPhoneDelete(t *testing.T) {
	_, phoneService := initPhoneTest()

	testCollector := []struct {
		phone submodel.Phone
		err   error
	}{
		{
			phone: submodel.Phone{
				ID:        1,
			},
			err: nil,
		},
		{
			phone: submodel.Phone{
				ID:        1111111,
			},
			err: errors.New("phone was not found to be deleted"),
		},
	}

	for _, value := range testCollector {
		_, err := phoneService.Delete(value)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			t.Errorf("ERROR FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.id, err, value.err)
		}
	}
}

*/
//TestPhoneFindByID() Commented Out(why?)
/*
	because service/FindByID() accepts idNode arg
/*
/*
func TestPhoneFindByID(t *testing.T) {
	_, phoneService := initPhoneTest()

	testCollector := []struct {
		phone submodel.Phone
		err   error
	}{
		{
			phone: submodel.Phone{
				ID:        1,
			},
			err: nil,
		},
		{
			phone: submodel.Phone{
				ID:        1324231,
			err: errors.New("there is no phone record"),
		},
	}

	for _, value := range testCollector {
		_, err := phoneService.Delete(value)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			t.Errorf("ERROR FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.id, err, value.err)
		}
	}
}
*/

func TestPhoneFindByPhone(t *testing.T) {
	_, phoneService := initPhoneTest()

	testCollector := []struct {
		phone submodel.Phone
		err   error
	}{
		{
			phone: submodel.Phone{
				Phone: "07701001111",
			},
			err: nil,
		},
		{
			phone: submodel.Phone{
				Phone: "12345678",
			},
			err: errors.New("there is no phone record"),
		},
	}

	for _, value := range testCollector {
		returnedPhone, err := phoneService.FindByPhone(value.phone.Phone)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			t.Errorf("ERROR FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.phone.Phone, returnedPhone.Phone, value.err)
		}
	}
}

func TestPhoneList(t *testing.T) {
	_, phoneService := initPhoneTest()
	regularParam := getRegularParam("base_phone.id asc")
	regularParam.Filter = "Note[like]'original'"
	testCollector := []struct {
		params param.Param
		count  uint64
		err    error
	}{
		{
			params: param.Param{},
			count:  4,
			err:    nil,
		},
		{
			params: regularParam,
			err:    nil,
			count:  1,
		},
	}

	for _, value := range testCollector {
		_, count, err := phoneService.List(value.params)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) {
			t.Errorf("ERROR FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v:::", value.params, count, value.count)
		}
	}
}

func TestPhoneExcel(t *testing.T) {
	_, phoneService := initPhoneTest()
	regularParam := getRegularParam("bas_phone.id asc")

	testCollector := []struct {
		params param.Param
		count  uint64
		err    error
	}{
		{
			params: param.Param{},
			count:  3,
			err:    nil,
		},
		{
			params: regularParam,
			err:    nil,
			count:  1,
		},
	}

	for _, value := range testCollector {
		returnedPhone, err := phoneService.Excel(value.params)
		if (value.err == nil && err != nil) || (value.err != nil && err == nil) || uint64(len(returnedPhone)) < value.count {
			t.Errorf("FOR ::::%+v::: \nRETURNS :::%+v:::, \nIT SHOULD BE :::%+v::: \nErr :::%+v:::",
				value.params, uint64(len(returnedPhone)), value.count, err)
		}
	}
}
