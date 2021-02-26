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

// MatTagServ for injecting auth matrepo
type MatTagServ struct {
	Repo   matrepo.TagRepo
	Engine *core.Engine
}

// ProvideMatTagService for tag is used in wire
func ProvideMatTagService(p matrepo.TagRepo) MatTagServ {
	return MatTagServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting tag by it's id
func (p *MatTagServ) FindByID(fix types.FixedCol) (tag matmodel.Tag, err error) {
	if tag, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7194522", "can't fetch the tag", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	return
}

// List of tags, it support pagination and search and return back count
func (p *MatTagServ) List(params param.Param) (tags []matmodel.Tag,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" mat_tags.company_id = '%v' ", params.CompanyID)
	}

	if tags, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in tags list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in tags count")
	}

	return
}

// Create a tag
func (p *MatTagServ) Create(tag matmodel.Tag) (createdTag matmodel.Tag, err error) {

	if err = tag.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7148742", "validation failed in creating the tag", tag)
		return
	}

	if createdTag, err = p.Repo.Create(tag); err != nil {
		err = corerr.Tick(err, "E7130682", "tag not created", tag)
		return
	}

	return
}

// Save a tag, if it is exist update it, if not create it
func (p *MatTagServ) Save(tag matmodel.Tag) (savedTag matmodel.Tag, err error) {
	if err = tag.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7132797", corerr.ValidationFailed, tag)
		return
	}

	if savedTag, err = p.Repo.Save(tag); err != nil {
		err = corerr.Tick(err, "E7177049", "tag not saved")
		return
	}

	return
}

// Delete tag, it is soft delete
func (p *MatTagServ) Delete(fix types.FixedCol) (tag matmodel.Tag, err error) {
	if tag, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7165277", "tag not found for deleting")
		return
	}

	if err = p.Repo.Delete(tag); err != nil {
		err = corerr.Tick(err, "E7166026", "tag not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *MatTagServ) Excel(params param.Param) (tags []matmodel.Tag, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", matmodel.TagTable)

	if tags, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E7125861", "cant generate the excel list for tags")
		return
	}

	return
}
