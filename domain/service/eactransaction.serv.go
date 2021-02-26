package service

import (
	"errors"
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/eaccounting/enum/transactionstatus"
	"omono/domain/eaccounting/enum/transactiontype"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
	"time"

	"gorm.io/gorm"
)

// EacTransactionServ for injecting auth eacrepo
type EacTransactionServ struct {
	Repo     eacrepo.TransactionRepo
	Engine   *core.Engine
	SlotServ EacSlotServ
}

// ProvideEacTransactionService for transaction is used in wire
func ProvideEacTransactionService(p eacrepo.TransactionRepo, slotServ EacSlotServ) EacTransactionServ {
	return EacTransactionServ{
		Repo:     p,
		Engine:   p.Engine,
		SlotServ: slotServ,
	}
}

// FindByID for getting transaction by it's id
func (p *EacTransactionServ) FindByID(fix types.FixedCol) (transaction eacmodel.Transaction, err error) {
	if transaction, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1463467", "can't fetch the transaction", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	if transaction.Status == transactionstatus.Unapproved {
		voucherServ := fetchVoucherService(p.Engine)
		if transaction.Slots, err = voucherServ.TempSlotServ.TemporarySlots(transaction.ID,
			transaction.CompanyID, transaction.NodeID); err != nil {
			err = corerr.Tick(err, "E1477159", "can't fetch the temproray slots",
				fix.CompanyID, fix.NodeID, transaction.ID)
			return
		}
	} else {
		if transaction.Slots, err = p.SlotServ.TransactionSlot(transaction.ID,
			transaction.CompanyID, transaction.NodeID); err != nil {
			err = corerr.Tick(err, "E1464213", "can't fetch the transactions slots",
				fix.CompanyID, fix.NodeID, transaction.ID)
			return
		}
	}

	return
}

// LastYearCounterByType for getting last transaction year counter based on the type
func (p *EacTransactionServ) LastYearCounterByType(transaction eacmodel.Transaction, year string) (yearCounter uint64, err error) {

	var dateStart time.Time
	if err = transaction.Validate(coract.Fetch); err != nil {
		err = corerr.TickValidate(err, "E1488349", "validation failed in fetching the transaction", transaction)
		return
	}
	if dateStart, err = time.Parse("2006", year); err != nil {
		err = corerr.Tick(err, "E1494575", "can't convert the year to date", year)
		return
	}
	tempDate := dateStart.AddDate(0, 11, 0)
	y, m, _ := tempDate.Date()
	dateEnd := time.Date(y, m+1, 0, 23, 59, 59, 0, time.UTC)
	//getting the similar type of voucher
	similarType := transactiontype.SimilarType(transaction.Type)
	if yearCounter, err = p.Repo.LastYearCounterByType(transaction, similarType, dateStart, dateEnd); err != nil {
		err = corerr.Tick(err, "E1412057", "can't fetch the transaction year counter", transaction.CompanyID, transaction.Type, year)
		return
	}

	//now we add +1 to year counter
	yearCounter++

	return
}

// List of transactions, it support pagination and search and return back count
func (p *EacTransactionServ) List(params param.Param) (transactions []eacmodel.Transaction,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" eac_transactions.company_id = '%v' ", params.CompanyID)
	}

	if transactions, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in transactions list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in transactions count")
	}

	return
}

// Transfer activate create for special transfering
func (p *EacTransactionServ) Transfer(transaction eacmodel.Transaction) (createdTransaction eacmodel.Transaction, err error) {
	slots := []eacmodel.Slot{
		{
			AccountID:  transaction.Pioneer,
			Credit:     transaction.Amount,
			CurrencyID: transaction.CurrencyID,
			PostDate:   transaction.PostDate,
		},
		{
			AccountID:  transaction.Follower,
			Debit:      transaction.Amount,
			CurrencyID: transaction.CurrencyID,
			PostDate:   transaction.PostDate,
		},
	}

	slots[0].CompanyID = transaction.CompanyID
	slots[1].CompanyID = transaction.CompanyID
	slots[0].NodeID = transaction.NodeID
	slots[1].NodeID = transaction.NodeID

	createdTransaction, err = p.Create(transaction, slots)

	return
}

