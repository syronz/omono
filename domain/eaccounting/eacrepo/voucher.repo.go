package eacrepo

import (
	"errors"
	"github.com/syronz/dict"
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
	"omono/pkg/helper"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// VoucherRepo for injecting engine
type VoucherRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideVoucherRepo is used in wire and initiate the Cols
func ProvideVoucherRepo(engine *core.Engine) VoucherRepo {
	return VoucherRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(eacmodel.Transaction{}), eacmodel.TransactionTable),
	}
}

// FindByID finds the transaction via its id
func (p *VoucherRepo) FindByID(fix types.FixedCol) (transaction eacmodel.Transaction, err error) {
	err = p.Engine.ReadDB.Table(eacmodel.TransactionTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&transaction).Error

	transaction.ID = fix.ID
	err = p.dbError(err, "E1439159", transaction, corterm.List)

	return
}

// FindByYearCounter finds the transaction based on the year counter
func (p *VoucherRepo) FindByYearCounter(companyID, yearCounter uint64, transactionType, similarType types.Enum, startDate, endDate time.Time) (tranasction eacmodel.Transaction, err error) {

	err = p.Engine.ReadDB.Table(eacmodel.TransactionTable).
		Where("company_id = ? AND year_counter=? AND type IN (?,?) AND post_date BETWEEN ? AND ?", companyID, yearCounter, transactionType, similarType, startDate, endDate).
		First(&tranasction).Error

	err = p.dbError(err, "E1485524", tranasction, corterm.List)

	return
}

// List returns an array of transactions
func (p *VoucherRepo) List(params param.Param) (transactions []eacmodel.Transaction, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1449108").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1444986").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(eacmodel.TransactionTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&transactions).Error

	err = p.dbError(err, "E1454073", eacmodel.Transaction{}, corterm.List)

	return
}

// LastYearCounter is used for returning the last-year for generating invoice
func (p *VoucherRepo) LastYearCounter(db *gorm.DB, companyID uint64, tType types.Enum, lastYearDay time.Time) (transaction eacmodel.Transaction, err error) {
	err = db.Table(eacmodel.TransactionTable).
		Where("company_id = ? AND type = ? AND post_date < ?", companyID, tType, lastYearDay).
		Order("year_counter DESC").
		Last(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
			return
		}
		err = limberr.Take(err, "").Custom(corerr.InternalServerErr).Build()
		return
	}

	return
}

// LastYearCumulative is used for returning the year_cumulative for generating invoice
func (p *VoucherRepo) LastYearCumulative(db *gorm.DB, companyID uint64, lastYearDay time.Time) (transaction eacmodel.Transaction, err error) {
	err = db.Table(eacmodel.TransactionTable).
		Where("company_id = ? AND post_date < ?", companyID, lastYearDay).
		Order("year_cumulative DESC").
		Last(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
			return
		}
		err = limberr.Take(err, "E1445086").Custom(corerr.InternalServerErr).Build()
		return
	}

	return
}

// Count of transactions, mainly calls with List
// func (p *VoucherRepo) Count(params param.Param) (count int64, err error) {
// 	var whereStr string
// 	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
// 		err = limberr.Take(err, "E1436521").Custom(corerr.ValidationFailedErr).Build()
// 		return
// 	}

// 	err = p.Engine.ReadDB.Table(eacmodel.TransactionTable).
// 		Where(whereStr).
// 		Count(&count).Error

// 	err = p.dbError(err, "E1465399", eacmodel.Transaction{}, corterm.List)
// 	return
// }

// Save the transaction, in case it is not exist create it
func (p *VoucherRepo) Save(transaction eacmodel.Transaction) (u eacmodel.Transaction, err error) {
	if err = p.Engine.DB.Table(eacmodel.TransactionTable).Save(&transaction).Error; err != nil {
		err = p.dbError(err, "E1420013", transaction, corterm.Updated)
	}

	p.Engine.DB.Table(eacmodel.TransactionTable).Where("id = ?", transaction.ID).Find(&u)
	return
}

// Create a transaction
func (p *VoucherRepo) Create(transaction eacmodel.Transaction) (u eacmodel.Transaction, err error) {
	if err = p.Engine.DB.Table(eacmodel.TransactionTable).Create(&transaction).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1413616", transaction, corterm.Created)
	}
	return
}

// TxCreate a transaction
func (p *VoucherRepo) TxCreate(db *gorm.DB, transaction eacmodel.Transaction) (u eacmodel.Transaction, err error) {
	if err = db.Table(eacmodel.TransactionTable).Create(&transaction).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1487269", transaction, corterm.Created)
	}
	return
}

// Delete the transaction
func (p *VoucherRepo) Delete(transaction eacmodel.Transaction) (err error) {
	if err = p.Engine.DB.Table(eacmodel.TransactionTable).Delete(&transaction).Error; err != nil {
		err = p.dbError(err, "E1474760", transaction, corterm.Deleted)
	}
	return
}

// TxDelete the transaction via rollback facility
func (p *VoucherRepo) TxDelete(db *gorm.DB, transaction eacmodel.Transaction) (err error) {
	now := time.Now()
	transaction.DeletedAt = &now
	if err = db.Table(eacmodel.TransactionTable).Save(&transaction).Error; err != nil {
		err = p.dbError(err, "E1426674", transaction, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *VoucherRepo) dbError(err error, code string, transaction eacmodel.Transaction, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, transaction.ID, eacterm.Transactions)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(eacterm.Transaction), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(eacterm.Transaction), transaction.Hash).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, transaction.Hash)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
