package locrepo

import (
	"omono/domain/base/basmodel"
	"omono/domain/location/locmodel"
	"omono/domain/location/locterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
	"github.com/syronz/limberr"
	"reflect"
)

// StoreRepo for injecting engine
type StoreRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideStoreRepo is used in wire and initiate the Cols
func ProvideStoreRepo(engine *core.Engine) StoreRepo {
	return StoreRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(locmodel.Store{}), locmodel.StoreTable),
	}
}

// FindByID finds the store via its id
func (p *StoreRepo) FindByID(fix types.FixedCol) (store locmodel.Store, err error) {
	err = p.Engine.ReadDB.Table(locmodel.StoreTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&store).Error

	store.ID = fix.ID
	err = p.dbError(err, "E2822501", store, corterm.List)

	return
}

// GetStoreUsers returns all users which is related to a store
func (p *StoreRepo) GetStoreUsers(storeID types.RowID) (users []basmodel.User, err error) {
	err = p.Engine.ReadDB.Table(locmodel.StoreUserTable).
		Select("*").
		Joins("INNER JOIN bas_users on bas_users.id = loc_store_users.user_id").
		Where("loc_store_users.store_id = ?", storeID).
		Find(&users).Error

	err = p.dbError(err, "E2897667", locmodel.Store{}, corterm.List)

	return
}

// List returns an array of stores
func (p *StoreRepo) List(params param.Param) (stores []locmodel.Store, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E2864603").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E2830694").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(locmodel.StoreTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&stores).Error

	err = p.dbError(err, "E2819108", locmodel.Store{}, corterm.List)

	return
}

// Count of stores, mainly calls with List
func (p *StoreRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E2853744").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(locmodel.StoreTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E2894157", locmodel.Store{}, corterm.List)
	return
}

// Save the store, in case it is not exist create it
func (p *StoreRepo) Save(store locmodel.Store) (u locmodel.Store, err error) {
	if err = p.Engine.DB.Table(locmodel.StoreTable).Save(&store).Error; err != nil {
		err = p.dbError(err, "E2854747", store, corterm.Updated)
	}

	p.Engine.DB.Table(locmodel.StoreTable).Where("id = ?", store.ID).Find(&u)
	return
}

// Create a store
func (p *StoreRepo) Create(store locmodel.Store) (u locmodel.Store, err error) {
	if err = p.Engine.DB.Table(locmodel.StoreTable).Create(&store).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E2863306", store, corterm.Created)
	}
	return
}

// AddUser to add a user
func (p *StoreRepo) AddUser(gUser locmodel.StoreUser) (u locmodel.StoreUser, err error) {
	if err = p.Engine.DB.Table(locmodel.StoreUserTable).Create(&gUser).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E2816635", locmodel.Store{}, corterm.Created)
	}
	return
}

// DelUser Delete a user via its id
func (p *StoreRepo) DelUser(id types.RowID) (err error) {
	if err = p.Engine.DB.Table(locmodel.StoreUserTable).
		Where("id = ?", id).
		Delete(&locmodel.StoreUser{}).Error; err != nil {
		err = p.dbError(err, "E2831435", locmodel.Store{}, corterm.Created)
	}
	return
}

// Delete the store
func (p *StoreRepo) Delete(store locmodel.Store) (err error) {
	if err = p.Engine.DB.Table(locmodel.StoreTable).Delete(&store).Error; err != nil {
		err = p.dbError(err, "E2817146", store, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *StoreRepo) dbError(err error, code string, store locmodel.Store, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, store.ID, locterm.Stores)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(locterm.Stores),
				dict.R(locterm.Store), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(locterm.Store), store.Name).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, store.Name)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
