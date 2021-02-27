package matrepo

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/material/matmodel"
	"omono/domain/material/matterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/helper"
	"reflect"
)

// GroupRepo for injecting engine
type GroupRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideGroupRepo is used in wire and initiate the Cols
func ProvideGroupRepo(engine *core.Engine) GroupRepo {
	return GroupRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(matmodel.Group{}), matmodel.GroupTable),
	}
}

// FindByID finds the group via its id
func (p *GroupRepo) FindByID(fix types.FixedCol) (group matmodel.Group, err error) {
	err = p.Engine.ReadDB.Table(matmodel.GroupTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&group).Error

	group.ID = fix.ID
	err = p.dbError(err, "E7151164", group, corterm.List)

	return
}

// GetGroupProducts returns all products which is related to a group
func (p *GroupRepo) GetGroupProducts(groupID types.RowID) (products []matmodel.Product, err error) {
	err = p.Engine.ReadDB.Table(matmodel.GroupProductTable).
		Select("*").
		Joins("INNER JOIN mat_products on mat_products.id = mat_group_products.product_id").
		Where("mat_group_products.group_id = ?", groupID).
		Find(&products).Error

	err = p.dbError(err, "E7182389", matmodel.Group{}, corterm.List)

	return
}

// List returns an array of groups
func (p *GroupRepo) List(params param.Param) (groups []matmodel.Group, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E7188930").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E7156968").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(matmodel.GroupTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&groups).Error

	err = p.dbError(err, "E7186094", matmodel.Group{}, corterm.List)

	return
}

// Count of groups, mainly calls with List
func (p *GroupRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E7149464").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(matmodel.GroupTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E7129881", matmodel.Group{}, corterm.List)
	return
}

// Save the group, in case it is not exist create it
func (p *GroupRepo) Save(group matmodel.Group) (u matmodel.Group, err error) {
	if err = p.Engine.DB.Table(matmodel.GroupTable).Save(&group).Error; err != nil {
		err = p.dbError(err, "E7154808", group, corterm.Updated)
	}

	p.Engine.DB.Table(matmodel.GroupTable).Where("id = ?", group.ID).Find(&u)
	return
}

// Create a group
func (p *GroupRepo) Create(group matmodel.Group) (u matmodel.Group, err error) {
	if err = p.Engine.DB.Table(matmodel.GroupTable).Create(&group).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E7141174", group, corterm.Created)
	}
	return
}

// AddProduct to add a product
func (p *GroupRepo) AddProduct(gProduct matmodel.GroupProduct) (u matmodel.GroupProduct, err error) {
	if err = p.Engine.DB.Table(matmodel.GroupProductTable).Create(&gProduct).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E7121382", matmodel.Group{}, corterm.Created)
	}
	return
}

// DelProduct Delete a product via its id
func (p *GroupRepo) DelProduct(id types.RowID) (err error) {
	if err = p.Engine.DB.Table(matmodel.GroupProductTable).
		Where("id = ?", id).
		Delete(&matmodel.GroupProduct{}).Error; err != nil {
		err = p.dbError(err, "E7183775", matmodel.Group{}, corterm.Created)
	}
	return
}

// Delete the group
func (p *GroupRepo) Delete(group matmodel.Group) (err error) {
	if err = p.Engine.DB.Table(matmodel.GroupTable).Delete(&group).Error; err != nil {
		err = p.dbError(err, "E7169627", group, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *GroupRepo) dbError(err error, code string, group matmodel.Group, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, group.ID, matterm.Groups)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(matterm.Groups),
				dict.R(matterm.Group), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(matterm.Group), group.Name).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, group.Name)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
