package eacapi

import (
	"fmt"
	"net/http"
	"omono/domain/eaccounting"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/eaccounting/enum/transactiontype"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"

	"github.com/gin-gonic/gin"
)

// VoucherAPI for injecting transaction service
type VoucherAPI struct {
	Service service.EacVoucherServ
	Engine  *core.Engine
}

// ProvideVoucherAPI for voucher is used in wire
func ProvideVoucherAPI(c service.EacVoucherServ) VoucherAPI {
	return VoucherAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a transaction by it's id
// func (p *VoucherAPI) FindByID(c *gin.Context) {
// 	resp := response.New(p.Engine, c, eaccounting.Domain)
// 	var err error
// 	var transaction eacmodel.Transaction
// 	var fix types.FixedCol

// 	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1423147", eacterm.Transaction); err != nil {
// 		return
// 	}

// 	if !resp.CheckRange(fix.CompanyID) {
// 		return
// 	}

// 	if transaction, err = p.Service.FindByID(fix); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(eaccounting.ViewTransaction)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.VInfo, eacterm.Transaction).
// 		JSON(transaction)
// }

// List of transactions
// func (p *VoucherAPI) List(c *gin.Context) {
// 	resp, params := response.NewParam(p.Engine, c, eacmodel.TransactionTable, eaccounting.Domain)

// 	data := make(map[string]interface{})
// 	var err error

// 	if params.CompanyID, err = resp.GetCompanyID("E1446041"); err != nil {
// 		return
// 	}

// 	if !resp.CheckRange(params.CompanyID) {
// 		return
// 	}

// 	if data["list"], data["count"], err = p.Service.List(params); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(eaccounting.ListTransaction)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.ListOfV, eacterm.Transactions).
// 		JSON(data)
// }

// ManualTransfer transaction
// func (p *VoucherAPI) ManualTransfer(c *gin.Context) {
// 	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
// 	var transaction, createdTransaction eacmodel.Transaction
// 	var err error

// 	if transaction.CompanyID, transaction.NodeID, err = resp.GetCompanyNode("E1495400", eaccounting.Domain); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	if transaction.CompanyID, err = resp.GetCompanyID("E1483718"); err != nil {
// 		return
// 	}

// 	if !resp.CheckRange(transaction.CompanyID) {
// 		return
// 	}

// 	if err = resp.Bind(&transaction, "E1444992", eaccounting.Domain, eacterm.Transaction); err != nil {
// 		return
// 	}

// 	transaction.Type = transactiontype.Manual
// 	transaction.CreatedBy = params.UserID

// 	if createdTransaction, err = p.Service.Transfer(transaction); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.RecordCreate(eaccounting.ManualTransfer, transaction)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.VCreatedSuccessfully, eacterm.Transaction).
// 		JSON(createdTransaction)
// }

// JournalVoucher for creating a journal
func (p *VoucherAPI) JournalVoucher(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
	var journal, createdJournal eacmodel.Transaction
	var err error

	if journal.CompanyID, journal.NodeID, err = resp.GetCompanyNode("E1445208", eaccounting.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if journal.CompanyID, err = resp.GetCompanyID("E1472967"); err != nil {
		return
	}

	if !resp.CheckRange(journal.CompanyID) {
		return
	}

	if err = resp.Bind(&journal, "E1457076", eaccounting.Domain, eacterm.Transaction); err != nil {
		return
	}

	//if type is empty we set it as journal voucher by default
	if journal.Type == "" {
		journal.Type = transactiontype.JournalVoucher
	}
	journal.CreatedBy = params.UserID

	transactionCh := eacmodel.TransactionCh{
		Transaction: journal,
		Type:        "journal-voucher",
		Respond:     make(chan eacmodel.Transaction, 1),
	}

	p.Engine.TransactionCh <- transactionCh

	createdJournal = <-transactionCh.Respond
	if createdJournal.Err != nil {
		resp.Error(createdJournal.Err).JSON()
		return
	}
	// close(p.Engine.TransactionCh)
	// close(transactionCh.Respond)

	// glog.Debug(result)

	// if createdJournal, err = p.Service.JournalEntry(journal); err != nil {
	// 	resp.Error(err).JSON()
	// 	return
	// }

	// resp.RecordCreate(eaccounting.EnterJournal, journal)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, eacterm.Voucher).
		// JSON(createdJournal)
		JSON("ok")

}

//ApproveVoucher for apporving voucher
func (p *VoucherAPI) ApproveVoucher(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var voucher, updatedVoucher eacmodel.Transaction
	var err error
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1434268", eacterm.Transaction); err != nil {
		return
	}

	voucher.ID = fix.ID
	voucher.CompanyID = fix.CompanyID
	voucher.NodeID = fix.NodeID

	transactionCh := eacmodel.TransactionCh{
		Transaction: voucher,
		Type:        transactiontype.VoucherApprove,
		Respond:     make(chan eacmodel.Transaction, 1),
	}

	p.Engine.TransactionCh <- transactionCh

	updatedVoucher = <-transactionCh.Respond
	if updatedVoucher.Err != nil {
		resp.Error(updatedVoucher.Err).JSON()
		return
	}

	// if updatedVoucher, voucherBefore, err = p.Service.VoucherApprove(voucher); err != nil {
	// 	resp.Error(err).JSON()
	// 	return
	// }

	resp.Record(eaccounting.UpdateJournal, updatedVoucher.Before, voucher)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, eacterm.ApproveVoucher).
		JSON(updatedVoucher)

}

