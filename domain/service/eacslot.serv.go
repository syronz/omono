package service

import (
	"errors"
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/base/basmodel"
	"omono/domain/base/enum/accountstatus"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// EacSlotServ for injecting auth eacrepo
type EacSlotServ struct {
	Repo         eacrepo.SlotRepo
	Engine       *core.Engine
	CurrencyServ EacCurrencyServ
	AccountServ  BasAccountServ
}

// ProvideEacSlotService for slot is used in wire
func ProvideEacSlotService(p eacrepo.SlotRepo, currencyServ EacCurrencyServ,
	accountServ BasAccountServ) EacSlotServ {
	return EacSlotServ{
		Repo:         p,
		Engine:       p.Engine,
		CurrencyServ: currencyServ,
		AccountServ:  accountServ,
	}
}

// FindByID for getting slot by it's id
func (p *EacSlotServ) FindByID(fix types.FixedCol) (slot eacmodel.Slot, err error) {
	if slot, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1410775", "can't fetch the slot", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	return
}

// List of slots, it support pagination and search and return back count
func (p *EacSlotServ) List(params param.Param) (slots []eacmodel.Slot,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" eac_slots.company_id = '%v' ", params.CompanyID)
	}

	if slots, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in slots list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in slots count")
	}

	return
}

// TransactionSlot is used inside Transaction.FindByID
func (p *EacSlotServ) TransactionSlot(transactionID types.RowID,
	companyID, nodeID uint64) (slots []eacmodel.Slot,
	err error) {

	params := param.New()
	params.Limit = consts.MaxRowsCount
	params.PreCondition = fmt.Sprintf(" eac_slots.company_id = '%v' AND eac_slots.transaction_id = '%v' AND eac_slots.node_id = '%v' AND eac_slots.deleted_at is NULL",
		companyID, transactionID, nodeID)
	params.Order = " eac_slots.post_date asc, eac_slots.id asc "

	if slots, err = p.Repo.List(params); err != nil {
		return
	}

	return
}

// discreteSlots is used for distinguish between different type of slot based on old slots, it is
// used in update journal
func (p *EacSlotServ) discreteSlots(oldSlots, rawSlots []eacmodel.Slot) (newSlots, deletedSlots, updatedSlots []eacmodel.Slot) {

	// var newSlots, updatedSlots, deletedSlots []eacmodel.Slot

	type Compare struct {
		Old     eacmodel.Slot
		Updated eacmodel.Slot
	}

	mapSlots := make(map[types.RowID]Compare, 1)

	// new slots
	for _, v := range rawSlots {
		fmt.Println("the id for each slot:", v.ID)
		if v.ID == 0 {
			newSlots = append(newSlots, v)
		} else {
			compare := Compare{
				Updated: v,
			}
			mapSlots[v.ID] = compare
		}
	}

	// old slots
	for _, v := range oldSlots {
		enteredCompare, ok := mapSlots[v.ID]
		if ok {
			enteredCompare.Old = v
			mapSlots[v.ID] = enteredCompare
			k := enteredCompare.Updated

			var oldDescription, updatedDescription string
			if v.Description != nil {
				oldDescription = *v.Description
			}

			if k.Description != nil {
				updatedDescription = *k.Description
			}

			if v.AccountID != k.AccountID ||
				v.Debit != k.Debit ||
				v.Credit != k.Credit ||
				!v.PostDate.Equal(k.PostDate) ||
				// (v.PostDate.Year() != k.PostDate.Year() && v.PostDate.YearDay() != k.PostDate.YearDay()) ||
				oldDescription != updatedDescription ||
				v.CurrencyID != k.CurrencyID {
				updatedSlots = append(updatedSlots, k)
			}
		} else {
			compare := Compare{
				Old: v,
			}
			mapSlots[v.ID] = compare
			deletedSlots = append(deletedSlots, v)
		}
	}

	return
}

func (p *EacSlotServ) duplicateAccounts(slots, newSlots []eacmodel.Slot) (err error) {
	var uniqueAcc []types.RowID
	IsUnique := true
	for _, v := range slots {
		for _, k := range uniqueAcc {
			if k == v.AccountID {
				IsUnique = false
				break
			}

		}
		if IsUnique {
			uniqueAcc = append(uniqueAcc, v.AccountID)
		}
		IsUnique = true
	}

	if len(uniqueAcc) != len(slots) {
		err = corerr.Tick(err, "E1428907", "duplicate accounts exists")
		return
	}

	//validation for newSlots
	for _, v := range newSlots {
		for _, k := range uniqueAcc {
			if k == v.AccountID {
				err = errors.New("duplicate accounts exists")
				err = corerr.TickCustom(err, corerr.ValidationFailedErr, "E1428907", "duplicate accounts exists")
				return
			}

		}
	}
	return
}

