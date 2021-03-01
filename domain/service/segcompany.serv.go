package service

import (
	"fmt"
	"omono/domain/segment/segmodel"
	"omono/domain/segment/segrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// SegCompanyServ for injecting auth segrepo
type SegCompanyServ struct {
	Repo   segrepo.CompanyRepo
	Engine *core.Engine
}

// ProvideSegCompanyService for company is used in wire
func ProvideSegCompanyService(p segrepo.CompanyRepo) SegCompanyServ {
	return SegCompanyServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting company by it's id
func (p *SegCompanyServ) FindByID(id uint) (company segmodel.Company, err error) {
	if company, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E1049049", "can't fetch the company", id)
		return
	}

	return
}

// TxFindCompanyStatus will return the status of an company
func (p *SegCompanyServ) TxFindCompanyStatus(db *gorm.DB, id uint) (company segmodel.Company, err error) {
	if company, err = p.Repo.TxFindCompanyStatus(db, id); err != nil {
		err = corerr.Tick(err, "E1048403", "can't fetch the company's status", id)
		return
	}

	return
}

// List of companies, it support pagination and search and return back count
func (p *SegCompanyServ) List(params param.Param) (companies []segmodel.Company,
	count int64, err error) {

	if companies, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in companies list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in companies count")
	}

	return
}

// Create a company
func (p *SegCompanyServ) Create(company segmodel.Company) (createdCompany segmodel.Company, err error) {
	db := p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"users table"), "rollback recover create user")
			db.Rollback()
		}
	}()

	if createdCompany, err = p.TxCreate(p.Repo.Engine.DB, company); err != nil {
		err = corerr.Tick(err, "E1014394", "error in creating company for user", createdCompany)

		db.Rollback()
		return
	}

	db.Commit()

	return
}

// TxCreate is used for creating an company in case of transaction activated
func (p *SegCompanyServ) TxCreate(db *gorm.DB, company segmodel.Company) (createdCompany segmodel.Company, err error) {
	if err = company.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1076780", "validation failed in creating the company", company)
		return
	}

	if createdCompany, err = p.Repo.TxCreate(db, company); err != nil {
		err = corerr.Tick(err, "E1065508", "company not created", company)
		return
	}

	return
}

// Save a company, if it is exist update it, if not create it
func (p *SegCompanyServ) Save(company segmodel.Company) (savedCompany, companyBefore segmodel.Company, err error) {
	if companyBefore, err = p.FindByID(company.ID); err != nil {
		err = corerr.Tick(err, "E1073641", "company not exist")
		return
	}

	company.CreatedAt = companyBefore.CreatedAt

	savedCompany, err = p.TxSave(p.Engine.DB, company)
	return
}

// TxSave a company, if it is exist update it, if not create it
func (p *SegCompanyServ) TxSave(db *gorm.DB, company segmodel.Company) (savedCompany segmodel.Company, err error) {
	if err = company.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1064761", corerr.ValidationFailed, company)
		return
	}

	if savedCompany, err = p.Repo.TxSave(db, company); err != nil {
		err = corerr.Tick(err, "E1084087", "company not saved")
		return
	}

	return
}

// Delete company, it is soft delete
func (p *SegCompanyServ) Delete(id uint) (company segmodel.Company, err error) {
	if company, err = p.FindByID(id); err != nil {
		err = corerr.Tick(err, "E1038835", "company not found for deleting")
		return
	}

	if err = p.Repo.Delete(company); err != nil {
		err = corerr.Tick(err, "E1045410", "company not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *SegCompanyServ) Excel(params param.Param) (companies []segmodel.Company, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", segmodel.CompanyTable)

	if companies, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1023076", "cant generate the excel list for companies")
		return
	}

	return
}
