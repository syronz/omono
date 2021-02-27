package matapi

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"net/http"
	"omono/domain/material"
	"omono/domain/material/matmodel"
	"omono/domain/material/matterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// ProductAPI for injecting product service
type ProductAPI struct {
	Service service.MatProductServ
	Engine  *core.Engine
}

// ProvideProductAPI for product is used in wire
func ProvideProductAPI(c service.MatProductServ) ProductAPI {
	return ProductAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a product by it's id
func (p *ProductAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error
	var product matmodel.Product
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("productID"), "E7175382", matterm.Product); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if product, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ViewProduct)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, matterm.Product).
		JSON(product)
}

// List of products
func (p *ProductAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, matmodel.ProductTable, material.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7170086"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ListProduct)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, matterm.Products).
		JSON(data)
}

// Create product
func (p *ProductAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var product, createdProduct matmodel.Product
	var err error

	if product.CompanyID, product.NodeID, err = resp.GetCompanyNode("E7141917", material.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if product.CompanyID, err = resp.GetCompanyID("E7178594"); err != nil {
		return
	}

	if !resp.CheckRange(product.CompanyID) {
		return
	}

	if err = resp.Bind(&product, "E7148140", material.Domain, matterm.Product); err != nil {
		return
	}

	if createdProduct, err = p.Service.Create(product); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(material.CreateProduct, product)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, matterm.Product).
		JSON(createdProduct)
}

// AddTag create a tag for product
func (p *ProductAPI) AddTag(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error
	var productTag matmodel.ProductTag

	if productTag.CompanyID, productTag.NodeID, err =
		resp.GetCompanyNode("E7172484", material.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if productTag.ProductID, err = types.StrToRowID(c.Param("productID")); err != nil {
		resp.Error(err).JSON()
		return
	}

	if err = resp.Bind(&productTag, "E7131469", material.Domain, matterm.Product); err != nil {
		return
	}

	if _, err = p.Service.AddTag(productTag); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.CreateTag, productTag)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, corterm.Tag).
		JSON()
}

// DelTag Delete a tag for a product
func (p *ProductAPI) DelTag(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)

	var err error
	var ptID types.RowID
	if ptID, err = types.StrToRowID(c.Param("productTagID")); err != nil {
		err = limberr.Take(err, "E7119344").
			Message(corerr.InvalidVForV, dict.R(corterm.ID), "tag_id").
			Custom(corerr.ValidationFailedErr).Build()
		resp.Error(err).JSON()
		return
	}

	if err = p.Service.DelTag(ptID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.DeleteTag, ptID)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, corterm.Tag).
		JSON()
}

// Update product
func (p *ProductAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error

	var product, productBefore, productUpdated matmodel.Product
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("productID"), "E7117671", matterm.Product); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&product, "E7154238", material.Domain, matterm.Product); err != nil {
		return
	}

	if productBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	product.ID = fix.ID
	product.CompanyID = fix.CompanyID
	product.NodeID = fix.NodeID
	product.CreatedAt = productBefore.CreatedAt
	if productUpdated, err = p.Service.Save(product); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.UpdateProduct, productBefore, product)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, matterm.Product).
		JSON(productUpdated)
}

// Delete product
func (p *ProductAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error
	var product matmodel.Product
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("productID"), "E7196269", matterm.Product); err != nil {
		return
	}

	if product, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.DeleteProduct, product)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, matterm.Product).
		JSON()
}

// Excel generate excel files eaced on search
func (p *ProductAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, matterm.Products, material.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7155073"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	products, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	type ExProduct struct {
		matmodel.Product
		ExDescription string `json:"ex_description"`
	}

	exProducts := make([]ExProduct, len(products))

	for i, v := range products {
		exProducts[i].ID = v.ID
		exProducts[i].CompanyID = v.CompanyID
		exProducts[i].NodeID = v.NodeID
		exProducts[i].Name = v.Name
		if v.Description != nil {
			exProducts[i].ExDescription = *v.Description
		}
		exProducts[i].UpdatedAt = v.UpdatedAt
	}

	ex := excel.New("product")
	ex.AddSheet("Products").
		AddSheet("Summary").
		Active("Products").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "E", 15.3).
		SetColWidth("F", "F", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Products").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Description", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "ExDescription", "UpdatedAt").
		WriteData(exProducts).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ExcelProduct)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