// JournalEntryWatcher is used for watching transactions
func (p *EacTransactionServ) JournalEntryWatcher() {
	for t := range p.Engine.TransactionCh {

		var err error
		var before eacmodel.Transaction
		switch t.Type {
		case transactiontype.JournalEntry:
			t.Transaction, err = p.JournalEntry(t.Transaction)
			t.Transaction.Err = err

		case transactiontype.JournalUpdate:
			t.Transaction, before, err = p.JournalUpdate(t.Transaction)
			t.Transaction.Before = &before
			t.Transaction.Err = err

		case transactiontype.JournalVoucher:
			voucherServ := fetchVoucherService(p.Engine)
			//sending the request
			t.Transaction, err = voucherServ.JournalVoucher(t.Transaction)
			t.Transaction.Err = err

		case transactiontype.VoucherApprove:
			voucherServ := fetchVoucherService(p.Engine)

			//sending the request
			t.Transaction, _, err = voucherServ.ApproveVoucher(t.Transaction)
			t.Transaction.Err = err
		case transactiontype.VoucherUpdate:
			voucherServ := fetchVoucherService(p.Engine)
			//sending the request
			t.Transaction, before, err = voucherServ.VoucherUpdate(t.Transaction)
			t.Transaction.Before = &before
			t.Transaction.Err = err
		default:
			{
				t.Transaction.Err = limberr.New("transaction type is not valid", "E1459638").
					Custom(corerr.ForbiddenErr).
					Message(eacterm.TransactionTypeIsNotValid).Build()
			}
		}

		t.Respond <- t.Transaction
		// close(t.Respond)
	}
}

// JournalEntry insert a journal to the system
func (p *EacTransactionServ) JournalEntry(journal eacmodel.Transaction) (createdJournal eacmodel.Transaction, err error) {
	var zero float64

	//we set the transaction to approved since its journal entry
	journal.Status = transactionstatus.Approved

	//we make validation for the transaction
	if err = journal.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1475662", "validation failed in creating the journal", journal)
		return
	}

	//we set companyId nodeId postdate and currencyId for each slot of the
	//we also compare to see if the difference between credit and debit is zero
	for i := range journal.Slots {
		journal.Slots[i].CompanyID = journal.CompanyID
		journal.Slots[i].NodeID = journal.NodeID
		journal.Slots[i].PostDate = journal.PostDate
		journal.Slots[i].CurrencyID = journal.CurrencyID
		zero += journal.Slots[i].Debit - journal.Slots[i].Credit
	}

	//if not zero return error
	if zero != 0 {
		err = limberr.New("difference is not zero", "E1453653").
			Message(eacterm.DifferenceIsNotZeroV, zero).
			Custom(corerr.ForbiddenErr).Build()
		return
	}

	//we start db instance
	db := p.Engine.DB.Begin()

	//generate the invoice for the journal, and get the yearcumulative and journal invoice
	if journal.YearCounter, journal.YearCumulative, journal.Invoice, err = p.InvoiceGenerator(db, journal.CompanyID,
		journal.Type, journal.PostDate); err != nil {
		err = corerr.Tick(err, "E1415692", "error in generating invoice for journal", journal)
		db.Rollback()
		return
	}

	//recover go routine panicking routine
	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"eac_transactions table"), "rollback recover tx-create transaction")
			db.Rollback()
		}
	}()

	///create the journal, inside p.TXcreat() the slots will also be inserted
	if createdJournal, err = p.TxCreate(db, journal, journal.Slots); err != nil {
		err = corerr.Tick(err, "E1412707", "transaction not created", journal)
		db.Rollback()
		return
	}

	//commit the transaction
	db.Commit()

	return
}