//FindByYearCounter will fetch the transaction based on year counter
func (p *VoucherAPI) FindByYearCounter(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
	var err error
	var transaction eacmodel.Transaction
	var slots []eacmodel.DetailedSlots
	var transactionYear string
	var transacstionType types.Enum
	var counter string

	if params.CompanyID, err = resp.GetCompanyID("E1469104"); err != nil {
		return
	}

	if !resp.CheckRange(transaction.CompanyID) {
		return
	}

	//holding the Params, which include: year of transacion, type of transaction, and the year counter of the requested transaction
	transacstionType = types.Enum(c.Param("type"))
	transactionYear = c.Param("year")
	counter = c.Param("counter")

	if transaction, slots, err = p.Service.FindByYearCounter(params.CompanyID, transacstionType, transactionYear, counter); err != nil {
		resp.Error(err).JSON()
		return
	}

	//setting the slots of the transaction to null since we have alreaedy fetched the deatil slots
	transaction.Slots = nil

	resp.Record(eaccounting.ViewTransaction)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, eacterm.Transaction).
		JSON(gin.H{"transaction": transaction, "slots": slots})

}

// VoucherUpdate is used for updating a voucher
func (p *VoucherAPI) VoucherUpdate(c *gin.Context) {
	resp := response.New(p.Engine, c, eaccounting.Domain)
	var voucher, updatedVoucher eacmodel.Transaction
	var err error
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("voucherID"), "E1466244", eacterm.Transaction); err != nil {
		return
	}

	if err = resp.Bind(&voucher, "E1445086", eaccounting.Domain, eacterm.Transaction); err != nil {
		return
	}

	fmt.Println(fix.ID.ToString(), fix.CompanyID, fix.NodeID)

	voucher.ID = fix.ID
	voucher.CompanyID = fix.CompanyID
	voucher.NodeID = fix.NodeID

	transactionCh := eacmodel.TransactionCh{
		Transaction: voucher,
		Type:        transactiontype.VoucherUpdate,
		Respond:     make(chan eacmodel.Transaction, 1),
	}

	p.Engine.TransactionCh <- transactionCh

	updatedVoucher = <-transactionCh.Respond
	if updatedVoucher.Err != nil {
		resp.Error(updatedVoucher.Err).JSON()
		return
	}

	resp.Record(eaccounting.UpdateJournal, updatedVoucher.Before, voucher)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, eacterm.Journal).
		JSON(updatedVoucher)
}

// Update transaction
// func (p *VoucherAPI) Update(c *gin.Context) {
// 	resp := response.New(p.Engine, c, eaccounting.Domain)
// 	var err error

// 	var transaction, transactionBefore, transactionUpdated eacmodel.Transaction
// 	var fix types.FixedCol

// 	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1452724", eacterm.Transaction); err != nil {
// 		return
// 	}

// 	if !resp.CheckRange(fix.CompanyID) {
// 		return
// 	}

// 	if err = resp.Bind(&transaction, "E1451486", eaccounting.Domain, eacterm.Transaction); err != nil {
// 		return
// 	}

