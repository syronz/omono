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

// TempSlotAPI for injecting slot service
type TempSlotAPI struct {
	Service service.EacTempSlotServ
	Engine  *core.Engine
}

// ProvideTempSlotAPI for slot is used in wire
func ProvideTempSlotAPI(c service.EacTempSlotServ) TempSlotAPI {
	return TempSlotAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a slot by it's id
func (p *TempSlotAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error
	var slot eacmodel.Slot
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("slotID"), "E1461257", eacterm.TempSlots); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if slot, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, eacterm.TempSlots).
		JSON(slot)
}

// List of slots
func (p *TempSlotAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacmodel.TempSlotTable, eaccounting.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1489791"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, eacterm.TempSlots).
		JSON(data)
}

// Create slot
func (p *TempSlotAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var slot, createdSlot eacmodel.Slot
	var err error

	if slot.CompanyID, slot.NodeID, err = resp.GetCompanyNode("E1496370", eaccounting.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if slot.CompanyID, err = resp.GetCompanyID("E1447626"); err != nil {
		return
	}

	if !resp.CheckRange(slot.CompanyID) {
		return
	}

	if err = resp.Bind(&slot, "E1434729", eaccounting.Domain, eacterm.TempSlots); err != nil {
		return
	}

	if createdSlot, err = p.Service.Create(slot); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, eacterm.Slot).
		JSON(createdSlot)
}

// Update slot
func (p *TempSlotAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error

	var slot, slotBefore, slotUpdated eacmodel.Slot
	var fix types.FixedCol
	_ = slotBefore

	if fix, err = resp.GetFixedCol(c.Param("slotID"), "E1478308", eacterm.TempSlots); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&slot, "E1450991", eaccounting.Domain, eacterm.TempSlots); err != nil {
		return
	}

	if slotBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	slot.ID = fix.ID
	slot.CompanyID = fix.CompanyID
	slot.NodeID = fix.NodeID
	if slotUpdated, err = p.Service.Save(slot); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, eacterm.TempSlots).
		JSON(slotUpdated)
}

// Excel generate excel files eaced on search
func (p *TempSlotAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.TempSlots, eaccounting.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1460149"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	slots, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("slot")
	ex.AddSheet("Slots").
		AddSheet("Summary").
		Active("Slots").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "F", 15.3).
		SetColWidth("G", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Slots").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
		WriteData(slots).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
