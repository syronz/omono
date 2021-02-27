package service

import (
	"fmt"
	"omono/domain/subscriber/enum/accountstatus"
	"omono/domain/subscriber/submodel"
	"omono/domain/subscriber/subrepo"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
	"strconv"
	"time"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"

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
func (p *SubAccountServ) FindByID(fix types.FixedCol) (account submodel.Account, err error) {
	if account, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1049049", "can't fetch the account", fix.ID)
		return
	}

	if account.Phones, err = p.PhoneServ.AccountsPhones(fix); err != nil {
		err = corerr.Tick(err, "E1017084", "can't fetch the account's phones", fix.ID)
		return
	}

	return
}

// TxFindAccountStatus will return the status of an account
func (p *SubAccountServ) TxFindAccountStatus(db *gorm.DB, fix types.FixedCol) (account submodel.Account, err error) {
	if account, err = p.Repo.TxFindAccountStatus(db, fix); err != nil {
		err = corerr.Tick(err, "E1048403", "can't fetch the account's status", fix.ID)
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

// GetAllAccounts will fetch all of the accounts. Currently used for balancesheet
func (p *SubAccountServ) GetAllAccounts(params param.Param) (accounts []submodel.Account,
	count int64, err error) {

	if accounts, err = p.Repo.GetAllAccounts(params); err != nil {
		glog.CheckError(err, "error in fetching all accounts")
		return
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
func (p *SubAccountServ) Save(account submodel.Account) (savedAccount submodel.Account, err error) {
	return p.TxSave(p.Engine.DB, account)
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
func (p *SubAccountServ) Delete(fix types.FixedCol) (account submodel.Account, err error) {
	if account, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1038835", "account not found for deleting")
		return
	}

	//check if account is only read-only
	if account.ReadOnly {
		//this temp variable will be given to the name of account if exists in english
		var tmp string
		if account.NameEn != nil {
			tmp = *account.NameEn
		}
		err = limberr.New("account has a child", "E1082665").
			Message(corerr.VHasChildThereforeNotDeleted, tmp).
			Custom(corerr.ForeignErr).Build()
		return
	}
	// check child accounts
	// params := param.NewForDelete("bas_accounts", "parent_id", fix.ID)

	// var accounts []submodel.Account
	// if accounts, err = p.Repo.List(params); err != nil {
	// 	err = corerr.Tick(err, "E1036442", "accounts not fetch for delete an account")
	// 	return
	// }

	// if len(accounts) > 0 {
	// 	var tmp string
	// 	if account.NameEn != nil {
	// 		tmp = *account.NameEn
	// 	}
	// 	err = limberr.New("account has a child", "E1082665").
	// 		Message(corerr.VHasChildThereforeNotDeleted, tmp).
	// 		Custom(corerr.ForeignErr).Build()
	// 	return
	// }

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
func (p *SubAccountServ) IsActive(fix types.FixedCol) (bool, submodel.Account, error) {
	var account submodel.Account
	var err error
	if account, err = p.FindByID(fix); err != nil {
		return false, account, corerr.Tick(err, "E1059307", "account not exist", fix.ID)
	}

	return account.Status == accountstatus.Active, account, nil
}

func treeChartOfAccounts(accounts []submodel.Account) (root submodel.Tree) {
	arr := make([]submodel.Tree, len(accounts))

	for i, v := range accounts {
		arr[i].ID = v.ID
		arr[i].ParentID = v.ParentID
		arr[i].Code = v.Code
		if v.NameEn != nil {
			arr[i].Name = *v.NameEn
		}
		arr[i].Type = v.Type
	}

	pMap := make(map[types.RowID]*submodel.Tree, 1)

	pMap[0] = &root

	exceed := submodel.Tree{
		Name: "exceed",
	}

	for i, v := range arr {
		pMap[v.ID] = &arr[i]
	}

	for i, v := range arr {
		pID := parseParent(v.ParentID)

		pMap[pID].Counter++
		if pMap[pID].Counter < consts.MaxChildrenForChartOfAccounts {
			pMap[pID].Children = append(pMap[pID].Children, &arr[i])
		} else {
			if pMap[pID].Counter == consts.MaxChildrenForChartOfAccounts {
				exceed.ParentID = v.ParentID
				pMap[pID].Children = append(pMap[pID].Children, &exceed)
			}
		}

	}

	return
}

func parseParent(pID *types.RowID) types.RowID {
	if pID == nil {
		return 0
	}
	return *pID
}

// ChartOfAccountRefresh is a tree shape of accounts implemented in the nested app
func (p *SubAccountServ) ChartOfAccountRefresh(params param.Param) (root submodel.Tree,
	err error) {

	var accounts []submodel.Account
	params.Limit = consts.MaxRowsCount
	params.Order = "bas_accounts.code ASC"

	if accounts, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in accounts list")
		return
	}

	root = treeChartOfAccounts(accounts)
	now := time.Now()
	root.LastRefresh = &now
	cacheChartOffAccount = &root

	return
}

// ChartOfAccount is a tree shape of accounts implemented in the nested app
func (p *SubAccountServ) ChartOfAccount(params param.Param) (root submodel.Tree,
	err error) {

	params.Limit = consts.MaxRowsCount
	params.Order = "bas_accounts.code ASC"

	if cacheChartOffAccount == nil {
		if root, err = p.ChartOfAccountRefresh(params); err != nil {
			glog.CheckError(err, "error in accounts list")
			return
		}
		return
	}

	return *cacheChartOffAccount, nil
}

// SearchLeafs is used for searching among accounts
func (p *SubAccountServ) SearchLeafs(search string, lang dict.Lang) (accounts []submodel.Account,
	err error) {

	//unfilteredAccs ..
	var unfilteredAcc []submodel.Account

	params := param.New()
	params.PreCondition = "bas_accounts.status = 'active' AND bas_accounts.read_only = 0 AND "
	code, errConvert := strconv.Atoi(search)
	if errConvert != nil {
		params.PreCondition += fmt.Sprintf("bas_accounts.name_%v LIKE '%v%%'", lang, search)
	} else {
		params.PreCondition += fmt.Sprintf("bas_accounts.code LIKE '%v%%'", code)
	}

	params.Order = "bas_accounts.code ASC"

	if unfilteredAcc, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in accounts list")
		return
	}

	//using a loop to filter the inactice accounts
	for _, v := range unfilteredAcc {
		if v.Status == accountstatus.Inactive {
			continue
		}
		accounts = append(accounts, v)
	}

	return
}
