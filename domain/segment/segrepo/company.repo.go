package segrepo

import (
	"omono/domain/base/basterm"
	"omono/domain/segment/segmodel"
	"omono/domain/segment/segterm"
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

// CompanyRepo for injecting engine
type CompanyRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideCompanyRepo is used in wire and initiate the Cols
func ProvideCompanyRepo(engine *core.Engine) CompanyRepo {
	return CompanyRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(segmodel.Company{}), segmodel.CompanyTable),
	}
}

// FindByID finds the company via its id
func (p *CompanyRepo) FindByID(id uint) (company segmodel.Company, err error) {
	err = p.Engine.ReadDB.Table(segmodel.CompanyTable).
		Where("id = ? AND sub_companies.deleted_at is null", id).
		First(&company).Error

	company.ID = id
	err = p.dbError(err, "E1045869", company, corterm.List)

	return
}

// TxFindCompanyStatus finds the company via its id and return back the status
func (p *CompanyRepo) TxFindCompanyStatus(db *gorm.DB, id uint) (company segmodel.Company, err error) {
	// err = db.Clauses(clause.Locking{Strength: "UPDATE"}).Table(segmodel.CompanyTable).
	err = db.Table(segmodel.CompanyTable).
		Where("id = ? AND sub_companies.deleted_at is null", id).
		First(&company).Error

	company.ID = id
	err = p.dbError(err, "E1042082", company, corterm.List)

	return
}

// List returns an array of companies
func (p *CompanyRepo) List(params param.Param) (companies []segmodel.Company, err error) {
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

	err = p.Engine.ReadDB.Table(segmodel.CompanyTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&companies).Error

	err = p.dbError(err, "E1082445", segmodel.Company{}, corterm.List)

	return
}

// Count of companies, mainly calls with List
func (p *CompanyRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1037218").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(segmodel.CompanyTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1056203", segmodel.Company{}, corterm.List)
	return
}

// TxSave the company, in case it is not exist create it
func (p *CompanyRepo) TxSave(db *gorm.DB, company segmodel.Company) (u segmodel.Company, err error) {
	if err = db.Table(segmodel.CompanyTable).Save(&company).Error; err != nil {
		err = p.dbError(err, "E1070874", company, corterm.Updated)
	}

	db.Table(segmodel.CompanyTable).Where("id = ?", company.ID).Find(&u)
	return
}

// Create a company
func (p *CompanyRepo) Create(company segmodel.Company) (u segmodel.Company, err error) {
	if err = p.Engine.DB.Table(segmodel.CompanyTable).Create(&company).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1054044", company, corterm.Created)
	}
	return
}

// TxCreate a company
func (p *CompanyRepo) TxCreate(db *gorm.DB, company segmodel.Company) (u segmodel.Company, err error) {
	if err = db.Table(segmodel.CompanyTable).Create(&company).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1054044", company, corterm.Created)
	}
	return
}

// Delete the company
func (p *CompanyRepo) Delete(company segmodel.Company) (err error) {
	if err = p.Engine.DB.Unscoped().Table(segmodel.CompanyTable).Delete(&company).Error; err != nil {
		err = p.dbError(err, "E1095299", company, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper database error
func (p *CompanyRepo) dbError(err error, code string, company segmodel.Company, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, company.ID, segterm.Companies)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(segterm.Company), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(segterm.Company), company.Name).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, company.Name)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
