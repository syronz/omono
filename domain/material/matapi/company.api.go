package matapi

import (
	"net/http"
	"omono/domain/base/message/basterm"
	"omono/domain/material"
	"omono/domain/material/matmodel"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// CompanyAPI for injecting company service
type CompanyAPI struct {
	Service service.MatCompanyServ
	Engine  *core.Engine
}

// ProvideCompanyAPI for company is used in wire
func ProvideCompanyAPI(c service.MatCompanyServ) CompanyAPI {
	return CompanyAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a company by it's id
func (p *CompanyAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error
	var company matmodel.Company
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("compID"), "E7142676", basterm.Company); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if company, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ViewCompany)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Company).
		JSON(company)
}

// List of companies
func (p *CompanyAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, matmodel.CompanyTable, material.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7137019"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ListCompany)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Companies).
		JSON(data)
}

// Create company
func (p *CompanyAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var company, createdCompany matmodel.Company
	var err error

	if company.CompanyID, company.NodeID, err = resp.GetCompanyNode("E7148684", material.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if company.CompanyID, err = resp.GetCompanyID("E7128449"); err != nil {
		return
	}

	if !resp.CheckRange(company.CompanyID) {
		return
	}

	if err = resp.Bind(&company, "E7129072", material.Domain, basterm.Company); err != nil {
		return
	}

	if createdCompany, err = p.Service.Create(company); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(material.CreateCompany, company)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.Company).
		JSON(createdCompany)
}

// Update company
func (p *CompanyAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error

	var company, companyBefore, companyUpdated matmodel.Company
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("compID"), "E7163288", basterm.Company); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&company, "E7147246", material.Domain, basterm.Company); err != nil {
		return
	}

	if companyBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	company.ID = fix.ID
	company.CompanyID = fix.CompanyID
	company.NodeID = fix.NodeID
	company.CreatedAt = companyBefore.CreatedAt
	if companyUpdated, err = p.Service.Save(company); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.UpdateCompany, companyBefore, company)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, basterm.Company).
		JSON(companyUpdated)
}

// Delete company
func (p *CompanyAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error
	var company matmodel.Company
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("compID"), "E7168083", basterm.Company); err != nil {
		return
	}

	if company, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.DeleteCompany, company)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.Company).
		JSON()
}

// Excel generate excel files eaced on search
func (p *CompanyAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Companies, material.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7166407"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	companies, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("company")
	ex.AddSheet("Companies").
		AddSheet("Summary").
		Active("Companies").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "F", 15.3).
		SetColWidth("G", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Companies").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
		WriteData(companies).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ExcelCompany)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
