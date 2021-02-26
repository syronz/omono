package service

import (
	"fmt"
	"omono/domain/material/matmodel"
	"omono/domain/material/matrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
)

// MatProductServ for injecting auth matrepo
type MatProductServ struct {
	Repo   matrepo.ProductRepo
	Engine *core.Engine
}

// ProvideMatProductService for product is used in wire
func ProvideMatProductService(p matrepo.ProductRepo) MatProductServ {
	return MatProductServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting product by it's id
func (p *MatProductServ) FindByID(fix types.FixedCol) (product matmodel.Product, err error) {
	if product, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7143172", "can't fetch the product", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	if product.Tags, err = p.GetProductTags(product.ID); err != nil {
		err = corerr.Tick(err, "E7131307", "can't fetch the product's tag inside find by id", product)
		return
	}

	return
}

// FindByRowID for getting product by it's id
func (p *MatProductServ) FindByRowID(rowID types.RowID) (product matmodel.Product, err error) {
	if product, err = p.Repo.FindByRowID(rowID); err != nil {
		err = corerr.Tick(err, "E7156405", "can't fetch the product", rowID)
		return
	}

	if product.Tags, err = p.GetProductTags(product.ID); err != nil {
		err = corerr.Tick(err, "E7156405", "can't fetch the product's tag inside find by id", product)
		return
	}

	return
}

// GetProductTags is used for returning the tags for a product
func (p *MatProductServ) GetProductTags(productID types.RowID) (pTags []matmodel.Tag, err error) {
	if pTags, err = p.Repo.GetProductTags(productID); err != nil {
		err = corerr.Tick(err, "E7179436", "can't fetch the product's tags", productID)
		return
	}

	return
}

// List of products, it support pagination and search and return back count
func (p *MatProductServ) List(params param.Param) (products []matmodel.Product,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" mat_products.company_id = '%v' ", params.CompanyID)
	}

	if products, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in products list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in products count")
	}

	return
}

// Create a product
func (p *MatProductServ) Create(product matmodel.Product) (createdProduct matmodel.Product, err error) {

	if err = product.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7115032", "validation failed in creating the product", product)
		return
	}

	if createdProduct, err = p.Repo.Create(product); err != nil {
		err = corerr.Tick(err, "E7158866", "product not created", product)
		return
	}

	for _, v := range product.Tags {
		pt := matmodel.ProductTag{
			ProductID: createdProduct.ID,
			TagID:     v.ID,
		}
		pt.CompanyID = createdProduct.CompanyID
		pt.NodeID = createdProduct.NodeID

		if _, err := p.AddTag(pt); err != nil {
			err = corerr.Tick(err, "E7142167", "tag not added to product", product)
		}
	}

	return
}

// AddTag is used for connect product with a tag
func (p *MatProductServ) AddTag(productTag matmodel.ProductTag) (createdTag matmodel.ProductTag, err error) {
	if createdTag, err = p.Repo.AddTag(productTag); err != nil {
		err = corerr.Tick(err, "E7194777", "tag not added to the product", productTag)
		return
	}

	return
}

// DelTag Delete a tag via its id
func (p *MatProductServ) DelTag(id types.RowID) (err error) {
	if err = p.Repo.DelTag(id); err != nil {
		err = corerr.Tick(err, "E7171452", "tag not removed from product", id)
		return
	}

	return
}

// Save a product, if it is exist update it, if not create it
func (p *MatProductServ) Save(product matmodel.Product) (savedProduct matmodel.Product, err error) {
	if err = product.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7139271", corerr.ValidationFailed, product)
		return
	}

	if savedProduct, err = p.Repo.Save(product); err != nil {
		err = corerr.Tick(err, "E7178049", "product not saved")
		return
	}

	return
}

// Delete product, it is soft delete
func (p *MatProductServ) Delete(fix types.FixedCol) (product matmodel.Product, err error) {
	if product, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7136806", "product not found for deleting")
		return
	}

	if err = p.Repo.Delete(product); err != nil {
		err = corerr.Tick(err, "E7113691", "product not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *MatProductServ) Excel(params param.Param) (products []matmodel.Product, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", matmodel.ProductTable)

	if products, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E7155309", "cant generate the excel list for products")
		return
	}

	return
}
