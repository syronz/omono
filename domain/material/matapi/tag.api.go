package matapi

import (
	"net/http"
	"omono/domain/material"
	"omono/domain/material/matmodel"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// TagAPI for injecting tag service
type TagAPI struct {
	Service service.MatTagServ
	Engine  *core.Engine
}

// ProvideTagAPI for tag is used in wire
func ProvideTagAPI(c service.MatTagServ) TagAPI {
	return TagAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a tag by it's id
func (p *TagAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error
	var tag matmodel.Tag
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("tagID"), "E7188907", corterm.Tag); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if tag, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ViewTag)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, corterm.Tag).
		JSON(tag)
}

// List of tags
func (p *TagAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, matmodel.TagTable, material.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7172418"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ListTag)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, corterm.Tags).
		JSON(data)
}

// Create tag
func (p *TagAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var tag, createdTag matmodel.Tag
	var err error

	if tag.CompanyID, tag.NodeID, err = resp.GetCompanyNode("E7141177", material.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if tag.CompanyID, err = resp.GetCompanyID("E7114803"); err != nil {
		return
	}

	if !resp.CheckRange(tag.CompanyID) {
		return
	}

	if err = resp.Bind(&tag, "E7186853", material.Domain, corterm.Tag); err != nil {
		return
	}

	if createdTag, err = p.Service.Create(tag); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(material.CreateTag, tag)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, corterm.Tag).
		JSON(createdTag)
}

// Update tag
func (p *TagAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error

	var tag, tagBefore, tagUpdated matmodel.Tag
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("tagID"), "E7161213", corterm.Tag); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&tag, "E7198760", material.Domain, corterm.Tag); err != nil {
		return
	}

	if tagBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	tag.ID = fix.ID
	tag.CompanyID = fix.CompanyID
	tag.NodeID = fix.NodeID
	tag.CreatedAt = tagBefore.CreatedAt
	if tagUpdated, err = p.Service.Save(tag); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.UpdateTag, tagBefore, tag)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, corterm.Tag).
		JSON(tagUpdated)
}

// Delete tag
func (p *TagAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, material.Domain)
	var err error
	var tag matmodel.Tag
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("tagID"), "E7195765", corterm.Tag); err != nil {
		return
	}

	if tag, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.DeleteTag, tag)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, corterm.Tag).
		JSON()
}

// Excel generate excel files eaced on search
func (p *TagAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, corterm.Tags, material.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E7152019"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	tags, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("tag")
	ex.AddSheet("Tags").
		AddSheet("Summary").
		Active("Tags").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "F", 15.3).
		SetColWidth("G", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Tags").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Symbol", "Code", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Symbol", "Code", "UpdatedAt").
		WriteData(tags).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(material.ExcelTag)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
