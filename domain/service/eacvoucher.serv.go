package service

import (
	"errors"
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/base/basrepo"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/eaccounting/enum/transactionstatus"
	"omono/domain/eaccounting/enum/transactiontype"
	"omono/domain/notification"
	"omono/domain/notification/notmodel"
	"omono/domain/notification/notrepo"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/types"
	"omono/pkg/glog"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// EacVoucherServ for injecting auth eacrepo
type EacVoucherServ struct {
	Repo         eacrepo.VoucherRepo
	Engine       *core.Engine
	TempSlotServ EacTempSlotServ
}

// ProvideEacVoucherService for transaction is used in wire
func ProvideEacVoucherService(p eacrepo.VoucherRepo, tempSlotServ EacTempSlotServ) EacVoucherServ {
	return EacVoucherServ{
		Repo:         p,
		Engine:       p.Engine,
		TempSlotServ: tempSlotServ,
	}
}

// JournalVoucher insert a journal voucher to the system
func (p *EacVoucherServ) JournalVoucher(voucher eacmodel.Transaction) (createdVoucher eacmodel.Transaction, err error) {
	var zero float64
	//differrence validation between credit and debit
	for i := range voucher.Slots {
		voucher.Slots[i].CompanyID = voucher.CompanyID
		voucher.Slots[i].NodeID = voucher.NodeID
		voucher.Slots[i].PostDate = voucher.PostDate
		voucher.Slots[i].CurrencyID = voucher.CurrencyID
		zero += voucher.Slots[i].Debit - voucher.Slots[i].Credit
	}

	if zero != 0 {
		err = limberr.New("difference is not zero", "E1462728").
			Message(eacterm.DifferenceIsNotZeroV, zero).
			Custom(corerr.ForbiddenErr).Build()
		return
	}
	db := p.Engine.DB.Begin()

	//p

	//fetchin transactionServ
	transactionServ, accountServ := fetchTransactionServ(p.Engine)
	if voucher.YearCounter, voucher.YearCumulative, voucher.Invoice, err = transactionServ.InvoiceGenerator(db, voucher.CompanyID,
		voucher.Type, voucher.PostDate); err != nil {
		err = corerr.Tick(err, "E1459079", "error in generating invoice for voucher", voucher)
		db.Rollback()
		return
	}

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"eac_transactions table"), "rollback recover tx-create transaction")
			db.Rollback()
		}
	}()

	if createdVoucher, err = p.TxCreateVoucher(db, voucher, voucher.Slots); err != nil {
		err = corerr.Tick(err, "E1422568", "transaction not created", voucher)
		db.Rollback()
	}

	//now lets create a notification if it exists

	if voucher.UserID != 0 {
		fix := types.FixedNode{
			ID:        voucher.UserID,
			CompanyID: voucher.CompanyID,
			NodeID:    voucher.NodeID,
		}
		if _, err := accountServ.FindByID(fix); err != nil {
			err = corerr.Tick(err, "E1414189", "the send To account does not exist, hence notification cannot be generated", voucher)
			db.Rollback()
		}
		var message notmodel.Message
		//setting up the notification
		message.CompanyID = voucher.CompanyID
		message.NodeID = voucher.NodeID
		message.CreatedBy = &voucher.CreatedBy
		message.RecepientID = voucher.UserID
		message.Message = voucher.Message
		message.URI = generateURI(createdVoucher, p.Engine.Envs[notification.AppURL])

		//now we fetch the notification serivce
		notificationRepo := notrepo.ProvideMessageRepo(p.Engine)
		notificationServ := ProvideNotMessageService(notificationRepo)

		//create the notification, if error occur will perform rollback
		if _, err = notificationServ.Create(message); err != nil {
			err = corerr.Tick(err, "E1455166", "notification cannot be created for voucher", message)
			db.Rollback()
		}

	}
	db.Commit()

	return
}

func generateURI(voucher eacmodel.Transaction, url string) (generatedURI string) {
	_ = "/#/main/accounting/voucher/journalVoucher/id/4"
	generatedURI = url
	generatedURI += "/#/main/accounting/voucher/journalVoucher/id/" + voucher.ID.ToString()

	return
}

