package basrepo

import (
	"omono/domain/base/basmodel"
	"omono/domain/base/message/basterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
	"github.com/syronz/limberr"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// AccountRepo for injecting engine
type AccountRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideAccountRepo is used in wire and initiate the Cols
func ProvideAccountRepo(engine *core.Engine) AccountRepo {
	return AccountRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(basmodel.Account{}), basmodel.AccountTable),
	}
}

// FindByID finds the account via its id
func (p *AccountRepo) FindByID(fix types.FixedNode) (account basmodel.Account, err error) {
	err = p.Engine.ReadDB.Table(basmodel.AccountTable).
		Where("id = ? AND company_id = ? AND node_id = ? AND bas_accounts.deleted_at is null", fix.ID.ToUint64(), fix.CompanyID, fix.NodeID).
		First(&account).Error

	account.ID = fix.ID
	err = p.dbError(err, "E1045869", account, corterm.List)

	return
}

// TxFindAccountStatus finds the account via its id and return back the status
func (p *AccountRepo) TxFindAccountStatus(db *gorm.DB, fix types.FixedNode) (account basmodel.Account, err error) {
	// err = db.Clauses(clause.Locking{Strength: "UPDATE"}).Table(basmodel.AccountTable).
	err = db.Table(basmodel.AccountTable).
		Where("id = ? AND company_id = ? AND bas_accounts.deleted_at is null", fix.ID.ToUint64(), fix.CompanyID).
		First(&account).Error

	account.ID = fix.ID
	err = p.dbError(err, "E1042082", account, corterm.List)

	return
}

// List returns an array of accounts
func (p *AccountRepo) List(params param.Param) (accounts []basmodel.Account, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1050070").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1084619").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.AccountTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&accounts).Error

	err = p.dbError(err, "E1082445", basmodel.Account{}, corterm.List)

	return
}

//GetAllAccounts will fetch all accounts with specified companyID
func (p *AccountRepo) GetAllAccounts(params param.Param) (account []basmodel.Account, err error) {
	err = p.Engine.ReadDB.Table(basmodel.AccountTable).
		Where("company_id = ?", params.CompanyID).Find(&account).Error

	err = p.dbError(err, "E1011232", basmodel.Account{}, corterm.List)

	return
}

// Count of accounts, mainly calls with List
func (p *AccountRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1037218").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.AccountTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1056203", basmodel.Account{}, corterm.List)
	return
}

// TxSave the account, in case it is not exist create it
func (p *AccountRepo) TxSave(db *gorm.DB, account basmodel.Account) (u basmodel.Account, err error) {
	if err = db.Table(basmodel.AccountTable).Save(&account).Error; err != nil {
		err = p.dbError(err, "E1070874", account, corterm.Updated)
	}

	db.Table(basmodel.AccountTable).Where("id = ?", account.ID).Find(&u)
	return
}

// Create a account
func (p *AccountRepo) Create(account basmodel.Account) (u basmodel.Account, err error) {
	if err = p.Engine.DB.Table(basmodel.AccountTable).Create(&account).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1054044", account, corterm.Created)
	}
	return
}

// TxCreate a account
func (p *AccountRepo) TxCreate(db *gorm.DB, account basmodel.Account) (u basmodel.Account, err error) {
	if err = db.Table(basmodel.AccountTable).Create(&account).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1054044", account, corterm.Created)
	}
	return
}

// Delete the account
func (p *AccountRepo) Delete(account basmodel.Account) (err error) {
	now := time.Now()
	account.DeletedAt = &now
	if err = p.Engine.DB.Table(basmodel.AccountTable).Save(&account).Error; err != nil {
		err = p.dbError(err, "E1095299", account, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper database error
func (p *AccountRepo) dbError(err error, code string, account basmodel.Account, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, account.ID, basterm.Accounts)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(basterm.Account), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(basterm.Account), account.NameEn).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, account.NameEn)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
