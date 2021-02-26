package bilapi

import (
	"net/http"
	"omono/domain/bill"
	"omono/domain/bill/bilmodel"
	"omono/domain/bill/bilterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// InvoiceAPI for injecting invoice service
type InvoiceAPI struct {
	Service service.BilInvoiceServ
	Engine  *core.Engine
}

// ProvideInvoiceAPI for invoice is used in wire
func ProvideInvoiceAPI(c service.BilInvoiceServ) InvoiceAPI {
	return InvoiceAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a invoice by it's id
func (p *InvoiceAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, bill.Domain)
	var err error
	var invoice bilmodel.Invoice
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("invoiceID"), "E7763443", bilterm.Invoice); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if invoice, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(bill.ViewInvoice)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, bilterm.Invoice).
		JSON(invoice)
}

// List of invoices
func (p *InvoiceAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, bilmodel.InvoiceTable, bill.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7779804"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(bill.ListInvoice)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, bilterm.Invoices).
		JSON(data)
}

// Create invoice
func (p *InvoiceAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, bilterm.Invoices, bill.Domain)
	var invoice, createdInvoice bilmodel.Invoice
	var err error

	if invoice.CompanyID, invoice.NodeID, err = resp.GetCompanyNode("E7751963", bill.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if invoice.CompanyID, err = resp.GetCompanyID("E7742882"); err != nil {
		return
	}

	if !resp.CheckRange(invoice.CompanyID) {
		return
	}

	if err = resp.Bind(&invoice, "E7761572", bill.Domain, bilterm.Invoice); err != nil {
		return
	}

	invoice.CreatedBy = params.UserID

	if createdInvoice, err = p.Service.Create(invoice); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(bill.CreateInvoice, invoice)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, bilterm.Invoice).
		JSON(createdInvoice)
}

// Update invoice
/*
func (p *InvoiceAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, bill.Domain)
	var err error

	var invoice, invoiceBefore, invoiceUpdated bilmodel.Invoice
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("invoiceID"), "E7720855", bilterm.Invoice); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&invoice, "E7780593", bill.Domain, bilterm.Invoice); err != nil {
		return
	}

	if invoiceBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	invoice.ID = fix.ID
	invoice.CompanyID = fix.CompanyID
	invoice.NodeID = fix.NodeID
	invoice.CreatedAt = invoiceBefore.CreatedAt
	if invoiceUpdated, err = p.Service.Save(invoice); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(bill.UpdateInvoice, invoiceBefore, invoice)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, bilterm.Invoice).
		JSON(invoiceUpdated)
}
*/

// Delete invoice
func (p *InvoiceAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, bill.Domain)
	var err error
	var invoice bilmodel.Invoice
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("invoiceID"), "E7745492", bilterm.Invoice); err != nil {
		return
	}

	if invoice, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(bill.DeleteInvoice, invoice)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, bilterm.Invoice).
		JSON()
}

// Excel generate excel files eaced on search
func (p *InvoiceAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, bilterm.Invoices, bill.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7738430"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	invoices, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("invoice")
	ex.AddSheet("Invoices").
		AddSheet("Summary").
		Active("Invoices").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "F", 15.3).
		SetColWidth("G", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Invoices").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
		WriteData(invoices).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(bill.ExcelInvoice)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
