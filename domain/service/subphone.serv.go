package service

import (
	"fmt"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"

	"github.com/syronz/limberr"

	"gorm.io/gorm"
)

// SubPhoneServ for injecting auth subrepo
type SubPhoneServ struct {
	Repo   subrepo.PhoneRepo
	Engine *core.Engine
}

// ProvideSubPhoneService for phone is used in wire
func ProvideSubPhoneService(p subrepo.PhoneRepo) SubPhoneServ {
	return SubPhoneServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting phone by it's id
func (p *SubPhoneServ) FindByID(fix types.FixedCol) (phone submodel.Phone, err error) {
	if phone, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1057387", "can't fetch the phone", fix.ID)
		return
	}

	return
}

// FindByPhone for getting phone by it's id
func (p *SubPhoneServ) FindByPhone(phoneNumber string) (phone submodel.Phone, err error) {
	if phone, err = p.Repo.FindByPhone(phoneNumber); err != nil {
		// do not log error if it is not-found
		if limberr.GetCustom(err) != corerr.NotFoundErr {
			err = corerr.Tick(err, "E1048291", "can't fetch the phone by phone-number", phoneNumber)
		}
		return
	}

	return
}

// AccountsPhones return list of phones assigned to an account
func (p *SubPhoneServ) AccountsPhones(fix types.FixedCol) (phones []submodel.Phone, err error) {
	if phones, err = p.Repo.AccountsPhones(fix); err != nil {
		err = corerr.Tick(err, "E1067138", "can't get account's phone", fix)
		return
	}

	return
}

// List of phones, it support pagination and search and return back count
func (p *SubPhoneServ) List(params param.Param) (phones []submodel.Phone,
	count int64, err error) {

	if phones, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in phones list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in phones count")
	}

	return
}

// Create a phone
func (p *SubPhoneServ) Create(phone submodel.Phone) (createdPhone submodel.Phone, err error) {
	return p.TxCreate(p.Repo.Engine.DB, phone)

}

// TxCreate used in case of transaction activated
func (p *SubPhoneServ) TxCreate(db *gorm.DB, phone submodel.Phone) (createdPhone submodel.Phone, err error) {
	if err = phone.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1067746", "validation failed in creating the phone", phone)
		return
	}

	var phoneExist submodel.Phone

	phoneExist, err = p.FindByPhone(phone.Phone)

	var account submodel.Account
	account.ID = phone.AccountID

	switch {

	//not found
	case limberr.GetCustom(err) == corerr.NotFoundErr:
		if createdPhone, err = p.Repo.TxCreate(db, phone); err != nil {
			err = corerr.Tick(err, "E1091571", "phone not created", phone)
			return
		}
		phone = createdPhone

		// database error
	case err != nil:
		err = corerr.Tick(err, "E1064472", "can't fetch the phone by phone-number", phone.Phone)
		return

		// found
	default:
		phone = phoneExist

	}

	// var joiner submodel.AccountPhone
	if _, err = p.Repo.JoinAccountPhone(db, account, phone, phone.Default); err != nil {
		err = corerr.Tick(err, "E1062524", "phone-join not created", phone)
		return
	}

	return
}

// Save a phone, if it is exist update it, if not create it
func (p *SubPhoneServ) Save(phone submodel.Phone) (savedPhone submodel.Phone, err error) {
	if err = phone.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1031666", corerr.ValidationFailed, phone)
		return
	}

	if savedPhone, err = p.Repo.Save(phone); err != nil {
		err = corerr.Tick(err, "E1031295", "phone not saved")
		return
	}

	return
}

// Delete phone, it is soft delete
func (p *SubPhoneServ) Delete(fix types.FixedCol) (phone submodel.Phone, err error) {
	if phone, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1044187", "phone not found for deleting")
		return
	}

	if err = p.Repo.Delete(phone); err != nil {
		err = corerr.Tick(err, "E1032085", "phone not deleted")
		return
	}

	return
}

// Separate phone, it is soft delete
func (p *SubPhoneServ) Separate(fix types.FixedCol) (aPhone submodel.AccountPhone, err error) {
	if aPhone, err = p.Repo.FindAccountPhoneByID(fix); err != nil {
		err = corerr.Tick(err, "E1049677", "account-phone not found for deleting")
		return
	}

	if err = p.Repo.SeparateAccountPhone(aPhone); err != nil {
		err = corerr.Tick(err, "E1040009", "account-phone not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *SubPhoneServ) Excel(params param.Param) (phones []submodel.Phone, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", submodel.PhoneTable)

	if phones, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1066621", "cant generate the excel list for phones")
		return
	}

	return
}
