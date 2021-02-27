package bilrepo

import (
	"errors"
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/bill/bilmodel"
	"omono/domain/bill/bilterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/helper"
	"reflect"

	"gorm.io/gorm"
)

// InvoiceRepo for injecting engine
type InvoiceRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideInvoiceRepo is used in wire and initiate the Cols
func ProvideInvoiceRepo(engine *core.Engine) InvoiceRepo {
	return InvoiceRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(bilmodel.Invoice{}), bilmodel.InvoiceTable),
	}
}

// FindByID finds the invoice via its id
func (p *InvoiceRepo) FindByID(fix types.FixedCol) (invoice bilmodel.Invoice, err error) {
	err = p.Engine.ReadDB.Table(bilmodel.InvoiceTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&invoice).Error

	invoice.ID = fix.ID
	err = p.dbError(err, "E7758754", invoice, corterm.List)

	return
}

// TxLastInvoice returns back the last invoice for a location
func (p *InvoiceRepo) TxLastInvoice(db *gorm.DB, companyID uint64, storeID types.RowID) (invoice bilmodel.Invoice, err error) {
	err = db.Table(bilmodel.InvoiceTable).
		Where("company_id = ? AND store_id = ?", companyID, storeID).
		Last(&invoice).Error

	switch {
	case err == nil:
		return
	case errors.Is(err, gorm.ErrRecordNotFound):
		err = nil
		return
	default:
		err = p.dbError(err, "E7714861", invoice, corterm.List)
	}

	return
}

// List returns an array of invoices
func (p *InvoiceRepo) List(params param.Param) (invoices []bilmodel.Invoice, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E7795269").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E7791380").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(bilmodel.InvoiceTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&invoices).Error

	err = p.dbError(err, "E7740085", bilmodel.Invoice{}, corterm.List)

	return
}

// Count of invoices, mainly calls with List
func (p *InvoiceRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E7795741").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(bilmodel.InvoiceTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E7780267", bilmodel.Invoice{}, corterm.List)
	return
}

// TxSave the invoice, in case it is not exist create it
func (p *InvoiceRepo) TxSave(db *gorm.DB, invoice bilmodel.Invoice) (u bilmodel.Invoice, err error) {
	if err = db.Table(bilmodel.InvoiceTable).Save(&invoice).Error; err != nil {
		err = p.dbError(err, "E7731161", invoice, corterm.Updated)
	}

	p.Engine.DB.Table(bilmodel.InvoiceTable).Where("id = ?", invoice.ID).Find(&u)
	return
}

// TxCreate create an invoice with rollback availability
func (p *InvoiceRepo) TxCreate(db *gorm.DB, invoice bilmodel.Invoice) (u bilmodel.Invoice, err error) {
	if err = db.Table(bilmodel.InvoiceTable).Create(&invoice).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E7789853", invoice, corterm.Created)
	}
	return
}

// Delete the invoice
func (p *InvoiceRepo) Delete(invoice bilmodel.Invoice) (err error) {
	if err = p.Engine.DB.Table(bilmodel.InvoiceTable).Delete(&invoice).Error; err != nil {
		err = p.dbError(err, "E7755143", invoice, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *InvoiceRepo) dbError(err error, code string, invoice bilmodel.Invoice, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, invoice.ID, bilterm.Invoices)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(bilterm.Invoices),
				dict.R(bilterm.Invoice), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(bilterm.Invoice), invoice.Invoice).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, invoice.Invoice)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
