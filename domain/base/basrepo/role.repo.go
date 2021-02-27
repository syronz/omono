package basrepo

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/base/basmodel"
	"omono/domain/base/message/basterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/helper"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// RoleRepo for injecting engine
type RoleRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideRoleRepo is used in wire and initiate the Cols
func ProvideRoleRepo(engine *core.Engine) RoleRepo {
	return RoleRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(basmodel.Role{}), basmodel.RoleTable),
	}
}

// FindByID finds the role via its id
func (p *RoleRepo) FindByID(fix types.FixedCol) (role basmodel.Role, err error) {
	err = p.Engine.ReadDB.Table(basmodel.RoleTable).
		Where("id = ?", fix.ID.ToUint64()).
		First(&role).Error

	role.ID = fix.ID
	err = p.dbError(err, "E1072991", role, corterm.List)

	return
}

// List returns an array of roles
func (p *RoleRepo) List(params param.Param) (roles []basmodel.Role, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1084438").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1032278").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.RoleTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&roles).Error

	err = p.dbError(err, "E1032861", basmodel.Role{}, corterm.List)

	return
}

// Count of roles, mainly calls with List
func (p *RoleRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1032288").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.RoleTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1039820", basmodel.Role{}, corterm.List)
	return
}

// Save the role, in case it is not exist create it
func (p *RoleRepo) Save(role basmodel.Role) (u basmodel.Role, err error) {
	if err = p.Engine.DB.Table(basmodel.RoleTable).Save(&role).Error; err != nil {
		err = p.dbError(err, "E1054817", role, corterm.Updated)
	}

	p.Engine.DB.Table(basmodel.RoleTable).Where("id = ?", role.ID).Find(&u)
	return
}

// TxCreate a role
func (p *RoleRepo) TxCreate(db *gorm.DB, role basmodel.Role) (u basmodel.Role, err error) {
	if err = db.Table(basmodel.RoleTable).Create(&role).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1053287", role, corterm.Created)
	}
	return
}

// Delete the role
func (p *RoleRepo) Delete(role basmodel.Role) (err error) {
	now := time.Now()
	role.DeletedAt = &now
	if err = p.Engine.DB.Table(basmodel.RoleTable).Save(&role).Error; err != nil {
		err = p.dbError(err, "E1067392", role, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper database error
func (p *RoleRepo) dbError(err error, code string, role basmodel.Role, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, role.ID, basterm.Roles)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(basterm.Role), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(basterm.Role), role.Name).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, role.Name)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
