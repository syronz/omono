package service

import (
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/bill"
	"omono/domain/bill/bilmodel"
	"omono/domain/bill/bilrepo"
	"omono/domain/bill/enum/invoicestatus"
	"omono/domain/bill/enum/pricemode"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/location/enum/storestatus"
	"omono/domain/location/locrepo"
	"omono/domain/location/locterm"
	"omono/domain/material/enum/productstatus"
	"omono/domain/material/matmodel"
	"omono/domain/material/matrepo"
	"omono/domain/material/matterm"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/glog"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BilInvoiceServ for injecting auth bilrepo
type BilInvoiceServ struct {
	Repo   bilrepo.InvoiceRepo
	Engine *core.Engine
}

// ProvideBilInvoiceService for invoice is used in wire
func ProvideBilInvoiceService(p bilrepo.InvoiceRepo) BilInvoiceServ {
	return BilInvoiceServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting invoice by it's id
func (p *BilInvoiceServ) FindByID(fix types.FixedCol) (invoice bilmodel.Invoice, err error) {
	if invoice, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7799821", "can't fetch the invoice", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	return
}

// List of invoices, it support pagination and search and return back count
func (p *BilInvoiceServ) List(params param.Param) (invoices []bilmodel.Invoice,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" bil_invoices.company_id = '%v' ", params.CompanyID)
	}

	if invoices, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in invoices list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in invoices count")
	}

	return
}

// Create a invoice
func (p *BilInvoiceServ) Create(invoice bilmodel.Invoice) (createdInvoice bilmodel.Invoice, err error) {

	invoice.Status = invoicestatus.New
	invoice.Year = time.Now().Year()

	if err = invoice.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7739340", "validation failed in creating the invoice", invoice)
		return
	}

	// store fetched and check activation and users inside the list
	storeFix := types.ExtractFixedCol(invoice)
	storeFix.ID = invoice.StoreID

	storeServ := ProvideLocStoreService(locrepo.ProvideStoreRepo(p.Engine))
	if invoice.Store, err = storeServ.FindByID(storeFix); err != nil {
		err = corerr.Tick(err, "E7723745", "store not found in creating invoice", invoice)
		return
	}

	if invoice.Store.Status != storestatus.Active {
		err = limberr.New("store is inactive", "E7761102").
			Message(corerr.VisInactive, dict.R(locterm.Store)).
			Custom(corerr.ForbiddenErr).
			Domain(bill.Domain).
			Build()
		return
	}

	// check if user defined in the store
	if err = storeServ.IsUserInStore(storeFix, invoice.CreatedBy); err != nil {
		return
	}

	// get currency rate
	rateServ := ProvideEacRateService(eacrepo.ProvideRateRepo(p.Engine))
	if invoice.Rate, err = rateServ.GetRate(invoice.Store.CityID, invoice.CurrencyID); err != nil {
		err = corerr.Tick(err, "E7788941", "can't fetch currency rate", invoice.Store.CityID, invoice.CurrencyID)
		return
	}
	invoice.CurrencyRate = invoice.Rate.Rate

	glog.Debug(invoice)

	/*
		TRANSACTION PART
	*/
	db := p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"bil_invoices table"), "rollback recover create invoice")
			db.Rollback()
		}
	}()

	if invoice.YearCounter, invoice.YearCumulative, invoice.Invoice, err =
		p.TxInvoiceGenerator(db, invoice); err != nil {
		err = corerr.Tick(err, "E7741555", "invoice number can not be generated", invoice)
		db.Rollback()
		return
	}

	if createdInvoice, err = p.Repo.TxCreate(db, invoice); err != nil {
		err = corerr.Tick(err, "E7722102", "invoice not created", invoice)
		db.Rollback()
		return
	}

	// insert prodcuts to the bil_invoice_products
	invoiceProductServ := ProvideBilInvoiceProductService(bilrepo.ProvideInvoiceProductRepo(p.Engine))
	productServ := ProvideMatProductService(matrepo.ProvideProductRepo(p.Engine))
	var total float64

	for i, v := range invoice.Products {
		var product matmodel.Product
		if product, err = productServ.FindByRowID(v.ProductID); err != nil {
			err = corerr.Tick(err, "E7790863", "product not found for creating invoice, invoice_id & product:", invoice.ID, v)
			db.Rollback()
			return
		}

		if product.Status != productstatus.Active {
			err = limberr.New("product is inactive").
				Message(matterm.ProductVIsInactive, product.Name).
				Custom(corerr.ForbiddenErr).
				Domain(bill.Domain).
				Build()
			err = corerr.Tick(err, "E7789417", fmt.Sprintf("product is not active, produt: %v - %s", product.ID, product.Name))
			db.Rollback()
			return
		}

		if err = p.getPriceForMode(i, &invoice, &product); err != nil {
			db.Rollback()
			return
		}

		var discount float64
		if invoice.Products[i].Discount != nil {
			discount = *invoice.Products[i].Discount
		}
		total += invoice.Products[i].Price - discount

		invoice.Products[i].InvoiceID = createdInvoice.ID
		if invoice.Products[i], err = invoiceProductServ.TxCreate(db, invoice.Products[i]); err != nil {
			db.Rollback()
			return
		}

	} // products for

	createdInvoice.Total = total
	if _, err = p.TxSave(db, createdInvoice); err != nil {
		err = corerr.Tick(err, "E7772071", fmt.Sprintf("can't save total to the invoice: %v", invoice.Invoice))
		db.Rollback()
		return
	}

	db.Commit()

	createdInvoice.Products = invoice.Products

	return
}

