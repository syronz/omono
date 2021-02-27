package basapi

import (
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/message/basterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// CityAPI for injecting city service
type CityAPI struct {
	Service service.BasCityServ
	Engine  *core.Engine
}

// ProvideCityAPI for city is used in wire
func ProvideCityAPI(c service.BasCityServ) CityAPI {
	return CityAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a city by it's id
func (p *CityAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var city basmodel.City
	var fix types.FixedCol

	if fix.ID, err = types.StrToRowID(c.Param("cityID")); err != nil {
		return
	}

	if city, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewCity)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.City).
		JSON(city)
}

// List of cities
func (p *CityAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basmodel.CityTable, base.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1066610"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListCity)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Cities).
		JSON(data)
}

// Create city
func (p *CityAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var city, createdCity basmodel.City
	var err error

	if err = resp.Bind(&city, "E1025334", base.Domain, basterm.City); err != nil {
		return
	}

	if createdCity, err = p.Service.Create(city); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(base.CreateCity, city)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.City).
		JSON(createdCity)
}

// Update city
func (p *CityAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error

	var city, cityBefore, cityUpdated basmodel.City
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("cityID"), "E1064608", basterm.City); err != nil {
		return
	}

	if err = resp.Bind(&city, "E1093884", base.Domain, basterm.City); err != nil {
		return
	}

	if cityBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	city.ID = fix.ID
	city.CreatedAt = cityBefore.CreatedAt
	if cityUpdated, err = p.Service.Save(city); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.UpdateCity, cityBefore, city)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, basterm.City).
		JSON(cityUpdated)
}

// Delete city
func (p *CityAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var city basmodel.City
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("cityID"), "E1057845", basterm.City); err != nil {
		return
	}

	if city, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.DeleteCity, city)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.City).
		JSON()
}

// Excel generate excel files based on search
func (p *CityAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Cities, base.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1074098"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	cities, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("city")
	ex.AddSheet("Cities").
		AddSheet("Summary").
		Active("Cities").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "G", 15.3).
		SetColWidth("H", "H", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Cities").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Code", "Type", "Status", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Code", "Type", "Status", "UpdatedAt").
		WriteData(cities).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ExcelCity)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