// TxCreateVoucher is used for activating rollback in case of error. the following will store the slots in the
func (p *EacVoucherServ) TxCreateVoucher(db *gorm.DB, transaction eacmodel.Transaction,
	slots []eacmodel.Slot) (createdTransaction eacmodel.Transaction, err error) {

	if err = transaction.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1464263", "validation failed in tx-creating the voucher", transaction)
		return
	}

	now := time.Now()
	transaction.Hash = now.Format(consts.HashTimeLayout)

	if createdTransaction, err = p.Repo.TxCreate(db, transaction); err != nil {
		err = corerr.Tick(err, "E1445014", "transaction not created", transaction)
		return
	}

	for _, v := range slots {

		v.TransactionID = createdTransaction.ID
		v.TransactionID = createdTransaction.ID
		// GO and fix the TX create in the tempslot service, make sure it will be added to tempslot table
		if _, err = p.TempSlotServ.TxCreate(db, v); err != nil {
			err = corerr.Tick(err, "E1475598", "slot not saved in TxCreate", v)
			return
		}
	}
	return
}

// FindVoucherByID for getting transaction by it's id
func (p *EacVoucherServ) FindVoucherByID(fix types.FixedCol) (voucher eacmodel.Transaction, err error) {
	if voucher, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1466743", "can't fetch the transaction", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	if voucher.Slots, err = p.TempSlotServ.TemporarySlots(voucher.ID,
		voucher.CompanyID, voucher.NodeID); err != nil {
		err = corerr.Tick(err, "E1468950", "can't fetch the vouchers slots",
			fix.CompanyID, fix.NodeID, voucher.ID)
		return
	}

	return
}

// FindByYearCounter for getting transaction based on year counter
func (p *EacVoucherServ) FindByYearCounter(companyID uint64, transacstionType types.Enum, transactionYear string, counter string) (transaction eacmodel.Transaction, detailedSlots []eacmodel.DetailedSlots, err error) {
	//validation of the transaction type
	transaction.Type = transacstionType
	if err = transaction.Validate(coract.Fetch); err != nil {
		err = corerr.TickValidate(err, "REMINDER : ADD ERROR", "validation failed in fetching the transaction", transaction)
		return
	}

	//getting the  start and end of the input year
	var dateStart time.Time
	if dateStart, err = time.Parse("2006", transactionYear); err != nil {
		err = corerr.Tick(err, "E1415914", "can't convert the year to date", transactionYear)
		return
	}
	tempDate := dateStart.AddDate(0, 11, 0)
	y, m, _ := tempDate.Date()
	dateEnd := time.Date(y, m+1, 0, 23, 59, 59, 0, time.UTC)

	//parsing the counter to uint64
	var yearCounter uint64
	if yearCounter, err = strconv.ParseUint(counter, 10, 64); err != nil {
		err = corerr.Tick(err, "E1417726", "cannot parse year counter")
	}

	//getting the similar type of voucher
	similarType := transactiontype.SimilarType(transacstionType)
	if transaction, err = p.Repo.FindByYearCounter(companyID, yearCounter, transacstionType, similarType, dateStart, dateEnd); err != nil {
		err = limberr.New("transaction was not found by this year counter", "E1464285").
			Message(eacterm.Voucher, yearCounter).
			Custom(corerr.NotFoundErr).Build()
		return
	}

	eacTransactionServ, _ := fetchTransactionServ(p.Engine)

	//adding the slots

	//in case voucher has not been approved we fetch from temproray slots table
	if transaction.Status == transactionstatus.Unapproved {
		if transaction.Slots, err = p.TempSlotServ.TemporarySlots(transaction.ID, transaction.CompanyID, transaction.NodeID); err != nil {
			err = corerr.Tick(err, "E1490976", "can't fetch the voucher temporary slots",
				transaction.CompanyID, transaction.NodeID, transaction.ID)
			return
		}
	} else {

		if transaction.Slots, err = eacTransactionServ.SlotServ.TransactionSlot(transaction.ID, transaction.CompanyID, transaction.NodeID); err != nil {
			err = corerr.Tick(err, "E1443942", "can't fetch the voucher temporary slots",
				transaction.CompanyID, transaction.NodeID, transaction.ID)
			return
		}
	}

	//fetchimg the deatailed slots including account name and code
	fixedNode := types.FixedNode{
		ID:        transaction.ID,
		CompanyID: transaction.CompanyID,
		NodeID:    transaction.NodeID,
	}
	if detailedSlots, _, err = eacTransactionServ.FetchDetailedSlots(transaction.Slots, fixedNode, 0); err != nil {
		err = errors.New("cannot fetch the detailed slots")
		err = corerr.TickCustom(err, corerr.ValidationFailedErr, "E1441504", "cannot fetch the detailed slots")
		return
	}

	//now we add +1 to year counter
	// yearCounter++

	return
}

