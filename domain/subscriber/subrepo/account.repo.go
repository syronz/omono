package subrepo

import (
	"omono/domain/base/basterm"
	"omono/domain/subscriber/submodel"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/pkg/helper"
	"reflect"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"

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
		Cols:   helper.TagExtracter(reflect.TypeOf(submodel.Account{}), submodel.AccountTable),
	}
}

// FindByID finds the account via its id
func (p *AccountRepo) FindByID(id uint) (account submodel.Account, err error) {
	err = p.Engine.ReadDB.Table(submodel.AccountTable).
		Where("id = ? AND sub_accounts.deleted_at is null", id).
		First(&account).Error

	account.ID = id
	err = p.dbError(err, "E1045869", account, corterm.List)

	return
}

// TxFindAccountStatus finds the account via its id and return back the status
func (p *AccountRepo) TxFindAccountStatus(db *gorm.DB, id uint) (account submodel.Account, err error) {
	// err = db.Clauses(clause.Locking{Strength: "UPDATE"}).Table(submodel.AccountTable).
	err = db.Table(submodel.AccountTable).
		Where("id = ? AND sub_accounts.deleted_at is null", id).
		First(&account).Error

	account.ID = id
	err = p.dbError(err, "E1042082", account, corterm.List)

	return
}

// List returns an array of accounts
func (p *AccountRepo) List(params param.Param) (accounts []submodel.Account, err error) {
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

	err = p.Engine.ReadDB.Table(submodel.AccountTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&accounts).Error

	err = p.dbError(err, "E1082445", submodel.Account{}, corterm.List)

	return
}

//GetAllAccounts will fetch all accounts with specified companyID
func (p *AccountRepo) GetAllAccounts(params param.Param) (account []submodel.Account, err error) {
	err = p.Engine.ReadDB.Table(submodel.AccountTable).Find(&account).Error

	err = p.dbError(err, "E1011232", submodel.Account{}, corterm.List)

	return
}

// Count of accounts, mainly calls with List
func (p *AccountRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1037218").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(submodel.AccountTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1056203", submodel.Account{}, corterm.List)
	return
}

// TxSave the account, in case it is not exist create it
func (p *AccountRepo) TxSave(db *gorm.DB, account submodel.Account) (u submodel.Account, err error) {
	if err = db.Table(submodel.AccountTable).Save(&account).Error; err != nil {
		err = p.dbError(err, "E1070874", account, corterm.Updated)
	}

	db.Table(submodel.AccountTable).Where("id = ?", account.ID).Find(&u)
	return
}

// Create a account
func (p *AccountRepo) Create(account submodel.Account) (u submodel.Account, err error) {
	if err = p.Engine.DB.Table(submodel.AccountTable).Create(&account).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1054044", account, corterm.Created)
	}
	return
}

// TxCreate a account
func (p *AccountRepo) TxCreate(db *gorm.DB, account submodel.Account) (u submodel.Account, err error) {
	if err = db.Table(submodel.AccountTable).Create(&account).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1054044", account, corterm.Created)
	}
	return
}

// Delete the account
func (p *AccountRepo) Delete(account submodel.Account) (err error) {
	if err = p.Engine.DB.Unscoped().Table(submodel.AccountTable).Delete(&account).Error; err != nil {
		err = p.dbError(err, "E1095299", account, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper database error
func (p *AccountRepo) dbError(err error, code string, account submodel.Account, action string) error {
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
