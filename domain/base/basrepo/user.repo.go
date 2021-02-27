package basrepo

import (
	// "github.com/cockroachdb/errors"

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

// UserRepo for injecting engine
type UserRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideUserRepo is used in wire and initiate the Cols
func ProvideUserRepo(engine *core.Engine) UserRepo {
	return UserRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(basmodel.User{}), basmodel.UserTable),
	}
}

// FindByID finds the user via its id
func (p *UserRepo) FindByID(fix types.FixedCol) (user basmodel.User, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, "*"); err != nil {
		err = limberr.Take(err, "E1084438").Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.UserTable).
		Select(colsStr).
		Joins("INNER JOIN bas_roles ON bas_roles.id = bas_users.role_id").
		Where("bas_users.id = ? AND bas_users.deleted_at IS NULL", fix.ID.ToUint64()).
		First(&user).Error

	user.ID = fix.ID
	err = p.dbError(err, "E1063251", user, corterm.List)

	return
}

// FindByUsername finds the user via its username
func (p *UserRepo) FindByUsername(username string) (user basmodel.User, err error) {
	err = p.Engine.ReadDB.Table(basmodel.UserTable).
		Select("bas_users.*, bas_roles.resources, bas_roles.name as role").
		Where("bas_users.username = ?", username).
		Joins("INNER JOIN bas_roles on bas_roles.id = bas_users.role_id").
		Scan(&user).Error

	user.Username = username
	err = p.dbError(err, "E1043108", user, corterm.List)

	return
}

// List returns an array of users
func (p *UserRepo) List(params param.Param) (users []basmodel.User, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1084438").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1043328").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.UserTable).Select(colsStr).
		Joins("INNER JOIN bas_roles ON bas_roles.id = bas_users.role_id").
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&users).Error

	err = p.dbError(err, "E1077340", basmodel.User{}, corterm.List)

	for i := range users {
		users[i].Password = ""
	}

	return
}

// Count of users, mainly calls with List
func (p *UserRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1042198").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.UserTable).
		Joins("INNER JOIN bas_roles ON bas_roles.id = bas_users.role_id").
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1042199", basmodel.User{}, corterm.List)
	return
}

// TxSave the user, in case it is not exist create it
func (p *UserRepo) TxSave(db *gorm.DB, user basmodel.User) (u basmodel.User, err error) {
	if err = db.Table(basmodel.UserTable).Save(&user).Error; err != nil {
		err = p.dbError(err, "E1056429", user, corterm.Updated)
	}

	db.Table(basmodel.UserTable).Where("id = ?", user.ID).Find(&u)
	return
}

// Create a user
func (p *UserRepo) Create(user basmodel.User) (u basmodel.User, err error) {
	if err = p.Engine.DB.Table(basmodel.UserTable).Create(&user).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1095328", user, corterm.Created)
	}
	return
}

// TxCreate a user
func (p *UserRepo) TxCreate(db *gorm.DB, user basmodel.User) (u basmodel.User, err error) {
	if err = db.Table(basmodel.UserTable).Create(&user).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1047736", user, corterm.Created)
	}
	return
}

// Delete the user
func (p *UserRepo) Delete(user basmodel.User) (err error) {
	now := time.Now()
	user.DeletedAt = &now
	if err = p.Engine.DB.Table(basmodel.UserTable).Save(&user).Error; err != nil {
		err = p.dbError(err, "E1044329", user, corterm.Deleted)
	}
	return
}

// dbError is an internal method for create proper database error
func (p *UserRepo) dbError(err error, code string, user basmodel.User, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, user.ID, basterm.Users)

	case corerr.ForeignErr:
		if action == corterm.Created {
			err = limberr.Take(err, code).
				Message(corerr.VNotExist, dict.R(basterm.Role)).
				Custom(corerr.ForeignErr).Build()
		} else {
			err = limberr.Take(err, code).
				Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(corterm.Fields),
					dict.R(basterm.User), dict.R(action)).
				Custom(corerr.ForeignErr).Build()
		}

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(basterm.User), user.Username).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "username", corerr.VisAlreadyExist, user.Username)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