// ApproveVoucher is used for approving voucher and moving the slots from temp_eac_slots to eac_slots
func (p *EacVoucherServ) ApproveVoucher(voucher eacmodel.Transaction) (approvedVoucher, voucherBefore eacmodel.Transaction, err error) {
	fix := types.FixedCol{
		CompanyID: voucher.CompanyID,
		NodeID:    voucher.NodeID,
		ID:        voucher.ID,
	}

	// TODO: validation for journal-update

	//checking if voucher exists
	if voucherBefore, err = p.FindVoucherByID(fix); err != nil {
		err = corerr.Tick(err, "E1434122", "voucher  not found for approval", fix)
		return
	}

	voucher.CreatedBy = voucherBefore.CreatedBy
	voucher.CurrencyID = voucherBefore.CurrencyID
	voucher.Hash = voucherBefore.Hash
	voucher.Type = voucherBefore.Type
	voucher.Invoice = voucherBefore.Invoice
	voucher.YearCounter = voucherBefore.YearCounter
	voucher.CreatedAt = voucherBefore.CreatedAt
	voucher.PostDate = voucherBefore.PostDate
	voucher.YearCumulative = voucherBefore.YearCumulative
	//setting the voucher to approved
	voucher.Status = transactionstatus.Approved

	//check if transaciton has been already approved
	if voucherBefore.Status == transactionstatus.Approved {
		err = limberr.New("voucher has been already approved", "E1431610").
			Custom(corerr.Nil).
			Message(eacterm.ApproveVoucher).Build()
		return
	}

	transactionServ, _ := fetchTransactionServ(p.Engine)
	//newSlots, deletedSlots, updatedSlots := transactionServ.SlotServ.discreteSlots(voucherBefore.Slots, voucher.Slots)

	//glog.Debug(newSlots, deletedSlots, updatedSlots)

	db := p.Engine.DB.Begin()

	//updating the voucher
	if approvedVoucher, err = p.Repo.Save(voucher); err != nil {
		err = corerr.Tick(err, "E1487800", "transaction not saved")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in %v",
				"eac_transactions table"), "rollback recover approvevoucher")
			db.Rollback()
		}
	}()

	for _, v := range voucherBefore.Slots {
		fmt.Println("tempo slot:", v)
		v.TransactionID = approvedVoucher.ID
		v.PostDate = approvedVoucher.PostDate
		//delete slots from eac_temp_slots
		fix := types.FixedCol{
			ID:        v.ID,
			CompanyID: v.CompanyID,
			NodeID:    v.NodeID,
		}

		v.ID = 0
		//save slots to eac_slots from eac_temp_slots
		if _, err = transactionServ.SlotServ.TxCreate(db, v); err != nil {
			err = corerr.Tick(err, "E1436771", "slot not created in voucher approval", v)
			db.Rollback()
			return
		}

		if _, err = p.TempSlotServ.HardDelete(db, fix); err != nil {
			err = corerr.Tick(err, "E1415068", "slot not deleted in voucher approval", v)
			db.Rollback()
			return
		}

	}

	db.Commit()

	return
}

