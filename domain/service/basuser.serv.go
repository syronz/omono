package service

import (
	"fmt"
	"github.com/syronz/limberr"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/base/enum/accountstatus"
	"omono/domain/base/enum/accounttype"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
	"omono/pkg/password"
)

// BasUserServ for injecting auth basrepo
type BasUserServ struct {
	Repo   basrepo.UserRepo
	Engine *core.Engine
}

// ProvideBasUserService for user is used in wire
func ProvideBasUserService(p basrepo.UserRepo) BasUserServ {
	return BasUserServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting user by it's id
func (p *BasUserServ) FindByID(fix types.FixedCol) (user basmodel.User, err error) {
	if user, err = p.Repo.FindByID(fix); err != nil {
		// err = corerr.Tick(err, "E1066324", "can't fetch the user", fix.CompanyID, fix.NodeID, fix.ID)
		err = limberr.AddCode(err, "E1066324")
		return
	}

	return
}

// FindByUsername find user with username, used for auth
func (p *BasUserServ) FindByUsername(username string) (user basmodel.User, err error) {
	if user, err = p.Repo.FindByUsername(username); err != nil {
		err = corerr.Tick(err, "E1088844", "can't fetch the user by username", username)
		return
	}

	return
}

// List of users, it support pagination and search and return back count
func (p *BasUserServ) List(params param.Param) (users []basmodel.User,
	count int64, err error) {

	if users, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in users list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in users count")
	}

	return
}

// Create a user
func (p *BasUserServ) Create(user basmodel.User) (createdUser basmodel.User, err error) {

	if err = user.Validate(coract.Create); err != nil {
		err = corerr.TickValidate(err, "E1043810", "validatation failed in creating user", user)
		return
	}

	roleServ := ProvideBasRoleService(basrepo.ProvideRoleRepo(p.Engine))
	fix := types.FixedCol{
		ID:        user.RoleID,
		CompanyID: user.CompanyID,
		NodeID:    user.NodeID,
	}
	if _, err = roleServ.FindByID(fix); err != nil {
		err = corerr.TickValidate(err, "E1093935", "role_id is out of scope", user)
		return
	}

	db := p.Engine.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"users table"), "rollback recover create user")
			db.Rollback()
		}
	}()

	phoneServ := ProvideBasPhoneService(basrepo.ProvidePhoneRepo(p.Engine))
	accountServ := ProvideBasAccountService(basrepo.ProvideAccountRepo(p.Engine), phoneServ)

	account := basmodel.Account{
		NameEn:   &user.Name,
		NameKu:   &user.Name,
		NameAr:   &user.Name,
		Code:     user.Code,
		ParentID: types.RowIDPointer(p.Engine.Envs.ToUint64("DEFAULT_USER_PARENT_ID")),
		Type:     accounttype.User,
		Status:   accountstatus.Active,
	}
	account.CompanyID = user.CompanyID
	account.NodeID = user.NodeID
	account.Phones = user.Phones

	var createdAccount basmodel.Account
	if createdAccount, err = accountServ.TxCreate(db, account); err != nil {
		err = corerr.Tick(err, "E1032795", "error in creating account for user", user)

		db.Rollback()
		return
	}

	user.ID = createdAccount.ID
	user.Password, err = password.Hash(user.Password, p.Engine.Envs[base.PasswordSalt])
	glog.CheckError(err, fmt.Sprintf("Hashing password failed for %+v", user))

	if createdUser, err = p.Repo.TxCreate(db, user); err != nil {
		err = corerr.Tick(err, "E1064180", "error in creating user", user)

		db.Rollback()
		return
	}

	db.Commit()
	createdUser.Password = ""

	return

}

// Save user
func (p *BasUserServ) Save(user basmodel.User) (createdUser basmodel.User, err error) {
	if err = user.Validate(coract.Update); err != nil {
		err = corerr.TickValidate(err, "E1098252", corerr.ValidationFailed, user)
		return
	}

	roleServ := ProvideBasRoleService(basrepo.ProvideRoleRepo(p.Engine))
	fix := types.FixedCol{
		ID:        user.RoleID,
		CompanyID: user.CompanyID,
		NodeID:    user.NodeID,
	}
	if _, err = roleServ.FindByID(fix); err != nil {
		err = corerr.TickValidate(err, "E1072771", "role_id is out of scope", user)
		return
	}

	var oldUser basmodel.User
	fix = types.FixedCol{
		ID:        user.ID,
		CompanyID: user.CompanyID,
		NodeID:    user.NodeID,
	}
	oldUser, _ = p.FindByID(fix)

	// clonedEngine := p.Engine.Clone()
	// clonedEngine.DB = clonedEngine.DB.Begin()

	db := p.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"users table"), "rollback recover save user")
			db.Rollback()
		}
	}()

	if user.Password != "" {
		if user.Password, err = password.Hash(user.Password, p.Engine.Envs[base.PasswordSalt]); err != nil {
			err = corerr.Tick(err, "E1057832", "error in saving user", user)
		}
	} else {
		user.Password = oldUser.Password
	}

	userRepo := basrepo.ProvideUserRepo(p.Engine)
	phoneServ := ProvideBasPhoneService(basrepo.ProvidePhoneRepo(p.Engine))
	accountServ := ProvideBasAccountService(basrepo.ProvideAccountRepo(p.Engine), phoneServ)
	account := basmodel.Account{
		NameEn:   &user.Name,
		Code:     user.Code,
		ParentID: types.RowIDPointer(p.Engine.Envs.ToUint64("DEFAULT_USER_PARENT_ID")),
		Type:     accounttype.User,
		Status:   user.Status,
	}
	account.ID = user.ID
	account.CompanyID = user.CompanyID
	account.NodeID = user.NodeID

	if _, err = accountServ.Save(account); err != nil {
		err = corerr.Tick(err, "E1098648", "error in saving account inside the user", user)

		db.Rollback()
		return
	}

	if createdUser, err = userRepo.TxSave(db, user); err != nil {
		err = corerr.Tick(err, "E1062983", "error in saving user", user)

		db.Rollback()
		return
	}

	BasAccessDeleteFromCache(user.ID)

	db.Commit()
	createdUser.Password = ""

	return
}

// Delete user, it is hard delete, by deleting account related to the user
func (p *BasUserServ) Delete(fix types.FixedCol) (user basmodel.User, err error) {
	if user, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1031839", "user not found for deleting")
		return
	}

	glog.Debug(user)

	if err = p.Repo.Delete(user); err != nil {
		err = corerr.Tick(err, "E1088344", "user not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *BasUserServ) Excel(params param.Param) (users []basmodel.User, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", basmodel.UserTable)

	if users, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1064328", "cant generate the excel list")
	}

	return
}
