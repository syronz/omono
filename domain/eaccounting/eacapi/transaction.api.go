package eacapi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"omono/cmd/restapi/enum/settingfields"
	"omono/domain/base/basmodel"
	"omono/domain/eaccounting"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/eaccounting/enum/transactiontype"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	numtoword "github.com/syronz/numberprinter"
)

// TransactionAPI for injecting transaction service
type TransactionAPI struct {
	Service service.EacTransactionServ
	Engine  *core.Engine
}

// ProvideTransactionAPI for transaction is used in wire
func ProvideTransactionAPI(c service.EacTransactionServ) TransactionAPI {
	return TransactionAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a transaction by it's id
func (p *TransactionAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error
	var transaction eacmodel.Transaction
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1423147", eacterm.Transaction); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if transaction, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ViewTransaction)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, eacterm.Transaction).
		JSON(transaction)
}

//LastYearCounterByType for getting last tranacstion  year counter plus 1. the given input is the type of the transaction
func (p *TransactionAPI) LastYearCounterByType(c *gin.Context) {
	resp, _ := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
	var err error
	var transaction eacmodel.Transaction
	var transactionYear string
	var lastYearCounter uint64

	if transaction.CompanyID, err = resp.GetCompanyID("E1415152"); err != nil {
		return
	}

	if !resp.CheckRange(transaction.CompanyID) {
		return
	}

	transaction.Type = types.Enum(c.Param("type"))
	transactionYear = c.Param("year")

	if lastYearCounter, err = p.Service.LastYearCounterByType(transaction, transactionYear); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ViewTransaction)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, eacterm.LastYearCounter).
		JSON(lastYearCounter)
}

// List of transactions
func (p *TransactionAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacmodel.TransactionTable, eaccounting.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1446041"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ListTransaction)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, eacterm.Transactions).
		JSON(data)
}

// ManualTransfer transaction
func (p *TransactionAPI) ManualTransfer(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
	var transaction, createdTransaction eacmodel.Transaction
	var err error

	if transaction.CompanyID, transaction.NodeID, err = resp.GetCompanyNode("E1495400", eaccounting.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if transaction.CompanyID, err = resp.GetCompanyID("E1483718"); err != nil {
		return
	}

	if !resp.CheckRange(transaction.CompanyID) {
		return
	}

	if err = resp.Bind(&transaction, "E1444992", eaccounting.Domain, eacterm.Transaction); err != nil {
		return
	}

	transaction.Type = transactiontype.Manual
	transaction.CreatedBy = params.UserID

	if createdTransaction, err = p.Service.Transfer(transaction); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(eaccounting.ManualTransfer, transaction)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, eacterm.Transaction).
		JSON(createdTransaction)
}

// JournalEntry for creating a journal
func (p *TransactionAPI) JournalEntry(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
	var journal, createdJournal eacmodel.Transaction
	var err error

	if journal.CompanyID, journal.NodeID, err = resp.GetCompanyNode("E1425869", eaccounting.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if journal.CompanyID, err = resp.GetCompanyID("E1479868"); err != nil {
		return
	}

	if !resp.CheckRange(journal.CompanyID) {
		return
	}

	if err = resp.Bind(&journal, "E1497901", eaccounting.Domain, eacterm.Transaction); err != nil {
		return
	}

	if journal.Type == "" {
		journal.Type = transactiontype.JournalEntry
	}
	journal.CreatedBy = params.UserID

	transactionCh := eacmodel.TransactionCh{
		Transaction: journal,
		Type:        transactiontype.JournalEntry,
		Respond:     make(chan eacmodel.Transaction, 1),
	}

	p.Engine.TransactionCh <- transactionCh

	createdJournal = <-transactionCh.Respond
	if createdJournal.Err != nil {
		resp.Error(createdJournal.Err).JSON()
		return
	}

	resp.RecordCreate(eaccounting.EnterJournal, journal)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, eacterm.Journal).
		JSON(createdJournal)
}

// Update transaction
func (p *TransactionAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error

	var transaction, transactionBefore, transactionUpdated eacmodel.Transaction
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1452724", eacterm.Transaction); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&transaction, "E1451486", eaccounting.Domain, eacterm.Transaction); err != nil {
		return
	}

	if transactionBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	transaction.CreatedBy = transactionBefore.CreatedBy
	transaction.CreatedAt = transactionBefore.CreatedAt
	transaction.Hash = transactionBefore.Hash
	transaction.Type = transactionBefore.Type
	transaction.Invoice = transactionBefore.Invoice
	transaction.YearCounter = transactionBefore.YearCounter
	transaction.ID = fix.ID
	transaction.CompanyID = fix.CompanyID
	transaction.NodeID = fix.NodeID
	if transactionUpdated, err = p.Service.Save(transaction); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.UpdateTransaction, transactionBefore, transaction)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, eacterm.Transaction).
		JSON(transactionUpdated)
}

