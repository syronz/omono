package service

import (
	"fmt"
	"omono/domain/subscriber/enum/accountstatus"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// SubAccountServ for injecting auth subrepo
type SubAccountServ struct {
	Repo      subrepo.AccountRepo
	Engine    *core.Engine
	PhoneServ SubPhoneServ
}

var cacheChartOffAccount *submodel.Tree

// ProvideSubAccountService for account is used in wire
func ProvideSubAccountService(p subrepo.AccountRepo, phoneServ SubPhoneServ) SubAccountServ {
	return SubAccountServ{
		Repo:      p,
		Engine:    p.Engine,
		PhoneServ: phoneServ,
	}
}

// FindByID for getting account by it's id
func (p *SubAccountServ) FindByID(id uint) (account submodel.Account, err error) {
	if account, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E1049049", "can't fetch the account", id)
		return
	}

	if account.Phones, err = p.PhoneServ.AccountsPhones(id); err != nil {
		err = corerr.Tick(err, "E1017084", "can't fetch the account's phones", id)
		return
	}

	return
}

// TxFindAccountStatus will return the status of an account
func (p *SubAccountServ) TxFindAccountStatus(db *gorm.DB, id uint) (account submodel.Account, err error) {
	if account, err = p.Repo.TxFindAccountStatus(db, id); err != nil {
		err = corerr.Tick(err, "E1048403", "can't fetch the account's status", id)
		return
	}

	return
}

// List of accounts, it support pagination and search and return back count
func (p *SubAccountServ) List(params param.Param) (accounts []submodel.Account,
	count int64, err error) {

	if accounts, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in accounts list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in accounts count")
	}

	return
}

// Create a account
func (p *SubAccountServ) Create(account submodel.Account) (createdAccount submodel.Account, err error) {
	db := p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"users table"), "rollback recover create user")
			db.Rollback()
		}
	}()

	if createdAccount, err = p.TxCreate(p.Repo.Engine.DB, account); err != nil {
		err = corerr.Tick(err, "E1014394", "error in creating account for user", createdAccount)

		db.Rollback()
		return
	}

	db.Commit()

	return
}

// TxCreate is used for creating an account in case of transaction activated
func (p *SubAccountServ) TxCreate(db *gorm.DB, account submodel.Account) (createdAccount submodel.Account, err error) {
	if err = account.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1076780", "validation failed in creating the account", account)
		return
	}

	if createdAccount, err = p.Repo.TxCreate(db, account); err != nil {
		err = corerr.Tick(err, "E1065508", "account not created", account)
		return
	}

	for _, phone := range account.Phones {
		phone.AccountID = createdAccount.ID
		if _, err = p.PhoneServ.TxCreate(db, phone); err != nil {
			err = corerr.Tick(err, "E1040913", "error in creating phone for account", phone)

			return
		}
	}

	return
}

// Save a account, if it is exist update it, if not create it
func (p *SubAccountServ) Save(account submodel.Account) (savedAccount, accountBefore submodel.Account, err error) {
	if accountBefore, err = p.FindByID(account.ID); err != nil {
		err = corerr.Tick(err, "E1073641", "account not exist")
		return
	}

	account.CreatedAt = accountBefore.CreatedAt

	savedAccount, err = p.TxSave(p.Engine.DB, account)
	return
}

// TxSave a account, if it is exist update it, if not create it
func (p *SubAccountServ) TxSave(db *gorm.DB, account submodel.Account) (savedAccount submodel.Account, err error) {
	if err = account.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1064761", corerr.ValidationFailed, account)
		return
	}

	if savedAccount, err = p.Repo.TxSave(db, account); err != nil {
		err = corerr.Tick(err, "E1084087", "account not saved")
		return
	}

	return
}

// Delete account, it is soft delete
func (p *SubAccountServ) Delete(id uint) (account submodel.Account, err error) {
	if account, err = p.FindByID(id); err != nil {
		err = corerr.Tick(err, "E1038835", "account not found for deleting")
		return
	}

	if err = p.Repo.Delete(account); err != nil {
		err = corerr.Tick(err, "E1045410", "account not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *SubAccountServ) Excel(params param.Param) (accounts []submodel.Account, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", submodel.AccountTable)

	if accounts, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1023076", "cant generate the excel list for accounts")
		return
	}

	return
}

// IsActive check the status of an account
func (p *SubAccountServ) IsActive(id uint) (bool, submodel.Account, error) {
	var account submodel.Account
	var err error
	if account, err = p.FindByID(id); err != nil {
		return false, account, corerr.Tick(err, "E1059307", "account not exist", id)
	}

	return account.Status == accountstatus.Active, account, nil
}