// JournalUpdate is used for editing the journal
func (p *EacTransactionServ) JournalUpdate(journal eacmodel.Transaction) (journalUpdated, journalBefore eacmodel.Transaction, err error) {
	fix := types.FixedCol{
		CompanyID: journal.CompanyID,
		NodeID:    journal.NodeID,
		ID:        journal.ID,
	}
	//update validation
	if err = journal.Validate(coract.Update); err != nil {
		err = corerr.TickValidate(err, "E1420523", "validation failed in fetching the transaction", journal)
		return
	}

	// TODO: validation for journal-update

	if journalBefore, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1423702", "transaction not found for update", fix)
		return
	}

	journal.CreatedBy = journalBefore.CreatedBy
	journal.Hash = journalBefore.Hash
	journal.Type = journalBefore.Type
	journal.Invoice = journalBefore.Invoice
	journal.YearCounter = journalBefore.YearCounter
	journal.YearCumulative = journalBefore.YearCumulative

	journal.CreatedAt = journalBefore.CreatedAt
	journal.Status = journalBefore.Status

	newSlots, deletedSlots, updatedSlots := p.SlotServ.discreteSlots(journalBefore.Slots, journal.Slots)

	// if err = p.SlotServ.duplicateAccounts(journal.Slots, newSlots); err != nil {
	// 	err = corerr.Tick(err, "E1420523", "duplicate accounts exists which is invalid", journal.Slots)
	// 	return

	// }

	db := p.Engine.DB.Begin()

	fmt.Printf("the new slots \n %v\n", newSlots)
	fmt.Printf("the updated slots slots \n %v\n", updatedSlots)
	fmt.Printf("the deleted slots \n %v\n", deletedSlots)
	if journalUpdated, err = p.Repo.Save(journal); err != nil {
		err = corerr.Tick(err, "E1416320", "transaction not saved")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in %v",
				"eac_transactions table"), "rollback recover JournalUpdate")
			db.Rollback()
		}
	}()

	for _, v := range newSlots {
		v.TransactionID = journal.ID
		v.PostDate = journal.PostDate
		v.CompanyID = journal.CompanyID
		v.NodeID = journal.NodeID
		v.CurrencyID = journal.CurrencyID

		if _, err = p.SlotServ.TxCreate(db, v); err != nil {
			err = corerr.Tick(err, "E1496205", "slot not created in JournalUpdate", v)
			db.Rollback()
			return
		}
	}

	for _, v := range updatedSlots {
		v.TransactionID = journal.ID
		v.PostDate = journal.PostDate
		v.CompanyID = journal.CompanyID
		v.NodeID = journal.NodeID
		v.CurrencyID = journal.CurrencyID
		if _, err = p.SlotServ.TxSave(db, v); err != nil {
			err = corerr.Tick(err, "E1438909", "slot not saved in JournalUpdate", v)
			db.Rollback()
			return
		}
	}

	for _, v := range deletedSlots {
		v.TransactionID = journal.ID
		adjust := v.Credit - v.Debit
		if err = p.SlotServ.Repo.TxRegulateBalancesSave(db, v, adjust); err != nil {
			err = corerr.Tick(err, "E1418024", "error in regulating the balance", v)
			db.Rollback()
			return
		}

		if err = p.SlotServ.Repo.TxDelete(db, v); err != nil {
			err = corerr.Tick(err, "E1433649", "error in deleting the slot", v)
			db.Rollback()
			return
		}
	}

	db.Commit()

	return
}

// InvoiceGenerator is used for create a unique number for each transaction based on year and
// companyID
func (p *EacTransactionServ) InvoiceGenerator(db *gorm.DB, companyID uint64, tType types.Enum,
	postDate time.Time) (yearCounter uint64, yearCumulative uint64, invoice string, err error) {

	// secondtType
	lastDay := time.Date(postDate.Year()+1, 1, 1, 0, 0, 0, 0, time.Now().Location())
	var transaction eacmodel.Transaction

	//fethcing the similar type
	secondType := transactiontype.SimilarType(tType)

	if transaction, err = p.Repo.LastYearCounter(db, companyID, tType, secondType, lastDay); err != nil {
		err = corerr.Tick(err, "E1469786", "error in fetching the last transaction", companyID, tType, postDate)
		return
	}

	yearCounter = transaction.YearCounter + 1

	if transaction, err = p.Repo.LastYearCumulative(db, companyID, lastDay); err != nil {
		err = corerr.Tick(err, "E1457746", "error in fetching the last transaction for year_cumulative", companyID, tType, postDate)
		return
	}
	yearCumulative = transaction.YearCumulative + 1

	// invoice = fmt.Sprintf("%v-%v", postDate.Year(), yearCounter)
	invoice = fmt.Sprintf("%v", yearCounter)

	return
}