// 	if transactionBefore, err = p.Service.FindByID(fix); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	transaction.CreatedBy = transactionBefore.CreatedBy
// 	transaction.Hash = transactionBefore.Hash
// 	transaction.Type = transactionBefore.Type
// 	transaction.Invoice = transactionBefore.Invoice
// 	transaction.YearCounter = transactionBefore.YearCounter
// 	transaction.ID = fix.ID
// 	transaction.CompanyID = fix.CompanyID
// 	transaction.NodeID = fix.NodeID
// 	if transactionUpdated, err = p.Service.Save(transaction); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(eaccounting.UpdateTransaction, transactionBefore, transaction)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.VUpdatedSuccessfully, eacterm.Transaction).
// 		JSON(transactionUpdated)
// }

// JournalUpdate is used for updating the journal
// func (p *VoucherAPI) JournalUpdate(c *gin.Context) {
// 	resp := response.New(p.Engine, c, eaccounting.Domain)
// 	var journal, journalBefore, journalUpdated eacmodel.Transaction
// 	var err error
// 	var fix types.FixedCol

// 	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1432163", eacterm.Transaction); err != nil {
// 		return
// 	}

// 	if err = resp.Bind(&journal, "E1431708", eaccounting.Domain, eacterm.Transaction); err != nil {
// 		return
// 	}

// 	journal.ID = fix.ID
// 	journal.CompanyID = fix.CompanyID
// 	journal.NodeID = fix.NodeID

// 	/* move to service
// 	if journalBefore, err = p.Service.FindByID(fix); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	journal.CreatedBy = journalBefore.CreatedBy
// 	journal.Hash = journalBefore.Hash
// 	journal.Type = journalBefore.Type
// 	journal.Invoice = journalBefore.Invoice
// 	journal.YearCounter = journalBefore.YearCounter
// 	*/

// 	if journalUpdated, journalBefore, err = p.Service.JournalUpdate(journal); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(eaccounting.UpdateJournal, journalBefore, journal)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.VUpdatedSuccessfully, eacterm.Journal).
// 		JSON(journalUpdated)
// }

// // Delete transaction
// func (p *VoucherAPI) Delete(c *gin.Context) {
// 	resp := response.New(p.Engine, c, eaccounting.Domain)
// 	var err error
// 	var transaction eacmodel.Transaction
// 	var fix types.FixedCol

// 	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1413663", eacterm.Transaction); err != nil {
// 		return
// 	}

// 	if transaction, err = p.Service.Delete(fix); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(eaccounting.DeleteTransaction, transaction)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.VDeletedSuccessfully, eacterm.Transaction).
// 		JSON()
// }

// Excel generate excel files eaced on search
// func (p *VoucherAPI) Excel(c *gin.Context) {
// 	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
// 	var err error

// 	if params.CompanyID, err = resp.GetCompanyID("E1469354"); err != nil {
// 		return
// 	}

// 	if !resp.CheckRange(params.CompanyID) {
// 		return
// 	}

// 	transactions, err := p.Service.Excel(params)
// 	if err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	ex := excel.New("transaction")
// 	ex.AddSheet("Transactions").
// 		AddSheet("Summary").
// 		Active("Transactions").
// 		SetPageLayout("landscape", "A4").
// 		SetPageMargins(0.2).
// 		SetHeaderFooter().
// 		SetColWidth("B", "F", 15.3).
// 		SetColWidth("G", "G", 40).
// 		Active("Summary").
// 		SetColWidth("A", "D", 20).
// 		Active("Transactions").
// 		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
// 		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
// 		WriteData(transactions).
// 		AddTable()

// 	buffer, downloadName, err := ex.Generate()
// 	if err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(eaccounting.ExcelTransaction)

// 	c.Header("Content-Description", "File Transfer")
// 	c.Header("Content-Disposition", "attachment; filename="+downloadName)
// 	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

// }

// //JournalPrint  is for default print for journals, such as receipt voucher, paymenet voucher,
// func (p *VoucherAPI) JournalPrint(c *gin.Context) {
// 	resp := response.New(p.Engine, c, eaccounting.Domain)

