package service

import (
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/eaccounting/eacterm"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
)

// EacRateServ for injecting auth eacrepo
type EacRateServ struct {
	Repo   eacrepo.RateRepo
	Engine *core.Engine
}

// ProvideEacRateService for rate is used in wire
func ProvideEacRateService(p eacrepo.RateRepo) EacRateServ {
	return EacRateServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// [cityID][currencyID]
var cacheRates map[types.RowID]map[types.RowID]eacmodel.Rate

func (p *EacRateServ) loadRatesToCache() {
	var err error
	var rates []eacmodel.Rate

	if rates, err = p.Repo.LastRates(); err != nil {
		glog.CheckError(err, "error in getting last rates")
		return
	}

	cacheRates = make(map[types.RowID]map[types.RowID]eacmodel.Rate)

	for _, v := range rates {
		if _, ok := cacheRates[v.CityID]; ok {
			cacheRates[v.CityID][v.CurrencyID] = v
		} else {
			cacheRates[v.CityID] = make(map[types.RowID]eacmodel.Rate)
			cacheRates[v.CityID][v.CurrencyID] = v
		}
	}
}

// RatesInCity returns last rate with currencies
func (p *EacRateServ) RatesInCity(cityID types.RowID) (rates []eacmodel.Rate, err error) {
	if len(cacheRates) == 0 {
		p.loadRatesToCache()
	}

	if mapRates, ok := cacheRates[cityID]; ok {
		for _, v := range mapRates {
			rates = append(rates, v)
		}
		return
	}

	err = limberr.New("no rates for city").
		Message(eacterm.ThereIsNoCurrencyRate).
		Custom(corerr.PreDataInsertedErr).
		Build()

	return
}

// GetRate when try to return specific rate per city and currency
func (p *EacRateServ) GetRate(cityID, currencyID types.RowID) (rate eacmodel.Rate, err error) {
	var rates []eacmodel.Rate

	if rates, err = p.RatesInCity(cityID); err != nil {
		return
	}

	for _, v := range rates {
		if v.CurrencyID == currencyID {
			rate = v
			return
		}
	}

	err = limberr.New("no rates for currency", "E7788617").
		Message(eacterm.ThereIsNoCurrencyRate).
		Custom(corerr.PreDataInsertedErr).
		Domain("EAccounting").
		Build()

	return

}

// FindByID for getting rate by it's id
func (p *EacRateServ) FindByID(fix types.FixedCol) (rate eacmodel.Rate, err error) {
	if rate, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1444863", "can't fetch the rate", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	return
}

// List of rates, it support pagination and search and return back count
func (p *EacRateServ) List(params param.Param) (rates []eacmodel.Rate,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" eac_rates.company_id = '%v' ", params.CompanyID)
	}

	if rates, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in rates list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in rates count")
	}

	return
}

// Create a rate
func (p *EacRateServ) Create(rate eacmodel.Rate) (createdRate eacmodel.Rate, err error) {

	if err = rate.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1462989", "validation failed in creating the rate", rate)
		return
	}

	if createdRate, err = p.Repo.Create(rate); err != nil {
		err = corerr.Tick(err, "E1473917", "rate not created", rate)
		return
	}

	p.loadRatesToCache()

	return
}

// Save a rate, if it is exist update it, if not create it
func (p *EacRateServ) Save(rate eacmodel.Rate) (savedRate eacmodel.Rate, err error) {
	if err = rate.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1446259", corerr.ValidationFailed, rate)
		return
	}

	if savedRate, err = p.Repo.Save(rate); err != nil {
		err = corerr.Tick(err, "E1471901", "rate not saved")
		return
	}

	p.loadRatesToCache()

	return
}

// Delete rate, it is soft delete
func (p *EacRateServ) Delete(fix types.FixedCol) (rate eacmodel.Rate, err error) {
	if rate, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1422538", "rate not found for deleting")
		return
	}

	if err = p.Repo.Delete(rate); err != nil {
		err = corerr.Tick(err, "E1467783", "rate not deleted")
		return
	}

	p.loadRatesToCache()

	return
}

// Excel is used for export excel file
func (p *EacRateServ) Excel(params param.Param) (rates []eacmodel.Rate, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", eacmodel.RateTable)

	if rates, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1453533", "cant generate the excel list for rates")
		return
	}

	return
}
