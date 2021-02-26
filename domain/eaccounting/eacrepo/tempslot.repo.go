package eacrepo

import (
	"errors"
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/base/message/basterm"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// TempSlotRepo for injecting engine
type TempSlotRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideTempSlotRepo is used in wire and initiate the Cols
func ProvideTempSlotRepo(engine *core.Engine) TempSlotRepo {
	return TempSlotRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(eacmodel.TempSlot{}), eacmodel.TempSlotTable),
	}
}

// FindByID finds the slot via its id
func (p *TempSlotRepo) FindByID(fix types.FixedCol) (slot eacmodel.Slot, err error) {
	err = p.Engine.ReadDB.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&slot).Error

	slot.ID = fix.ID
	err = p.dbError(err, "E1471037", slot, corterm.List)

	return
}

// List returns an array of slots
func (p *TempSlotRepo) List(params param.Param) (slots []eacmodel.Slot, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1426795").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1490266").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(eacmodel.TempSlotTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&slots).Error

	err = p.dbError(err, "E1474533", eacmodel.Slot{}, corterm.List)

	return
}

// Count of slots, mainly calls with List
func (p *TempSlotRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1428251").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(eacmodel.TempSlotTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1481282", eacmodel.Slot{}, corterm.List)
	return
}

// Save the slot, in case it is not exist create it
func (p *TempSlotRepo) Save(slot eacmodel.Slot) (u eacmodel.Slot, err error) {
	if err = p.Engine.DB.Table(eacmodel.TempSlotTable).Save(&slot).Error; err != nil {
		err = p.dbError(err, "E1484537", slot, corterm.Updated)
	}

	p.Engine.DB.Table(eacmodel.TempSlotTable).Where("id = ?", slot.ID).Find(&u)
	return
}

// TxSave is used for updating the slot via rollback
func (p *TempSlotRepo) TxSave(db *gorm.DB, slot eacmodel.Slot) (u eacmodel.Slot, err error) {
	if err = db.Table(eacmodel.TempSlotTable).Save(&slot).Error; err != nil {
		err = p.dbError(err, "E1435079", slot, corterm.Updated)
	}

	db.Table(eacmodel.TempSlotTable).Where("id = ?", slot.ID).Find(&u)
	return
}

// Create a slot
func (p *TempSlotRepo) Create(slot eacmodel.Slot) (u eacmodel.Slot, err error) {
	if err = p.Engine.DB.Table(eacmodel.TempSlotTable).Create(&slot).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1437304", slot, corterm.Created)
	}
	return
}

// TxCreate a slot by using rollback
func (p *TempSlotRepo) TxCreate(db *gorm.DB, slot eacmodel.Slot) (u eacmodel.Slot, err error) {
	if err = db.Table(eacmodel.TempSlotTable).Create(&slot).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1477259", slot, corterm.Created)
	}
	return
}

// TxDelete delete the slot via transaction connection
func (p *TempSlotRepo) TxDelete(db *gorm.DB, slot eacmodel.Slot) (err error) {
	now := time.Now()
	slot.DeletedAt = &now
	if err = db.Table(eacmodel.TempSlotTable).Save(&slot).Error; err != nil {
		err = p.dbError(err, "E1473780", slot, corterm.Deleted)
	}
	return
}

// HardDelete will make hard delete
func (p *TempSlotRepo) HardDelete(db *gorm.DB, slot eacmodel.Slot) (err error) {

	if err = db.Table(eacmodel.TempSlotTable).Delete(&slot).Error; err != nil {
		err = p.dbError(err, "E1470083", slot, corterm.Deleted)
	}
	return
}

// LastSlot returns the last slot before post_date
func (p *TempSlotRepo) LastSlot(slotIn eacmodel.Slot) (slot eacmodel.Slot, err error) {
	err = p.Engine.ReadDB.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND account_id = ? AND currency_id = ? AND post_date <= ? AND deleted_at IS NULL",
			slotIn.CompanyID, slotIn.AccountID, slotIn.CurrencyID, slotIn.PostDate).
		Order(" post_date DESC, id DESC ").
		Limit(1).
		Find(&slot).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	err = p.dbError(err, "E1471037", slot, corterm.List)

	return
}

