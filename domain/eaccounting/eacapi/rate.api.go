package eacapi

import (
	"net/http"
	"omono/domain/eaccounting"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// RateAPI for injecting rate service
type RateAPI struct {
	Service service.EacRateServ
	Engine  *core.Engine
}

// ProvideRateAPI for rate is used in wire
func ProvideRateAPI(c service.EacRateServ) RateAPI {
	return RateAPI{Service: c, Engine: c.Engine}
}

// RatesInCity return all currency and its rate for a city
func (p *RateAPI) RatesInCity(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error
	var rates []eacmodel.Rate
	var cityID types.RowID

	if cityID, err = types.StrToRowID(c.Param("cityID")); err != nil {
		resp.Error(err).JSON()
		return
	}

	if rates, err = p.Service.RatesInCity(cityID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, eacterm.Rates).
		JSON(rates)
}

// FindByID is used for fetch a rate by it's id
func (p *RateAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error
	var rate eacmodel.Rate
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("rateID"), "E1491903", eacterm.Rate); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if rate, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ViewRate)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, eacterm.Rate).
		JSON(rate)
}

// List of rates
func (p *RateAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacmodel.RateTable, eaccounting.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1449627"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ListRate)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, eacterm.Rates).
		JSON(data)
}

// Create rate
func (p *RateAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var rate, createdRate eacmodel.Rate
	var err error

	if rate.CompanyID, rate.NodeID, err = resp.GetCompanyNode("E1474003", eaccounting.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if rate.CompanyID, err = resp.GetCompanyID("E1420791"); err != nil {
		return
	}

	if !resp.CheckRange(rate.CompanyID) {
		return
	}

	if err = resp.Bind(&rate, "E1447895", eaccounting.Domain, eacterm.Rate); err != nil {
		return
	}

	if createdRate, err = p.Service.Create(rate); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(eaccounting.CreateRate, rate)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, eacterm.Rate).
		JSON(createdRate)
}

// Update rate
func (p *RateAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error

	var rate, rateBefore, rateUpdated eacmodel.Rate
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("rateID"), "E1417444", eacterm.Rate); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&rate, "E1466762", eaccounting.Domain, eacterm.Rate); err != nil {
		return
	}

	if rateBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	rate.ID = fix.ID
	rate.CompanyID = fix.CompanyID
	rate.NodeID = fix.NodeID
	rate.CreatedAt = rateBefore.CreatedAt
	if rateUpdated, err = p.Service.Save(rate); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.UpdateRate, rateBefore, rate)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, eacterm.Rate).
		JSON(rateUpdated)
}

// Delete rate
func (p *RateAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error
	var rate eacmodel.Rate
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("rateID"), "E1486606", eacterm.Rate); err != nil {
		return
	}

	if rate, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.DeleteRate, rate)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, eacterm.Rate).
		JSON()
}

// Excel generate excel files eaced on search
func (p *RateAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.Rates, eaccounting.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1464408"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	rates, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("rate")
	ex.AddSheet("Rates").
		AddSheet("Summary").
		Active("Rates").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "F", 15.3).
		SetColWidth("G", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Rates").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
		WriteData(rates).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ExcelRate)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
