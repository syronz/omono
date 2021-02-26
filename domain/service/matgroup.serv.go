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

// MatGroupServ for injecting auth matrepo
type MatGroupServ struct {
	Repo   matrepo.GroupRepo
	Engine *core.Engine
}

// ProvideMatGroupService for group is used in wire
func ProvideMatGroupService(p matrepo.GroupRepo) MatGroupServ {
	return MatGroupServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting group by it's id
func (p *MatGroupServ) FindByID(fix types.FixedCol) (group matmodel.Group, err error) {
	if group, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7182452", "can't fetch the group", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	if group.Products, err = p.GetGroupProducts(group.ID); err != nil {
		err = corerr.Tick(err, "E7136676", "can't fetch the group's prodcut inside find by id", group)
		return
	}

	return
}

// GetGroupProducts is used for returning the products for a group
func (p *MatGroupServ) GetGroupProducts(groupID types.RowID) (pProducts []matmodel.Product, err error) {
	if pProducts, err = p.Repo.GetGroupProducts(groupID); err != nil {
		err = corerr.Tick(err, "E7151154", "can't fetch the group's products", groupID)
		return
	}

	return
}

// List of groups, it support pagination and search and return back count
func (p *MatGroupServ) List(params param.Param) (groups []matmodel.Group,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" mat_groups.company_id = '%v' ", params.CompanyID)
	}

	if groups, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in groups list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in groups count")
	}

	return
}

// Create a group
func (p *MatGroupServ) Create(group matmodel.Group) (createdGroup matmodel.Group, err error) {

	if err = group.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7165050", "validation failed in creating the group", group)
		return
	}

	if createdGroup, err = p.Repo.Create(group); err != nil {
		err = corerr.Tick(err, "E7178030", "group not created", group)
		return
	}

	return
}

// AddProduct is used for connect product with a product
func (p *MatGroupServ) AddProduct(gProduct matmodel.GroupProduct) (createdProduct matmodel.GroupProduct, err error) {
	if createdProduct, err = p.Repo.AddProduct(gProduct); err != nil {
		err = corerr.Tick(err, "E7183147", "product not added to the group", gProduct)
		return
	}

	return
}

// DelProduct Delete a product via its id
func (p *MatGroupServ) DelProduct(id types.RowID) (err error) {
	if err = p.Repo.DelProduct(id); err != nil {
		err = corerr.Tick(err, "E7145638", "product not removed from group", id)
		return
	}

	return
}

// Save a group, if it is exist update it, if not create it
func (p *MatGroupServ) Save(group matmodel.Group) (savedGroup matmodel.Group, err error) {
	if err = group.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E7162676", corerr.ValidationFailed, group)
		return
	}

	if savedGroup, err = p.Repo.Save(group); err != nil {
		err = corerr.Tick(err, "E7190562", "group not saved")
		return
	}

	return
}

// Delete group, it is soft delete
func (p *MatGroupServ) Delete(fix types.FixedCol) (group matmodel.Group, err error) {
	if group, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E7128465", "group not found for deleting")
		return
	}

	if err = p.Repo.Delete(group); err != nil {
		err = corerr.Tick(err, "E7126065", "group not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *MatGroupServ) Excel(params param.Param) (groups []matmodel.Group, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", matmodel.GroupTable)

	if groups, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E7119907", "cant generate the excel list for groups")
		return
	}

	return
}