// VoucherUpdate is used for editing a voucher before approval
func (p *EacVoucherServ) VoucherUpdate(voucher eacmodel.Transaction) (voucherUpdated, voucherBefore eacmodel.Transaction, err error) {
	fix := types.FixedCol{
		CompanyID: voucher.CompanyID,
		NodeID:    voucher.NodeID,
		ID:        voucher.ID,
	}

	// validation for voucher update
	if err = voucher.Validate(coract.Update); err != nil {
		err = corerr.TickValidate(err, "E1420523", "validation failed in fetching the transaction", voucher)
		return
	}
	if voucherBefore, err = p.FindVoucherByID(fix); err != nil {
		err = corerr.Tick(err, "E1469607", "voucher not found for update", fix)
		return
	}
	// for _, v := range voucherBefore.Slots {
	// 	fmt.Println(v)
	// }
	if voucherBefore.Status != transactionstatus.Unapproved {
		err = corerr.Tick(err, "E1428939", "voucher has been already approved", fix)
		return
	}
	voucher.CreatedBy = voucherBefore.CreatedBy
	voucher.Hash = voucherBefore.Hash
	voucher.Type = voucherBefore.Type
	voucher.Invoice = voucherBefore.Invoice
	voucher.YearCounter = voucherBefore.YearCounter
	voucher.CreatedAt = voucherBefore.CreatedAt
	voucher.Status = voucherBefore.Status

	//using discreteslo() to differenitiate the new slots, deleted slots, and the updated slots from each othetr
	eacServ, _ := fetchTransactionServ(p.Engine)
	newSlots, deletedSlots, updatedSlots := eacServ.SlotServ.discreteSlots(voucherBefore.Slots, voucher.Slots)

	fmt.Printf("the new slots \n %v\n", newSlots)
	fmt.Printf("the updated slots slots \n %v\n", updatedSlots)
	fmt.Printf("the deleted slots \n %v\n", deletedSlots)

	db := p.Engine.DB.Begin()

	//updating the voucher
	if voucherUpdated, err = p.Repo.Save(voucher); err != nil {
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

	fmt.Println("voucher has been updated  ", voucher.CompanyID, voucher.NodeID)
	//adding the new slots
	for _, v := range newSlots {
		v.TransactionID = voucher.ID
		v.CompanyID = voucher.CompanyID
		v.NodeID = voucher.NodeID
		v.PostDate = voucher.PostDate
		v.CurrencyID = voucher.CurrencyID
		fmt.Println("THE ID :", v.ID, v.CompanyID)
		if _, err = p.TempSlotServ.TxCreate(db, v); err != nil {
			err = corerr.Tick(err, "E1469017", "slot not created in JournalUpdate", v)
			db.Rollback()
			return
		}
	}

	//updating the edited slots
	for _, v := range updatedSlots {

		v.TransactionID = voucher.ID
		v.PostDate = voucher.PostDate
		v.CompanyID = voucher.CompanyID
		v.NodeID = voucher.NodeID
		v.CurrencyID = voucher.CurrencyID
		if _, err = p.TempSlotServ.TxSave(db, v); err != nil {
			err = corerr.Tick(err, "E1458313", "slot not saved in JournalUpdate", v)
			db.Rollback()
			return
		}
	}

	//deleting the removed slots
	for _, v := range deletedSlots {
		v.TransactionID = voucher.ID
		v.CompanyID = voucher.CompanyID
		v.NodeID = voucher.NodeID
		adjust := v.Credit - v.Debit
		if err = p.TempSlotServ.Repo.TxRegulateBalancesSave(db, v, adjust); err != nil {
			err = corerr.Tick(err, "E1417880", "error in regulating the balance", v)
			db.Rollback()
			return
		}

		if err = p.TempSlotServ.Repo.TxDelete(db, v); err != nil {
			err = corerr.Tick(err, "E1489945", "error in deleting the slot", v)
			db.Rollback()
			return
		}
	}

	db.Commit()

	return
}

//function for fetching transaction service
func fetchTransactionServ(engine *core.Engine) (transactionServ EacTransactionServ, accountServ BasAccountServ) {
	transactionRepo := eacrepo.ProvideTransactionRepo(engine)
	slotRepo := eacrepo.ProvideSlotRepo(engine)
	currencyServ := ProvideEacCurrencyService(eacrepo.ProvideCurrencyRepo(engine))
	phoneServ := ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountServ = ProvideBasAccountService(basrepo.ProvideAccountRepo(engine), phoneServ)

	slotServ := ProvideEacSlotService(slotRepo, currencyServ, accountServ)
	transactionServ = ProvideEacTransactionService(transactionRepo, slotServ)

	return
}