// TxLastSlot returns the last slot before post_date
func (p *TempSlotRepo) TxLastSlot(db *gorm.DB, slotIn eacmodel.Slot) (slot eacmodel.Slot, err error) {
	// err = db.Clauses(clause.Locking{Strength: "UPDATE"}).Table(eacmodel.TempSlotTable).
	err = db.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND account_id = ? AND currency_id = ? AND deleted_at IS NULL",
			slotIn.CompanyID, slotIn.AccountID, slotIn.CurrencyID).
		Order(" post_date DESC, id DESC ").
		Limit(1).
		Find(&slot).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	slot.CompanyID = slotIn.CompanyID
	slot.NodeID = slotIn.NodeID
	slot.AccountID = slotIn.AccountID
	slot.CurrencyID = slotIn.CurrencyID

	err = p.dbError(err, "E1452166", slot, corterm.List)

	return
}

// LastSlotWithID returns the last slot before post_date
func (p *TempSlotRepo) LastSlotWithID(slotIn eacmodel.Slot) (slot eacmodel.Slot, err error) {
	err = p.Engine.ReadDB.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND account_id = ? AND currency_id = ? AND post_date <= ? AND id < ? AND deleted_at IS NULL",
			slotIn.CompanyID, slotIn.AccountID, slotIn.CurrencyID, slotIn.PostDate, slotIn.ID).
		Order(" post_date DESC, id DESC ").
		Limit(1).
		Find(&slot).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	err = p.dbError(err, "E1495037", slot, corterm.List)

	return
}

// RegulateBalances will adjust all balances after the post_date
func (p *TempSlotRepo) RegulateBalances(slotIn eacmodel.Slot, adjust float64) (err error) {
	err = p.Engine.DB.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND account_id = ? AND currency_id = ? AND post_date > ?",
			slotIn.CompanyID, slotIn.AccountID, slotIn.CurrencyID, slotIn.PostDate).
		Update("balance", gorm.Expr("balance + ?", adjust)).Error
	return
}

// TxRegulateBalances will adjust all balances after the post_date
func (p *TempSlotRepo) TxRegulateBalances(db *gorm.DB, slotIn eacmodel.Slot, adjust float64) (err error) {
	err = db.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND account_id = ? AND currency_id = ? AND post_date > ?",
			slotIn.CompanyID, slotIn.AccountID, slotIn.CurrencyID, slotIn.PostDate).
		Update("balance", gorm.Expr("balance + ?", adjust)).Error
	return
}

// RegulateBalancesSave will adjust all balances after and equal the post_date
func (p *TempSlotRepo) RegulateBalancesSave(slotIn eacmodel.Slot, adjust float64) (err error) {
	err = p.Engine.DB.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND account_id = ? AND currency_id = ? AND post_date >= ? AND id > ?",
			slotIn.CompanyID, slotIn.AccountID, slotIn.CurrencyID, slotIn.PostDate, slotIn.ID).
		Update("balance", gorm.Expr("balance + ?", adjust)).Error
	return
}

// TxRegulateBalancesSave will adjust all balances after and equal the post_date
func (p *TempSlotRepo) TxRegulateBalancesSave(db *gorm.DB, slotIn eacmodel.Slot, adjust float64) (err error) {
	err = db.Table(eacmodel.TempSlotTable).
		Where("company_id = ? AND account_id = ? AND currency_id = ? AND post_date >= ? AND id > ?",
			slotIn.CompanyID, slotIn.AccountID, slotIn.CurrencyID, slotIn.PostDate, slotIn.ID).
		Update("balance", gorm.Expr("balance + ?", adjust)).Error
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *TempSlotRepo) dbError(err error, code string, slot eacmodel.Slot, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, slot.ID, eacterm.Transaction)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(eacterm.Transaction), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}

// TxUpdateBalance will insert or update the balance based on the input data
func (p *TempSlotRepo) TxUpdateBalance(db *gorm.DB, b eacmodel.Balance) (err error) {
	query := fmt.Sprintf("INSERT INTO eac_balances(company_id,node_id, account_id,currency_id, balance) VALUES(%v, %v,%v,%v,%f) ON DUPLICATE KEY UPDATE balance = %f;",
		b.CompanyID, b.NodeID, b.AccountID, b.CurrencyID, b.Balance, b.Balance)

	err = db.Table(eacmodel.BalanceTable).
		Exec(query).Error
	err = p.dbError(err, "E1450987", eacmodel.Slot{}, eacterm.UpdateBalance)

	return
}