// Create a transaction, TODO: refactor based on TX design
func (p *EacTransactionServ) Create(transaction eacmodel.Transaction,
	slots []eacmodel.Slot) (createdTransaction eacmodel.Transaction, err error) {

	if err = transaction.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1433872", "validation failed in creating the transaction", transaction)
		return
	}

	clonedEngine := p.Engine.Clone()
	clonedEngine.DB = clonedEngine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"eac_transactions table"), "rollback recover create transaction")
			clonedEngine.DB.Rollback()
		}
	}()

	transactionRepo := eacrepo.ProvideTransactionRepo(clonedEngine)
	slotServ := ProvideEacSlotService(eacrepo.ProvideSlotRepo(clonedEngine),
		p.SlotServ.CurrencyServ, p.SlotServ.AccountServ)

	now := time.Now()
	transaction.Hash = now.Format(consts.HashTimeLayout)

	if transaction.YearCounter, transaction.YearCumulative, transaction.Invoice, err = p.InvoiceGenerator(clonedEngine.DB, transaction.CompanyID,
		transaction.Type, transaction.PostDate); err != nil {
		err = corerr.Tick(err, "E1420835", "error in generating invoice for transaction", transaction)
		return
	}

	if createdTransaction, err = transactionRepo.Create(transaction); err != nil {
		err = corerr.Tick(err, "E1479603", "transaction not created", transaction)

		clonedEngine.DB.Rollback()
		return
	}

	for _, v := range slots {
		v.TransactionID = createdTransaction.ID
		if _, err = slotServ.Create(v); err != nil {
			err = corerr.Tick(err, "E1420630", "slot not saved in transaction creation", v)
			clonedEngine.DB.Rollback()
			return
		}
	}

	clonedEngine.DB.Commit()

	return
}

// TxCreate is used for activating rollback in case of error
func (p *EacTransactionServ) TxCreate(db *gorm.DB, transaction eacmodel.Transaction,
	slots []eacmodel.Slot) (createdTransaction eacmodel.Transaction, err error) {

	if err = transaction.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1475662", "validation failed in tx-creating the transaction", transaction)
		return
	}

	now := time.Now()
	transaction.Hash = now.Format(consts.HashTimeLayout)

	if createdTransaction, err = p.Repo.TxCreate(db, transaction); err != nil {
		err = corerr.Tick(err, "E1446483", "transaction not created", transaction)
		return
	}

	for _, v := range slots {
		v.TransactionID = createdTransaction.ID
		if _, err = p.SlotServ.TxCreate(db, v); err != nil {
			err = corerr.Tick(err, "E1429080", "slot not saved in TxCreate", v)
			return
		}
	}

	return
}

// EditTransfer activate create for special transfering
func (p *EacTransactionServ) EditTransfer(tr eacmodel.Transaction) (updatedTr eacmodel.Transaction, err error) {
	slots := []eacmodel.Slot{
		{
			AccountID:  tr.Pioneer,
			Credit:     tr.Amount,
			CurrencyID: tr.CurrencyID,
			PostDate:   tr.PostDate,
		},
		{
			AccountID:  tr.Follower,
			Debit:      tr.Amount,
			CurrencyID: tr.CurrencyID,
			PostDate:   tr.PostDate,
		},
	}

	fix := types.FixedCol{
		CompanyID: tr.CompanyID,
		NodeID:    tr.NodeID,
		ID:        tr.ID,
	}

	var oldTr eacmodel.Transaction
	if oldTr, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1484155", "edit transfer can't find the transaction", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	slots[0].CompanyID = tr.CompanyID
	slots[0].NodeID = tr.NodeID
	slots[0].ID = oldTr.Slots[0].ID

	slots[1].CompanyID = tr.CompanyID
	slots[1].NodeID = tr.NodeID
	slots[1].ID = oldTr.Slots[1].ID

	updatedTr, err = p.Update(tr, slots)

	return
}

// Update is used when a transaction has been changed
func (p *EacTransactionServ) Update(tr eacmodel.Transaction,
	slots []eacmodel.Slot) (updatedTr eacmodel.Transaction, err error) {

	if err = tr.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1469927", "validation failed in updating the transaction", tr)
		return
	}

	clonedEngine := p.Engine.Clone()
	clonedEngine.DB = clonedEngine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"eac_transactions table"), "rollback recover update transaction")
			clonedEngine.DB.Rollback()
		}
	}()

	transactionRepo := eacrepo.ProvideTransactionRepo(clonedEngine)
	slotServ := ProvideEacSlotService(eacrepo.ProvideSlotRepo(clonedEngine),
		p.SlotServ.CurrencyServ, p.SlotServ.AccountServ)

	now := time.Now()
	tr.Hash = now.Format(consts.HashTimeLayout)

	if updatedTr, err = transactionRepo.Save(tr); err != nil {
		err = corerr.Tick(err, "E1479603", "transaction not updated", tr)

		clonedEngine.DB.Rollback()
		return
	}

	for _, v := range slots {
		v.TransactionID = updatedTr.ID
		if _, err = slotServ.Save(v); err != nil {
			err = corerr.Tick(err, "E1420630", "slot not saved in updating the transaction", v)
			clonedEngine.DB.Rollback()
			return
		}
	}

	clonedEngine.DB.Commit()

	return

}

