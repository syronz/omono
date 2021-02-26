package matrepo

import (
	"github.com/syronz/limberr"
	"omono/domain/base/message/basterm"
	"omono/domain/material/matmodel"
	"omono/domain/material/matterm"
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
)

// ProductRepo for injecting engine
type ProductRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideProductRepo is used in wire and initiate the Cols
func ProvideProductRepo(engine *core.Engine) ProductRepo {
	return ProductRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(matmodel.Product{}), matmodel.ProductTable),
	}
}

// FindByID finds the product via its id
func (p *ProductRepo) FindByID(fix types.FixedCol) (product matmodel.Product, err error) {
	err = p.Engine.ReadDB.Table(matmodel.ProductTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&product).Error

	product.ID = fix.ID
	err = p.dbError(err, "E7162410", product, corterm.List)

	return
}

// FindByRowID finds the product via its id
func (p *ProductRepo) FindByRowID(rowID types.RowID) (product matmodel.Product, err error) {
	err = p.Engine.ReadDB.Table(matmodel.ProductTable).
		Where("id = ?", rowID.ToUint64()).
		First(&product).Error

	product.ID = rowID
	err = p.dbError(err, "E7150620", product, corterm.List)

	return
}

// GetProductTags returns all tags which is related to a product
func (p *ProductRepo) GetProductTags(productID types.RowID) (pTags []matmodel.Tag, err error) {
	err = p.Engine.ReadDB.Table(matmodel.ProductTagTable).
		Select("mat_product_tags.id as id, mat_tags.tag as tag").
		Joins("INNER JOIN mat_tags on mat_tags.id = mat_product_tags.tag_id").
		Where("mat_product_tags.product_id = ?", productID).
		Find(&pTags).Error

	err = p.dbError(err, "E7196529", matmodel.Product{}, corterm.List)

	return
}

// List returns an array of products
func (p *ProductRepo) List(params param.Param) (products []matmodel.Product, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E7137966").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhereDelete(p.Cols); err != nil {
		err = limberr.Take(err, "E7175511").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(matmodel.ProductTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&products).Error

	err = p.dbError(err, "E7184392", matmodel.Product{}, corterm.List)

	return
}

// Count of products, mainly calls with List
func (p *ProductRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E7144133").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(matmodel.ProductTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E7181192", matmodel.Product{}, corterm.List)
	return
}

// Save the product, in case it is not exist create it
func (p *ProductRepo) Save(product matmodel.Product) (u matmodel.Product, err error) {
	if err = p.Engine.DB.Table(matmodel.ProductTable).Save(&product).Error; err != nil {
		err = p.dbError(err, "E7167023", product, corterm.Updated)
	}

	p.Engine.DB.Table(matmodel.ProductTable).Where("id = ?", product.ID).Find(&u)
	return
}

// Create a product
func (p *ProductRepo) Create(product matmodel.Product) (u matmodel.Product, err error) {
	if err = p.Engine.DB.Table(matmodel.ProductTable).Create(&product).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E7131707", product, corterm.Created)
	}
	return
}

// AddTag to add a product
func (p *ProductRepo) AddTag(productTag matmodel.ProductTag) (u matmodel.ProductTag, err error) {
	if err = p.Engine.DB.Table(matmodel.ProductTagTable).Create(&productTag).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E7146504", matmodel.Product{}, corterm.Created)
	}
	return
}

// DelTag Delete a tag via its id
func (p *ProductRepo) DelTag(id types.RowID) (err error) {
	if err = p.Engine.DB.Table(matmodel.ProductTagTable).
		Where("id = ?", id).
		Delete(&matmodel.ProductTag{}).Error; err != nil {
		err = p.dbError(err, "E7161644", matmodel.Product{}, corterm.Created)
	}
	return
}

// Delete the product
func (p *ProductRepo) Delete(product matmodel.Product) (err error) {
	now := time.Now()
	product.DeletedAt = &now
	if err = p.Engine.DB.Table(matmodel.ProductTable).Save(&product).Error; err != nil {
		err = p.dbError(err, "E7144646", product, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *ProductRepo) dbError(err error, code string, product matmodel.Product, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, product.ID, matterm.Products)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(matterm.Product), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(matterm.Product), product.Name).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, product.Name)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