// Create a slot
func (p *EacSlotServ) Create(slot eacmodel.Slot) (createdSlot eacmodel.Slot, err error) {
	if err = slot.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1437242", "validation failed in creating the slot", slot)
		return
	}

	var lastSlot eacmodel.Slot
	if lastSlot, err = p.Repo.LastSlot(slot); err != nil {
		err = corerr.Tick(err, "E1445069", "last slot not found", slot)
		return
	}

	adjust := slot.Debit - slot.Credit
	slot.Balance = lastSlot.Balance + adjust

	if createdSlot, err = p.Repo.Create(slot); err != nil {
		err = corerr.Tick(err, "E1434523", "slot not created", slot)
		return
	}

	if err = p.Repo.RegulateBalances(slot, adjust); err != nil {
		err = corerr.Tick(err, "E1466626", "regulate balances faced error in create", slot, adjust)
		return
	}

	return
}

// TxCreate a slot is used for activating rollback
func (p *EacSlotServ) TxCreate(db *gorm.DB, slot eacmodel.Slot) (createdSlot eacmodel.Slot, err error) {
	if err = slot.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1457207", "validation failed in creating the slot", slot)
		return
	}

	//validation for the accounts whether it's readonly or not
	if err = p.AccountReadOnlyValidation(slot); err != nil {
		err = corerr.Tick(err, "E1490575", "Account is readonly", slot)
		return
	}

	var lastSlot eacmodel.Slot
	if lastSlot, err = p.Repo.TxLastSlot(db, slot); err != nil {
		err = corerr.Tick(err, "E1428648", "last slot not found", slot)
		return
	}

	fmt.Println("this is last slot")

	adjust := slot.Debit - slot.Credit
	slot.Balance = lastSlot.Balance + adjust

	var account basmodel.Account
	fix := types.FixedNode{
		ID:        slot.AccountID,
		CompanyID: slot.CompanyID,
	}
	if account, err = p.AccountServ.TxFindAccountStatus(db, fix); err != nil {
		err = corerr.Tick(err, "E1491733", "error in getting account for creating a slot", fix)
		return
	}

	if account.Status == accountstatus.Inactive {
		err = limberr.New("account is inactive", "E1410609").
			Message("Account is inactive").
			Custom(corerr.ForbiddenErr).Build()
		return
	}

	if createdSlot, err = p.Repo.TxCreate(db, slot); err != nil {
		err = corerr.Tick(err, "E1424374", "slot not created", slot)
		return
	}

	if err = p.Repo.TxRegulateBalances(db, slot, adjust); err != nil {
		err = corerr.Tick(err, "E1452307", "regulate balances faced error in create", slot, adjust)
		return
	}

	if err = p.TxUpdateBalance(db, slot); err != nil {
		err = corerr.Tick(err, "E1455024", "update balance for account faced problem", slot)
		return
	}

	return
}

// Reset will remove affect of transaction on journal, similar to delete but don't delete the
// records
func (p *EacSlotServ) Reset(slot eacmodel.Slot) (err error) {
	adjust := slot.Credit - slot.Debit
	slot.Debit = 0
	slot.Credit = 0
	slot.Balance = 0
	if slot, err = p.Repo.Save(slot); err != nil {
		err = corerr.Tick(err, "E1445231", "error in resetting the slot", slot)
		return
	}

	if err = p.Repo.RegulateBalancesSave(slot, adjust); err != nil {
		err = corerr.Tick(err, "E1491272", "regulate balances faced error in reset", slot, adjust)
		return
	}

	return
}

// TxReset is used for reset the transaction's slot without deleting via rollback
func (p *EacSlotServ) TxReset(db *gorm.DB, slot eacmodel.Slot) (err error) {
	adjust := slot.Credit - slot.Debit
	slot.Debit = 0
	slot.Credit = 0
	slot.Balance = 0
	if slot, err = p.Repo.TxSave(db, slot); err != nil {
		err = corerr.Tick(err, "E1433784", "error in resetting the slot", slot)
		return
	}

	if err = p.Repo.TxRegulateBalancesSave(db, slot, adjust); err != nil {
		err = corerr.Tick(err, "E1443695", "regulate balances faced error in reset", slot, adjust)
		return
	}

	return
}

// Save a slot, if it is exist update it, if not create it
func (p *EacSlotServ) Save(slot eacmodel.Slot) (savedSlot eacmodel.Slot, err error) {
	if err = slot.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1445816", corerr.ValidationFailed, slot)
		return
	}

	fix := types.FixedCol{
		CompanyID: slot.CompanyID,
		NodeID:    slot.NodeID,
		ID:        slot.ID,
	}

	var oldSlot eacmodel.Slot

	if oldSlot, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1482909", "slot not found", slot)
		return
	}

	p.Reset(oldSlot)

	var lastSlot eacmodel.Slot
	if lastSlot, err = p.Repo.LastSlotWithID(slot); err != nil {
		err = corerr.Tick(err, "E1433617", "last slot not found in save transaction", slot)
		return
	}

	adjust := slot.Debit - slot.Credit
	slot.Balance = lastSlot.Balance + adjust

	if slot, err = p.Repo.Save(slot); err != nil {
		err = corerr.Tick(err, "E1475746", "error in saving the slot", slot)
		return
	}

	if err = p.Repo.RegulateBalancesSave(slot, adjust); err != nil {
		err = corerr.Tick(err, "E1454858", "regulate balances faced error in save", slot, adjust)
		return
	}

	return
}

