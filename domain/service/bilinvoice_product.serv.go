package service

import (
	"omono/domain/bill/bilmodel"
	"omono/domain/bill/bilrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"

	"gorm.io/gorm"
)

// BilInvoiceProductServ for injecting auth bilrepo
type BilInvoiceProductServ struct {
	Repo   bilrepo.InvoiceProductRepo
	Engine *core.Engine
}

// ProvideBilInvoiceProductService for invoiceProdcut is used in wire
func ProvideBilInvoiceProductService(p bilrepo.InvoiceProductRepo) BilInvoiceProductServ {
	return BilInvoiceProductServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// TxCreate create an invoiceProdcut with rollback availability
func (p *BilInvoiceProductServ) TxCreate(db *gorm.DB, invoiceProdcut bilmodel.InvoiceProduct) (createdInvoiceProduct bilmodel.InvoiceProduct, err error) {

	if err = invoiceProdcut.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7767635", "validation failed in creating the invoiceProdcut", invoiceProdcut)
		return
	}

	if createdInvoiceProduct, err = p.Repo.TxCreate(db, invoiceProdcut); err != nil {
		err = corerr.Tick(err, "E7796830", "invoiceProdcut not created", invoiceProdcut)
		return
	}

	return
}