// JournalUpdate is used for updating the journal
func (p *TransactionAPI) JournalUpdate(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var journal, updatedJournal eacmodel.Transaction
	var err error
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1432163", eacterm.Transaction); err != nil {
		return
	}

	if err = resp.Bind(&journal, "E1431708", eaccounting.Domain, eacterm.Transaction); err != nil {
		return
	}

	journal.ID = fix.ID
	journal.CompanyID = fix.CompanyID
	journal.NodeID = fix.NodeID

	transactionCh := eacmodel.TransactionCh{
		Transaction: journal,
		Type:        transactiontype.JournalUpdate,
		Respond:     make(chan eacmodel.Transaction, 1),
	}

	p.Engine.TransactionCh <- transactionCh

	updatedJournal = <-transactionCh.Respond
	if updatedJournal.Err != nil {
		resp.Error(updatedJournal.Err).JSON()
		return
	}

	// if updatedJournal, journalBefore, err = p.Service.JournalUpdate(journal); err != nil {
	// 	resp.Error(err).JSON()
	// 	return
	// }

	resp.Record(eaccounting.UpdateJournal, updatedJournal.Before, journal)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, eacterm.Journal).
		JSON(updatedJournal)
}

// Delete transaction
func (p *TransactionAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var err error
	var transaction eacmodel.Transaction
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1413663", eacterm.Transaction); err != nil {
		return
	}

	if transaction, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.DeleteTransaction, transaction)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, eacterm.Transaction).
		JSON()
}

// Excel generate excel files eaced on search
func (p *TransactionAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1469354"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	transactions, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("transaction")
	ex.AddSheet("Transactions").
		AddSheet("Summary").
		Active("Transactions").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "F", 15.3).
		SetColWidth("G", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Transactions").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
		WriteData(transactions).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ExcelTransaction)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}

//JournalPrint  is for default print for journals, such as receipt voucher, paymenet voucher,
func (p *TransactionAPI) JournalPrint(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)

	var transaction eacmodel.Transaction
	var slot eacmodel.Slot
	var account basmodel.Account
	var currency eacmodel.Currency
	var total float64
	var totalText string

	//var currency eacmodel.Currency
	var fix types.FixedCol
	var err error
	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1423147", eacterm.Transaction); err != nil {
		return
	}

	//print type
	//printType := c.Param("type")

	// switch printType {
	// case "receipt-voucher", "payment-voucher", "receipt-entry", "payment-entry":

	// }

	//fethcing the transaction
	if transaction, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	//fmt.Println(fix)

	//defining the template
	var theme string
	//fetching the debit/credit slot
	switch transaction.Type {
	case transactiontype.ReceiptEntry, transactiontype.ReceiptVoucher:
		if slot, err = p.Service.FindDebit(transaction.Slots, transaction.Type); err != nil {
			resp.Error(err).JSON()
			return
		}
		theme = "receiptVoucher.tmpl"

	case transactiontype.PaymentEntry, transactiontype.PaymentVoucher:
		if slot, err = p.Service.FindCredit(transaction.Slots, transaction.Type); err != nil {
			resp.Error(err).JSON()
			return
		}
		theme = "paymentVoucher.tmpl"

	default:

		resp.Error(errors.New("this type of voucher cannot be used for printing")).JSON()
		return
	}

	//initiating fixed node for account service
	fixedNode := types.FixedNode{
		ID:        slot.AccountID,
		CompanyID: fix.CompanyID,
		NodeID:    fix.NodeID,
	}

	fmt.Println(fixedNode)

	//fetching account Name
	if account, err = p.Service.FindAccount(fixedNode); err != nil {
		return
	}
	//fetching currency of transaction
	fix.ID = transaction.CurrencyID
	if currency, err = p.Service.FetchCurrency(fix); err != nil {
		return
	}
	//fetching detailed slots with name and code
	var slots []eacmodel.DetailedSlots
	if slots, total, err = p.Service.FetchDetailedSlots(transaction.Slots, fixedNode, account.ID); err != nil {
		return
	}

	tempint := int(total)
	tempuint := uint(tempint)
	totalText = numtoword.EnConverter(tempuint)
	fmt.Println("detailed slots", slots)

	fmt.Println("the phone number is ", p.Engine.Setting[settingfields.CompanyPhone].Value)

	_ = transaction
	resp.Record(eaccounting.PrintJorunal)
	data := gin.H{
		"company": gin.H{
			"address": p.Engine.Setting[settingfields.CompanyAddress].Value,
			"logo":    p.Engine.Setting[settingfields.CompanyLogo].Value,
			"phone":   p.Engine.Setting[settingfields.CompanyPhone].Value,
			"email":   p.Engine.Setting[settingfields.CompanyEmail].Value,
		},
		"transaction": gin.H{
			"invoice":     transaction.Invoice,
			"accountName": *account.NameEn,
			"currency":    strings.ToUpper(currency.Name),
			"postdate":    transaction.PostDate.Format("02-01-2006"),
			// "createdAt":    transaction.CreatedAt.Format(consts.TimeLayout),
			"note":      transaction.Description,
			"printDate": time.Now().Format("02-01-2006"),
		},
		"total":     total,     //total amount of credit or debit
		"totalText": totalText, //total amount of credit or debit in Text

		"slots": slots,
	}

	themeContet, err := ioutil.ReadFile(filepath.Join("public", "vouchers-themes", theme))

	if err != nil {
		resp.Error(err).JSON()
		return
	}

	t := template.Must(template.New("example").Funcs(template.FuncMap{"counter": counter}).Parse(string(themeContet)))

	t.Execute(c.Writer, data)

	// c.HTML(http.StatusOK, theme, data)
}

func counter() func() int {
	i := -1
	return func() int {
		i++
		return i
	}
}
