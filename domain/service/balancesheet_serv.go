package service

import (
	"errors"
	"fmt"
	"omono/cmd/restapi/enum/settingfields"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/internal/types"
	"strconv"
)

// EacBalanceSheetServ for injecting auth eacrepo
type EacBalanceSheetServ struct {
	Repo   eacrepo.BalanceSheetRepo
	Engine *core.Engine
}

// ProvideEacBalanceSheetService for transaction is used in wire
func ProvideEacBalanceSheetService(p eacrepo.BalanceSheetRepo) EacBalanceSheetServ {
	return EacBalanceSheetServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// BalanceSheet will get the balancesheet for each main account based on the level
func (p *EacBalanceSheetServ) BalanceSheet(params param.Param, level string) (balanceSheet []eacmodel.BalanceSheet, err error) {
	var depth uint64
	var mainAccountID types.RowID
	var settingField string
	var accounts []basmodel.Account
	var balances []eacmodel.Balance

	fmt.Println("the level", level)
	if depth, err = strconv.ParseUint(level, 10, 64); err != nil {
		err = errors.New("the level of depth for balancesheet is invalid")
		return
	}

	phoneServ := ProvideBasPhoneService(basrepo.ProvidePhoneRepo(p.Engine))
	accountServ := ProvideBasAccountService(basrepo.ProvideAccountRepo(p.Engine), phoneServ)
	//fetching all the accounts
	if accounts, _, err = accountServ.GetAllAccounts(params); err != nil {
		return
	}

	if balances, err = p.Repo.GetAllAccountBalances(params); err != nil {
		return
	}

	for i := 0; i < 3; i++ {
		var AccountType types.Setting
		var mainAccount eacmodel.BalanceSheet

		switch i {
		case 0:
			AccountType = settingfields.MainAssetID
		case 1:
			AccountType = settingfields.MainLiabilityID
		case 2:
			AccountType = settingfields.MainEquityID
		}
		//obtaining mainAccountID from setting
		settingField = p.Engine.Setting[AccountType].Value
		mainAccountID, _ = types.StrToRowID(settingField)

		//fetch account names for the main accounts
		mainAccount.NameEn, mainAccount.NameKu, mainAccount.NameAr = fetchNamesMain(mainAccountID, accounts)

		//fetching the balance and childs for main account
		if mainAccount.ChildBalances, mainAccount.Balance, err = p.FetchBalances(params, mainAccountID, depth, accounts, balances); err != nil {
			err = errors.New("cannot fetch balances for assets")
			return
		}

		//appending to the final balancesheet
		balanceSheet = append(balanceSheet, mainAccount)

	}

	return
}

// FetchBalances will get the balancesheet for each main account based on the level
func (p *EacBalanceSheetServ) FetchBalances(params param.Param, parentAccount types.RowID, depth uint64, accounts []basmodel.Account, balances []eacmodel.Balance) (balanceSheet []eacmodel.BalanceSheet, balance float64, err error) {

	var balanceAlreadyExist bool
	var tempBalanceSheet eacmodel.BalanceSheet
	if depth == 0 {
		return
	}
	//search among childs

	for _, v := range accounts {

		balanceAlreadyExist = false
		if v.ParentID == nil {
			continue
		}
		pid := v.ParentID.ToUint64()
		fmt.Println("the parent id is following :", pid)
		if *v.ParentID == parentAccount {
			if v.NameEn != nil {
				tempBalanceSheet.NameEn = *v.NameEn

			}
			if v.NameKu != nil {
				tempBalanceSheet.NameKu = *v.NameKu

			}
			if v.NameAr != nil {
				tempBalanceSheet.NameAr = *v.NameAr

			}

			fmt.Println("work fine")

			//check if the current account already has balance
			for _, b := range balances {
				if b.AccountID == v.ID {
					tempBalanceSheet.Balance = b.Balance
					balanceSheet = append(balanceSheet, tempBalanceSheet)
					balance += tempBalanceSheet.Balance
					balanceAlreadyExist = true
					break
				}

			}

			if balanceAlreadyExist {
				continue
			}

			tempBalanceSheet.ChildBalances, tempBalanceSheet.Balance, _ = p.FetchBalances(params, v.ID, depth-1, accounts, balances)

			//now we fetch the balance for the last childs
			if depth == 0 {
				for _, b := range balances {
					if b.AccountID == v.ID {
						tempBalanceSheet.Balance = b.Balance
						balanceSheet = append(balanceSheet, tempBalanceSheet)
						balance += tempBalanceSheet.Balance

					}

				}
			}
			balanceSheet = append(balanceSheet, tempBalanceSheet)
			balance += tempBalanceSheet.Balance
		}

	}
	return
}

func fetchNamesMain(mainAccount types.RowID, accounts []basmodel.Account) (nameEn, nameKu, nameAr string) {

	for _, v := range accounts {
		if v.ID == mainAccount {
			if v.NameEn != nil {
				nameEn = *v.NameEn
			}
			if v.NameKu != nil {
				nameKu = *v.NameKu
			}
			if v.NameAr != nil {
				nameAr = *v.NameAr
			}
			return
		}
	}
	return
}
