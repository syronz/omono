package locapi

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"net/http"
	"omono/domain/base/message/basterm"
	"omono/domain/location"
	"omono/domain/location/locmodel"
	"omono/domain/location/locterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// StoreAPI for injecting store service
type StoreAPI struct {
	Service service.LocStoreServ
	Engine  *core.Engine
}

// ProvideStoreAPI for store is used in wire
func ProvideStoreAPI(c service.LocStoreServ) StoreAPI {
	return StoreAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a store by it's id
func (p *StoreAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, location.Domain)
	var err error
	var store locmodel.Store
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("storeID"), "E2891845", locterm.Store); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if store, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(location.ViewStore)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, locterm.Store).
		JSON(store)
}

// List of stores
func (p *StoreAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, locmodel.StoreTable, location.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E2826488"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(location.ListStore)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, locterm.Stores).
		JSON(data)
}

// Create store
func (p *StoreAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, location.Domain)
	var store, createdStore locmodel.Store
	var err error

	if store.CompanyID, store.NodeID, err = resp.GetCompanyNode("E2843418", location.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if store.CompanyID, err = resp.GetCompanyID("E2892450"); err != nil {
		return
	}

	if !resp.CheckRange(store.CompanyID) {
		return
	}

	if err = resp.Bind(&store, "E2854172", location.Domain, locterm.Store); err != nil {
		return
	}

	if createdStore, err = p.Service.Create(store); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(location.CreateStore, store)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, locterm.Store).
		JSON(createdStore)
}

// AddUser create a user for store
func (p *StoreAPI) AddUser(c *gin.Context) {
	resp := response.New(p.Engine, c, location.Domain)
	var err error
	var storeUser locmodel.StoreUser

	if storeUser.CompanyID, storeUser.NodeID, err =
		resp.GetCompanyNode("E2833571", location.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if storeUser.StoreID, err = types.StrToRowID(c.Param("storeID")); err != nil {
		resp.Error(err).JSON()
		return
	}

	if err = resp.Bind(&storeUser, "E2844897", location.Domain, locterm.Store); err != nil {
		return
	}

	if _, err = p.Service.AddUser(storeUser); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(location.AddUser, storeUser)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.User).
		JSON()
}

// DelUser Delete a user for a store
func (p *StoreAPI) DelUser(c *gin.Context) {
	resp := response.New(p.Engine, c, location.Domain)

	var err error
	var ptID types.RowID
	if ptID, err = types.StrToRowID(c.Param("storeUserID")); err != nil {
		err = limberr.Take(err, "E2890212").
			Message(corerr.InvalidVForV, dict.R(corterm.ID), "user_id").
			Custom(corerr.ValidationFailedErr).Build()
		resp.Error(err).JSON()
		return
	}

	if err = p.Service.DelUser(ptID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(location.DelUser, ptID)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.User).
		JSON()
}

// Update store
func (p *StoreAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, location.Domain)
	var err error

	var store, storeBefore, storeUpdated locmodel.Store
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("storeID"), "E2859912", locterm.Store); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&store, "E2869195", location.Domain, locterm.Store); err != nil {
		return
	}

	if storeBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	store.ID = fix.ID
	store.CompanyID = fix.CompanyID
	store.NodeID = fix.NodeID
	store.CreatedAt = storeBefore.CreatedAt
	if storeUpdated, err = p.Service.Save(store); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(location.UpdateStore, storeBefore, store)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, locterm.Store).
		JSON(storeUpdated)
}

// Delete store
func (p *StoreAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, location.Domain)
	var err error
	var store locmodel.Store
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("storeID"), "E2839408", locterm.Store); err != nil {
		return
	}

	if store, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(location.DeleteStore, store)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, locterm.Store).
		JSON()
}

// Excel generate excel files eaced on search
func (p *StoreAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, locterm.Stores, location.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E2846062"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	stores, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("store")
	ex.AddSheet("Stores").
		AddSheet("Summary").
		Active("Stores").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "F", 15.3).
		SetColWidth("G", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Stores").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
		WriteData(stores).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(location.ExcelStore)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
