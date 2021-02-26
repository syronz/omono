package basrepo

import (
	"omono/domain/base/basmodel"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/pkg/helper"
	"github.com/syronz/limberr"
	"reflect"
)

// ActivityRepo for injecting engine
type ActivityRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideActivityRepo is used in wire
func ProvideActivityRepo(engine *core.Engine) ActivityRepo {
	return ActivityRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(basmodel.Activity{}), basmodel.ActivityTable),
	}
}

// Create ActivityRepo
func (p *ActivityRepo) Create(activity basmodel.Activity) (u basmodel.Activity, err error) {
	err = p.Engine.ActivityDB.
		Table(basmodel.ActivityTable).
		Create(&activity).Error
	return
}

// CreateBatch ActivityRepo
func (p *ActivityRepo) CreateBatch(activities []basmodel.Activity) (u basmodel.Activity, err error) {
	err = p.Engine.ActivityDB.
		Table(basmodel.ActivityTable).
		Create(&activities).Error
	return
}

// List of activities
func (p *ActivityRepo) List(params param.Param) (activities []basmodel.Activity, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E8964282").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E9367965").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ActivityDB.
		Table(basmodel.ActivityTable).
		Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&activities).Error

	return
}

// Count of activities
func (p *ActivityRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E9367965").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ActivityDB.
		Table(basmodel.ActivityTable).
		Select(params.Select).
		Where(whereStr).
		Count(&count).Error
	return
}