// Save a transaction, if it is exist update it, if not create it
func (p *EacTransactionServ) Save(transaction eacmodel.Transaction) (savedTransaction eacmodel.Transaction, err error) {
	if err = transaction.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1414478", corerr.ValidationFailed, transaction)
		return
	}

	clonedEngine := p.Engine.Clone()
	clonedEngine.DB = clonedEngine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"eac_transactions table"), "rollback recover update transaction")
			clonedEngine.DB.Rollback()
		}
	}()

	transactionRepo := eacrepo.ProvideTransactionRepo(clonedEngine)
	slotServ := ProvideEacSlotService(eacrepo.ProvideSlotRepo(clonedEngine),
		p.SlotServ.CurrencyServ, p.SlotServ.AccountServ)

	if savedTransaction, err = transactionRepo.Save(transaction); err != nil {
		err = corerr.Tick(err, "E1482909", "transaction not saved", transaction)

		clonedEngine.DB.Rollback()
		return
	}

	for _, v := range transaction.Slots {
		v.PostDate = transaction.PostDate
		v.TransactionID = transaction.ID
		if _, err = slotServ.Save(v); err != nil {
			err = corerr.Tick(err, "E1485649", "slot not saved in transaction edit", v)
			clonedEngine.DB.Rollback()
			return
		}
	}

	clonedEngine.DB.Commit()

	return
}

// Delete transaction, it is soft delete
func (p *EacTransactionServ) Delete(fix types.FixedCol) (transaction eacmodel.Transaction, err error) {
	if transaction, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1431984", "transaction not found for deleting")
		return
	}

	db := p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(errors.New("panic happened in rollback mode for deleteing a transaction"),
				"rollback recover delete a transaction")
			db.Rollback()
		}
	}()

	for _, v := range transaction.Slots {
		adjust := v.Credit - v.Debit
		if err = p.SlotServ.Repo.TxRegulateBalancesSave(db, v, adjust); err != nil {
			err = corerr.Tick(err, "E1414828", "error in regulating the balance", v)
			db.Rollback()
			return
		}

		if err = p.SlotServ.Repo.TxDelete(db, v); err != nil {
			err = corerr.Tick(err, "E1480067", "error in deleting the slot", v)
			db.Rollback()
			return
		}

		if err = p.SlotServ.TxUpdateBalance(db, v); err != nil {
			err = corerr.Tick(err, "E1490428", "update account's balance after delete slot", v)
			db.Rollback()
			return
		}

	}

	if err = p.Repo.TxDelete(db, transaction); err != nil {
		err = corerr.Tick(err, "E1440895", "error in deleting the transaction", transaction.ID)
		db.Rollback()
		return
	}

	db.Commit()

	return
}

// Excel is used for export excel file
func (p *EacTransactionServ) Excel(params param.Param) (transactions []eacmodel.Transaction, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", eacmodel.TransactionTable)

	if transactions, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1426905", "cant generate the excel list for transactions")
		return
	}

	return
}

//FindDebit will return the debit account among slots ... used for receipt voucher & entry currently
func (p *EacTransactionServ) FindDebit(slots []eacmodel.Slot, transactionType types.Enum) (debitAccount eacmodel.Slot, err error) {

	//inital checking whether the type of transaction is correct

	if transactionType != transactiontype.ReceiptEntry && transactionType != transactiontype.ReceiptVoucher {
		err = corerr.Tick(err, "E1426906", "Type of transaciton is not Receipt voucher or Receipt Entry")
		return
	}

	//going through the slots to find the debit account
	for _, v := range slots {
		if v.Debit > 0 {
			debitAccount = v
			break
		}
	}

	//check whether debit account has been selected
	if debitAccount == (eacmodel.Slot{}) {
		err = corerr.Tick(err, "E1426907", "debit slot was not found ")

	}

	return
}