func (p *BilInvoiceServ) getPriceForMode(i int, invoice *bilmodel.Invoice, product *matmodel.Product) (err error) {
	errTmp := limberr.New("product is inactive").
		Message(matterm.PriceNotAddedForThisModeVInProdcutV, invoice.PriceMode, product.Name).
		Custom(corerr.PreDataInsertedErr).
		Domain(bill.Domain).
		Build()

	switch invoice.PriceMode {
	case pricemode.Whole:
		if product.PriceWhole != nil {
			invoice.Products[i].Price = *product.PriceWhole
		} else {
			err = corerr.Tick(errTmp, "E7719850", fmt.Sprintf("price_mode %v not defined for product: %v - %s",
				pricemode.Whole, product.ID, product.Name))
			return
		}

	case pricemode.VIP:
		if product.PriceVIP != nil {
			invoice.Products[i].Price = *product.PriceVIP
		} else {
			err = corerr.Tick(errTmp, "E7735988", fmt.Sprintf("price_mode %v not defined for product: %v - %s",
				pricemode.VIP, product.ID, product.Name))
			return
		}

	case pricemode.Distributor:
		if product.PriceDistributor != nil {
			invoice.Products[i].Price = *product.PriceDistributor
		} else {
			err = corerr.Tick(errTmp, "E7714769", fmt.Sprintf("price_mode %v not defined for product: %v - %s",
				pricemode.Distributor, product.ID, product.Name))
			return
		}

	case pricemode.Export:
		if product.PriceExport != nil {
			invoice.Products[i].Price = *product.PriceExport
		} else {
			err = corerr.Tick(errTmp, "E7767315", fmt.Sprintf("price_mode %v not defined for product: %v - %s",
				pricemode.Export, product.ID, product.Name))
			return
		}

	default:
		invoice.Products[i].Price = product.PriceRetail
	}

	return

}

// TxInvoiceGenerator is used for create a unique number for each invoice based on year and locationID
func (p *BilInvoiceServ) TxInvoiceGenerator(db *gorm.DB, invoice bilmodel.Invoice) (yearCounter, cumulative uint64, invoiceNumber string, err error) {

	var preInvoice bilmodel.Invoice

	if preInvoice, err = p.Repo.TxLastInvoice(db, invoice.CompanyID, invoice.StoreID); err != nil {
		err = corerr.Tick(err, "E7798014", "error in fetching the last invoice", invoice.CompanyID, invoice.StoreID)
		return
	}

	yearCounter = preInvoice.YearCounter + 1
	cumulative = preInvoice.YearCumulative + 1

	invoiceNumber = p.Engine.Setting[bill.InvoiceNumberPattern].Value

	invoiceNumber = strings.Replace(
		invoiceNumber,
		consts.InvoicePatternYearCounter,
		fmt.Sprintf("%d", yearCounter), -1)

	invoiceNumber = strings.Replace(
		invoiceNumber,
		consts.InvoicePatternStoreCode,
		fmt.Sprintf("%v", invoice.Store.Code), -1)

	invoiceNumber = strings.Replace(
		invoiceNumber,
		consts.InvoicePatternYear,
		fmt.Sprintf("%d", invoice.Year), -1)

	return
}

// TxSave save an invoice, if it is exist update it, if not create it
func (p *BilInvoiceServ) TxSave(db *gorm.DB, invoice bilmodel.Invoice) (savedInvoice bilmodel.Invoice, err error) {
	if err = invoice.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7727507", corerr.ValidationFailed, invoice)
		return
	}

	if savedInvoice, err = p.Repo.TxSave(db, invoice); err != nil {
		err = corerr.Tick(err, "E7739556", "invoice not saved")
		return
	}

	return
}

// Delete invoice, it is soft delete
func (p *BilInvoiceServ) Delete(fix types.FixedCol) (invoice bilmodel.Invoice, err error) {
	if invoice, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7773432", "invoice not found for deleting")
		return
	}

	if err = p.Repo.Delete(invoice); err != nil {
		err = corerr.Tick(err, "E7729491", "invoice not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *BilInvoiceServ) Excel(params param.Param) (invoices []bilmodel.Invoice, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", bilmodel.InvoiceTable)

	if invoices, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E7764621", "cant generate the excel list for invoices")
		return
	}

	return
}
