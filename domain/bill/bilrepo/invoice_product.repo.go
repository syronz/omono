package bilrepo

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/bill/bilmodel"
	"omono/domain/bill/bilterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/pkg/helper"
	"reflect"

	"gorm.io/gorm"
)

// InvoiceProductRepo for injecting engine
type InvoiceProductRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideInvoiceProductRepo is used in wire and initiate the Cols
func ProvideInvoiceProductRepo(engine *core.Engine) InvoiceProductRepo {
	return InvoiceProductRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(bilmodel.InvoiceProduct{}), bilmodel.InvoiceProductTable),
	}
}

// TxCreate create an invoiceProduct via rollback
func (p *InvoiceProductRepo) TxCreate(db *gorm.DB, invoiceProduct bilmodel.InvoiceProduct) (u bilmodel.InvoiceProduct, err error) {
	if err = db.Table(bilmodel.InvoiceProductTable).Create(&invoiceProduct).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E7768701", invoiceProduct, corterm.Created)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *InvoiceProductRepo) dbError(err error, code string, invoiceProduct bilmodel.InvoiceProduct, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, invoiceProduct.ID, bilterm.Invoice)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(bilterm.Invoice),
				dict.R(bilterm.Invoice), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(bilterm.Invoice), invoiceProduct.InvoiceID).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, invoiceProduct.InvoiceID)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