//FindCredit will return the debit account among slots ... used for Payment voucher & entry currently
func (p *EacTransactionServ) FindCredit(slots []eacmodel.Slot, transactionType types.Enum) (creditAccount eacmodel.Slot, err error) {

	//inital checking whether the type of transaction is correct
	if transactionType != transactiontype.PaymentEntry && transactionType != transactiontype.PaymentVoucher {
		err = corerr.Tick(err, "E1426918", "Type of transaciton is not Payment voucher or Payment Entry")
		return
	}

	//going through the slots to find the debit account
	for _, v := range slots {
		if v.Credit > 0 {
			creditAccount = v
			break
		}
	}

	//check whether credit account has been selected
	if creditAccount == (eacmodel.Slot{}) {
		err = corerr.Tick(err, "E1426919", "credit slot was not found ")

	}

	return
}

//FindAccount will find the account that belongs to a transaction.. currently used for journal print to identify accounts
func (p *EacTransactionServ) FindAccount(fix types.FixedNode) (account basmodel.Account, err error) {

	//providing the services for account service
	basAccountRepo := basrepo.ProvideAccountRepo(p.Engine)
	basPhoneRepo := basrepo.ProvidePhoneRepo(p.Engine)
	basPhoneServ := ProvideBasPhoneService(basPhoneRepo)

	//providing account service
	basAccount := ProvideBasAccountService(basAccountRepo, basPhoneServ)

	if account, err = basAccount.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1426908", "account was not found")

	}

	return

}

//FetchCurrency will find the account that belongs to a transaction.. currently used for journal print to identify accounts
func (p *EacTransactionServ) FetchCurrency(fix types.FixedCol) (currency eacmodel.Currency, err error) {

	//providing the repo for currency service
	eacCurrencyRepo := eacrepo.ProvideCurrencyRepo(p.Engine)

	//providing account service
	currencyServ := ProvideEacCurrencyService(eacCurrencyRepo)

	if currency, err = currencyServ.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1426909", "account was not found")

	}

	return

}

//FetchDetailedSlots will created a detailed version for the slots
func (p *EacTransactionServ) FetchDetailedSlots(slots []eacmodel.Slot, fix types.FixedNode, exceptAccount types.RowID) (detailedSlots []eacmodel.DetailedSlots, totalAmount float64, err error) {

	//temporary slot
	var tempSlot eacmodel.DetailedSlots
	//temporary accountName
	var tempAcc basmodel.Account
	for _, v := range slots {

		//SKIPPING THE ACCOUNT THAT WILL BE NOT NEEDED
		if v.AccountID == exceptAccount {
			continue
		}
		tempSlot.ID = v.ID
		tempSlot.AccountID = v.AccountID
		tempSlot.CurrencyID = v.CurrencyID
		tempSlot.TransactionID = v.TransactionID
		tempSlot.Debit = v.Debit
		tempSlot.Credit = v.Credit
		tempSlot.Balance = v.Balance
		tempSlot.Description = v.Description
		tempSlot.PostDate = v.PostDate

		//fethcing account
		fix.ID = v.AccountID
		if tempAcc, err = p.FindAccount(fix); err != nil {
			return
		}
		tempSlot.AccountName = *tempAcc.NameEn
		tempSlot.AccountCode = tempAcc.Code

		//addding total based on the type receipt|Payment
		if v.Credit == 0 {
			totalAmount += tempSlot.Debit

		} else if v.Debit == 0 {
			totalAmount += tempSlot.Credit

		}

		detailedSlots = append(detailedSlots, tempSlot)
	}

	return
}

func fetchVoucherService(engine *core.Engine) (voucherServ EacVoucherServ) {
	transactionRepo := eacrepo.ProvideVoucherRepo(engine)
	tempSlotRepo := eacrepo.ProvideTempSlotRepo(engine)
	currencyServ := ProvideEacCurrencyService(eacrepo.ProvideCurrencyRepo(engine))
	phoneServ := ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountServ := ProvideBasAccountService(basrepo.ProvideAccountRepo(engine), phoneServ)

	tempSlotServ := ProvideEacTempSlotService(tempSlotRepo, currencyServ, accountServ)
	voucherServ = ProvideEacVoucherService(transactionRepo, tempSlotServ)

	return voucherServ
}
