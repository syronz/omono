package basrepo

import (
	"omono/domain/base/basmodel"
	"omono/domain/base/message/basterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/pkg/helper"
	"reflect"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"
)

// CityRepo for injecting engine
type CityRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideCityRepo is used in wire and initiate the Cols
func ProvideCityRepo(engine *core.Engine) CityRepo {
	return CityRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(basmodel.City{}), basmodel.CityTable),
	}
}

// FindByID finds the city via its id
func (p *CityRepo) FindByID(id uint) (city basmodel.City, err error) {
	err = p.Engine.ReadDB.Table(basmodel.CityTable).
		Where("id = ?", id).
		First(&city).Error

	city.ID = id
	err = p.dbError(err, "E1080299", city, corterm.List)

	return
}

// FindByCity finds the city via its id
func (p *CityRepo) FindByCity(cityNumber string) (city basmodel.City, err error) {
	err = p.Engine.ReadDB.Table(basmodel.CityTable).
		Where("city LIKE ?", cityNumber).First(&city).Error

	err = p.dbError(err, "E1073640", city, corterm.List)

	return
}

// List returns an array of cities
func (p *CityRepo) List(params param.Param) (cities []basmodel.City, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1091738").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1082911").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.CityTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&cities).Error

	err = p.dbError(err, "E1029474", basmodel.City{}, corterm.List)

	return
}

// Count of cities, mainly calls with List
func (p *CityRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1035100").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(basmodel.CityTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1067151", basmodel.City{}, corterm.List)
	return
}

// Save the city, in case it is not exist create it
func (p *CityRepo) Save(city basmodel.City) (u basmodel.City, err error) {
	if err = p.Engine.DB.Table(basmodel.CityTable).Save(&city).Error; err != nil {
		err = p.dbError(err, "E1020589", city, corterm.Updated)
	}

	p.Engine.DB.Table(basmodel.CityTable).Where("id = ?", city.ID).Find(&u)
	return
}

// Create a city
func (p *CityRepo) Create(city basmodel.City) (u basmodel.City, err error) {
	if err = p.Engine.DB.Table(basmodel.CityTable).Create(&city).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1037044", city, corterm.Created)
	}
	return
}

// Delete the city
func (p *CityRepo) Delete(city basmodel.City) (err error) {
	if err = p.Engine.DB.Unscoped().Table(basmodel.CityTable).Delete(&city).Error; err != nil {
		err = p.dbError(err, "E1026719", city, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper database error
func (p *CityRepo) dbError(err error, code string, city basmodel.City, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, city.ID, basterm.Cities)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(basterm.City), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(basterm.City), city.City).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "city", corerr.VisAlreadyExist, city.City)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
