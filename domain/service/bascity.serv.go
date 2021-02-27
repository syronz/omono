package service

import (
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
)

// BasCityServ for injecting auth basrepo
type BasCityServ struct {
	Repo   basrepo.CityRepo
	Engine *core.Engine
}

// ProvideBasCityService for city is used in wire
func ProvideBasCityService(p basrepo.CityRepo) BasCityServ {
	return BasCityServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting city by it's id
func (p *BasCityServ) FindByID(fix types.FixedCol) (city basmodel.City, err error) {
	if city, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1032610", "can't fetch the city", fix.ID)
		return
	}

	return
}

// FindByCity for getting city by it's id
func (p *BasCityServ) FindByCity(cityNumber string) (city basmodel.City, err error) {
	if city, err = p.Repo.FindByCity(cityNumber); err != nil {
		// do not log error if it is not-found
		if limberr.GetCustom(err) != corerr.NotFoundErr {
			err = corerr.Tick(err, "E1042894", "can't fetch the city by city-number", cityNumber)
		}
		return
	}

	return
}

// List of cities, it support pagination and search and return back count
func (p *BasCityServ) List(params param.Param) (cities []basmodel.City,
	count int64, err error) {

	if cities, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in cities list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in cities count")
	}

	return
}

// Create a city
func (p *BasCityServ) Create(city basmodel.City) (createdCity basmodel.City, err error) {
	if err = city.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1418032", corerr.ValidationFailed, city)
		return
	}

	if createdCity, err = p.Repo.Save(city); err != nil {
		err = corerr.Tick(err, "E1415152", "city not saved")
		return
	}

	return
}

// Save a city, if it is exist update it, if not create it
func (p *BasCityServ) Save(city basmodel.City) (savedCity basmodel.City, err error) {
	glog.Debug(city)
	if err = city.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1023474", corerr.ValidationFailed, city)
		return
	}

	if savedCity, err = p.Repo.Save(city); err != nil {
		err = corerr.Tick(err, "E1044237", "city not saved")
		return
	}

	return
}

// Delete city, it is soft delete
func (p *BasCityServ) Delete(fix types.FixedCol) (city basmodel.City, err error) {
	if city, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1079342", "city not found for deleting")
		return
	}

	if err = p.Repo.Delete(city); err != nil {
		err = corerr.Tick(err, "E1092207", "city not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *BasCityServ) Excel(params param.Param) (cities []basmodel.City, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", basmodel.CityTable)

	if cities, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1019610", "cant generate the excel list for cities")
		return
	}

	return
}