// 	var transaction eacmodel.Transaction
// 	var slot eacmodel.Slot
// 	var account basmodel.Account
// 	var currency eacmodel.Currency
// 	var total float64
// 	var totalText string

// 	//var currency eacmodel.Currency
// 	var fix types.FixedCol
// 	var err error
// 	if fix, err = resp.GetFixedCol(c.Param("transactionID"), "E1423147", eacterm.Transaction); err != nil {
// 		return
// 	}

// 	//print type
// 	//printType := c.Param("type")

// 	// switch printType {
// 	// case "receipt-voucher", "payment-voucher", "receipt-entry", "payment-entry":

// 	// }

// 	//fethcing the transaction
// 	if transaction, err = p.Service.FindByID(fix); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	//fmt.Println(fix)

// 	//defining the template
// 	var theme string
// 	//fetching the debit/credit slot
// 	switch transaction.Type {
// 	case transactiontype.ReceiptEntry, transactiontype.ReceiptVoucher:
// 		if slot, err = p.Service.FindDebit(transaction.Slots, transaction.Type); err != nil {
// 			resp.Error(err).JSON()
// 			return
// 		}
// 		theme = "receiptVoucher.tmpl"

// 	case transactiontype.PaymentEntry, transactiontype.PaymentVoucher:
// 		if slot, err = p.Service.FindCredit(transaction.Slots, transaction.Type); err != nil {
// 			resp.Error(err).JSON()
// 			return
// 		}
// 		theme = "paymentVoucher.tmpl"

// 	default:

// 		resp.Error(errors.New("this type of voucher cannot be used for printing")).JSON()
// 		return
// 	}

// 	//initiating fixed node for account service
// 	fixedNode := types.FixedNode{
// 		ID:        slot.AccountID,
// 		CompanyID: fix.CompanyID,
// 		NodeID:    fix.NodeID,
// 	}

// 	fmt.Println(fixedNode)

// 	//fetching account Name
// 	if account, err = p.Service.FindAccount(fixedNode); err != nil {
// 		return
// 	}
// 	//fetching currency of transaction
// 	fix.ID = transaction.CurrencyID
// 	if currency, err = p.Service.FetchCurrency(fix); err != nil {
// 		return
// 	}
// 	//fetching detailed slots with name and code
// 	var slots []eacmodel.DetailedSlots
// 	if slots, total, err = p.Service.FetchDetailedSlots(transaction.Slots, fixedNode, account.ID); err != nil {
// 		return
// 	}

// 	tempint := int(total)
// 	tempuint := uint(tempint)
// 	totalText = numtoword.EnConverter(tempuint)
// 	fmt.Println("detailed slots", slots)

// 	fmt.Println("the phone number is ", p.Engine.Setting[settingfields.CompanyPhone].Value)

// 	_ = transaction
// 	resp.Record(eaccounting.PrintJorunal)
// 	data := gin.H{
// 		"company": gin.H{
// 			"address": p.Engine.Setting[settingfields.CompanyAddress].Value,
// 			"logo":    p.Engine.Setting[settingfields.CompanyLogo].Value,
// 			"phone":   p.Engine.Setting[settingfields.CompanyPhone].Value,
// 			"email":   p.Engine.Setting[settingfields.CompanyEmail].Value,
// 		},
// 		"transaction": gin.H{
// 			"invoice":     transaction.Invoice,
// 			"accountName": *account.NameEn,
// 			"currency":    strings.ToUpper(currency.Name),
// 			"postdate":    transaction.PostDate.Format("02-01-2006"),
// 			// "createdAt":    transaction.CreatedAt.Format(consts.TimeLayout),
// 			"note":      transaction.Description,
// 			"printDate": time.Now().Format("02-01-2006"),
// 		},
// 		"total":     total,     //total amount of credit or debit
// 		"totalText": totalText, //total amount of credit or debit in Text

// 		"slots": slots,
// 	}

// 	themeContet, err := ioutil.ReadFile(filepath.Join("public", "vouchers-themes", theme))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	t := template.Must(template.New("example").Funcs(template.FuncMap{"counter": counter}).Parse(string(themeContet)))

// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	t.Execute(c.Writer, data)

// 	//c.HTML(http.StatusOK, theme, data)
// }

// func counter() func() int {
// 	i := -1
// 	return func() int {
// 		i++
// 		return i
// 	}
// }
