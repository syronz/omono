package eacrepo

import (
	"github.com/syronz/limberr"
	"omono/domain/base/message/basterm"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
	"reflect"
)

// RateRepo for injecting engine
type RateRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideRateRepo is used in wire and initiate the Cols
func ProvideRateRepo(engine *core.Engine) RateRepo {
	return RateRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(eacmodel.Rate{}), eacmodel.RateTable),
	}
}

// FindByID finds the rate via its id
func (p *RateRepo) FindByID(fix types.FixedCol) (rate eacmodel.Rate, err error) {
	err = p.Engine.ReadDB.Table(eacmodel.RateTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&rate).Error

	rate.ID = fix.ID
	err = p.dbError(err, "E1459246", rate, corterm.List)

	return
}

// List returns an array of rates
func (p *RateRepo) List(params param.Param) (rates []eacmodel.Rate, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1464677").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1467392").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(eacmodel.RateTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&rates).Error

	err = p.dbError(err, "E1446901", eacmodel.Rate{}, corterm.List)

	return
}

// LastRates return back the realtime value for each currency rate
func (p *RateRepo) LastRates() (rates []eacmodel.Rate, err error) {
	err = p.Engine.ReadDB.Raw(`WITH eac_rates AS (
	SELECT m.*, ROW_NUMBER() OVER (PARTITION BY currency_id,city_id ORDER BY id DESC) AS rn
	FROM eac_rates AS m)
	select r.*,c.name,c.code from eac_rates r 
	inner join eac_currencies c on c.id = r.currency_id 
	where rn = 1;`).Scan(&rates).Error

	err = p.dbError(err, "E1446901", eacmodel.Rate{}, corterm.List)

	return
}

// Count of rates, mainly calls with List
func (p *RateRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1497655").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(eacmodel.RateTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1464003", eacmodel.Rate{}, corterm.List)
	return
}

// Save the rate, in case it is not exist create it
func (p *RateRepo) Save(rate eacmodel.Rate) (u eacmodel.Rate, err error) {
	if err = p.Engine.DB.Table(eacmodel.RateTable).Save(&rate).Error; err != nil {
		err = p.dbError(err, "E1428244", rate, corterm.Updated)
	}

	p.Engine.DB.Table(eacmodel.RateTable).Where("id = ?", rate.ID).Find(&u)
	return
}

// Create a rate
func (p *RateRepo) Create(rate eacmodel.Rate) (u eacmodel.Rate, err error) {
	if err = p.Engine.DB.Table(eacmodel.RateTable).Create(&rate).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1436065", rate, corterm.Created)
	}
	return
}

// Delete the rate
func (p *RateRepo) Delete(rate eacmodel.Rate) (err error) {
	if err = p.Engine.DB.Table(eacmodel.RateTable).Delete(&rate).Error; err != nil {
		err = p.dbError(err, "E1412523", rate, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *RateRepo) dbError(err error, code string, rate eacmodel.Rate, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, rate.ID, eacterm.Rates)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(eacterm.Rate), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