// TxSave is used for saving an existing slot via rollback
func (p *EacSlotServ) TxSave(db *gorm.DB, slot eacmodel.Slot) (savedSlot eacmodel.Slot, err error) {
	//slot validation from slot model
	if err = slot.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1475432", corerr.ValidationFailed, slot)
		return
	}
	//validation for the accounts whether it's readonly or not
	if err = p.AccountReadOnlyValidation(slot); err != nil {
		err = corerr.Tick(err, "E1490575", "Account is readonly", slot)
		return
	}

	fix := types.FixedCol{
		CompanyID: slot.CompanyID,
		NodeID:    slot.NodeID,
		ID:        slot.ID,
	}

	var oldSlot eacmodel.Slot

	if oldSlot, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1428966", "slot not found", slot)
		return
	}

	var needReset bool
	if oldSlot.Debit != slot.Debit ||
		oldSlot.Credit != slot.Credit ||
		oldSlot.Balance != slot.Balance ||
		oldSlot.CurrencyID != slot.CurrencyID ||
		oldSlot.AccountID != slot.AccountID ||
		!oldSlot.PostDate.Equal(slot.PostDate) {
		needReset = true
	}
	adjust := slot.Debit - slot.Credit

	if needReset {
		p.TxReset(db, oldSlot)

		var lastSlot eacmodel.Slot
		if lastSlot, err = p.Repo.LastSlotWithID(slot); err != nil {
			err = corerr.Tick(err, "E1416071", "last slot not found in save transaction", slot)
			return
		}

		slot.Balance = lastSlot.Balance + adjust
	}

	//adding created_at ,since in update the created at will become null
	if oldSlot.CreatedAt != nil {
		slot.CreatedAt = oldSlot.CreatedAt
	}
	if slot, err = p.Repo.TxSave(db, slot); err != nil {
		err = corerr.Tick(err, "E1499071", "error in saving the slot", slot)
		return
	}

	if needReset {
		if err = p.Repo.TxRegulateBalancesSave(db, slot, adjust); err != nil {
			err = corerr.Tick(err, "E1481480", "regulate balances faced error in save", slot, adjust)
			return
		}

		if err = p.TxUpdateBalance(db, slot); err != nil {
			err = corerr.Tick(err, "E1460419", "update balance for account faced problem in save slot", slot)
			return
		}
	}

	return
}

// TxDelete is used for deleting a slot via rollback
func (p *EacSlotServ) TxDelete(db *gorm.DB, fix types.FixedCol) (slot eacmodel.Slot, err error) {
	if slot, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1498868", "slot not found for deleting")
		return
	}

	if err = p.Repo.TxDelete(db, slot); err != nil {
		err = corerr.Tick(err, "E1498728", "slot not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *EacSlotServ) Excel(params param.Param) (slots []eacmodel.Slot, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", eacmodel.SlotTable)

	if slots, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1485891", "cant generate the excel list for slots")
		return
	}

	return
}

// TxUpdateBalance is set the last balance from the eac_slots and put it to the eac_balances
func (p *EacSlotServ) TxUpdateBalance(db *gorm.DB, slot eacmodel.Slot) (err error) {
	if slot, err = p.Repo.TxLastSlot(db, slot); err != nil {
		err = corerr.Tick(err, "E1490538", "cant fetch last slot for updating the balance")
		return
	}

	var balance eacmodel.Balance
	balance.CompanyID = slot.CompanyID
	balance.NodeID = slot.NodeID
	balance.AccountID = slot.AccountID
	balance.Balance = slot.Balance
	balance.CurrencyID = slot.CurrencyID

	return p.Repo.TxUpdateBalance(db, balance)
}

//AccountReadOnlyValidation will validate whether an account's status is active and whether it's readonly or no
func (p *EacSlotServ) AccountReadOnlyValidation(slot eacmodel.Slot) (err error) {

	var acc basmodel.Account
	fix := types.FixedNode{
		ID:        slot.AccountID,
		CompanyID: slot.CompanyID,
		NodeID:    slot.NodeID,
	}
	if acc, err = p.AccountServ.FindByID(fix); err != nil {
		err = errors.New("Account was not found")
		return err
	}

	if acc.ReadOnly {
		err = errors.New("Account is read only in slot")
	}
	return

}
