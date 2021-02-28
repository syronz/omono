package service

import (
	"fmt"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/pkg/glog"
	"omono/pkg/helper/password"

	"github.com/syronz/limberr"
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
func (p *BasUserServ) FindByID(id uint) (user basmodel.User, err error) {
	if user, err = p.Repo.FindByID(id); err != nil {
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

	db := p.Engine.DB.Begin()

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
func (p *BasUserServ) Save(user basmodel.User) (updatedUser, userBefore basmodel.User, err error) {
	if err = user.Validate(coract.Update); err != nil {
		err = corerr.TickValidate(err, "E1098252", corerr.ValidationFailed, user)
		return
	}

	userBefore, _ = p.FindByID(user.ID)

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
		user.Password = userBefore.Password
	}

	userRepo := basrepo.ProvideUserRepo(p.Engine)
	user.CreatedAt = userBefore.CreatedAt

	if updatedUser, err = userRepo.TxSave(db, user); err != nil {
		err = corerr.Tick(err, "E1062983", "error in saving user", user)

		db.Rollback()
		return
	}

	BasAccessDeleteFromCache(user.ID)

	db.Commit()
	updatedUser.Password = ""

	return
}

// Delete user, it is hard delete, by deleting account related to the user
func (p *BasUserServ) Delete(id uint) (user basmodel.User, err error) {
	if user, err = p.FindByID(id); err != nil {
		err = corerr.Tick(err, "E1031839", "user not found for deleting")
		return
	}

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
