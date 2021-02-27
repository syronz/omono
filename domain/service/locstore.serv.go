package service

import (
	"fmt"
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/base/basmodel"
	"omono/domain/location"
	"omono/domain/location/locmodel"
	"omono/domain/location/locrepo"
	"omono/domain/location/locterm"
	"omono/internal/consts"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
)

// LocStoreServ for injecting auth locrepo
type LocStoreServ struct {
	Repo   locrepo.StoreRepo
	Engine *core.Engine
}

// ProvideLocStoreService for store is used in wire
func ProvideLocStoreService(p locrepo.StoreRepo) LocStoreServ {
	return LocStoreServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

var cacheStores map[types.RowID]locmodel.Store

func (p *LocStoreServ) loadStoresToCache() {
	params := param.New()
	params.Limit = consts.MaxRowsCount
	var err error
	var stores []locmodel.Store

	if stores, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in stores list")
		return
	}

	for i, v := range stores {
		if stores[i].Users, err = p.GetStoreUsers(v.ID); err != nil {
			err = corerr.Tick(err, "E2894526", "can't fetch the store's users inside list", v)
			// return
		}
	}

	cacheStores = make(map[types.RowID]locmodel.Store)

	for _, v := range stores {
		cacheStores[v.ID] = v
	}

}

// FindByID for getting store by it's id
func (p *LocStoreServ) FindByID(fix types.FixedCol) (store locmodel.Store, err error) {
	if v, ok := cacheStores[fix.ID]; ok {
		store = v
		return
	}

	if len(cacheStores) == 0 {
		p.loadStoresToCache()
	}

	if store, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E2824668", "can't fetch the store", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	if store.Users, err = p.GetStoreUsers(store.ID); err != nil {
		err = corerr.Tick(err, "E2834683", "can't fetch the store's prodcut inside find by id", store)
		return
	}

	return
}

// GetStoreUsers is used for returning the users for a store
func (p *LocStoreServ) GetStoreUsers(storeID types.RowID) (pUsers []basmodel.User, err error) {
	if pUsers, err = p.Repo.GetStoreUsers(storeID); err != nil {
		err = corerr.Tick(err, "E2872482", "can't fetch the store's users", storeID)
		return
	}

	return
}

// IsUserInStore check permission of user to location
func (p *LocStoreServ) IsUserInStore(storeFix types.FixedCol, userID types.RowID) (err error) {
	var store locmodel.Store

	if store, err = p.FindByID(storeFix); err != nil {
		return
	}

	for _, v := range store.Users {
		if v.ID == userID {
			return
		}
	}

	err = limberr.New("you don't have permission to this store", "E2815226").
		Message(corerr.YouDontHavePermissionToThisV, dict.R(locterm.Store)).
		Custom(corerr.ForbiddenErr).
		Domain(location.Domain).
		Build()

	return
}

// List of stores, it support pagination and search and return back count
func (p *LocStoreServ) List(params param.Param) (stores []locmodel.Store,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" loc_stores.company_id = '%v' ", params.CompanyID)
	}

	if stores, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in stores list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in stores count")
	}

	return
}

// Create a store
func (p *LocStoreServ) Create(store locmodel.Store) (createdStore locmodel.Store, err error) {

	if err = store.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E2827659", "validation failed in creating the store", store)
		return
	}

	if createdStore, err = p.Repo.Create(store); err != nil {
		err = corerr.Tick(err, "E2847932", "store not created", store)
		return
	}

	p.loadStoresToCache()

	return
}

// AddUser is used for connect user with a user
func (p *LocStoreServ) AddUser(gUser locmodel.StoreUser) (createdUser locmodel.StoreUser, err error) {
	if createdUser, err = p.Repo.AddUser(gUser); err != nil {
		err = corerr.Tick(err, "E2882170", "user not added to the store", gUser)
		return
	}
	p.loadStoresToCache()

	return
}

// DelUser Delete a user via its id
func (p *LocStoreServ) DelUser(id types.RowID) (err error) {
	if err = p.Repo.DelUser(id); err != nil {
		err = corerr.Tick(err, "E2859619", "user not removed from store", id)
		return
	}
	p.loadStoresToCache()

	return
}

// Save a store, if it is exist update it, if not create it
func (p *LocStoreServ) Save(store locmodel.Store) (savedStore locmodel.Store, err error) {
	if err = store.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E2893679", corerr.ValidationFailed, store)
		return
	}

	if savedStore, err = p.Repo.Save(store); err != nil {
		err = corerr.Tick(err, "E2826325", "store not saved")
		return
	}

	p.loadStoresToCache()

	return
}

// Delete store, it is soft delete
func (p *LocStoreServ) Delete(fix types.FixedCol) (store locmodel.Store, err error) {
	if store, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E2874187", "store not found for deleting")
		return
	}

	if err = p.Repo.Delete(store); err != nil {
		err = corerr.Tick(err, "E2831614", "store not deleted")
		return
	}

	p.loadStoresToCache()

	return
}

// Excel is used for export excel file
func (p *LocStoreServ) Excel(params param.Param) (stores []locmodel.Store, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", locmodel.StoreTable)

	if stores, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E2869900", "cant generate the excel list for stores")
		return
	}

	return
}
