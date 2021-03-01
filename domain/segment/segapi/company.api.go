package segapi

import (
	"net/http"
	"omono/domain/base/basterm"
	"omono/domain/segment"
	"omono/domain/segment/segmodel"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/pkg/helper/excel"

	"github.com/gin-gonic/gin"
)

// CompanyAPI for injecting company service
type CompanyAPI struct {
	Service service.SegCompanyServ
	Engine  *core.Engine
}

// ProvideCompanyAPI for company is used in wire
func ProvideCompanyAPI(c service.SegCompanyServ) CompanyAPI {
	return CompanyAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a company by it's id
func (p *CompanyAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, segment.Domain)
	var err error
	var company segmodel.Company
	var id uint

	if id, err = resp.GetID(c.Param("companyID"), "E1070061", basterm.Company); err != nil {
		return
	}

	if company, err = p.Service.FindByID(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(segment.ViewCompany)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Company).
		JSON(company)
}

// List of companies
func (p *CompanyAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, segmodel.CompanyTable, segment.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(segment.ListCompany)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Companies).
		JSON(data)
}

// Create company
func (p *CompanyAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, segment.Domain)
	var company, createdCompany segmodel.Company
	var err error

	if err = resp.Bind(&company, "E1057541", segment.Domain, basterm.Company); err != nil {
		return
	}

	if createdCompany, err = p.Service.Create(company); err != nil {
		resp.Error(err).JSON()
		return
	}

	// resp.RecordCreate(segment.CreateCompany, company)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.Company).
		JSON(createdCompany)
}

// Update company
func (p *CompanyAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, segment.Domain)
	var err error

	var company, companyBefore, companyUpdated segmodel.Company
	var id uint

	if id, err = resp.GetID(c.Param("companyID"), "E1076703", basterm.Company); err != nil {
		return
	}

	if err = resp.Bind(&company, "E1086162", segment.Domain, basterm.Company); err != nil {
		return
	}

	company.ID = id
	company.CreatedAt = companyBefore.CreatedAt
	if companyUpdated, companyBefore, err = p.Service.Save(company); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(segment.UpdateCompany, companyBefore, company)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, basterm.Company).
		JSON(companyUpdated)
}

// Delete company
func (p *CompanyAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, segment.Domain)
	var err error
	var company segmodel.Company
	var id uint

	if id, err = resp.GetID(c.Param("companyID"), "E1092196", basterm.Company); err != nil {
		return
	}

	if company, err = p.Service.Delete(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(segment.DeleteCompany, company)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.Company).
		JSON()
}

// Excel generate excel files based on search
func (p *CompanyAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Companies, segment.Domain)
	var err error

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
		SetColWidth("B", "G", 15.3).
		SetColWidth("H", "H", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Companies").
		WriteHeader("ID", "Name", "Code", "Type", "Status", "Updated At").
		SetSheetFields("ID", "Name", "Code", "Type", "Status", "UpdatedAt").
		WriteData(companies).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(segment.ExcelCompany)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
