package service

import (
	"fmt"
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/base/message/basterm"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"

	"gorm.io/gorm"
)

// BasRoleServ for injecting auth basrepo
type BasRoleServ struct {
	Repo   basrepo.RoleRepo
	Engine *core.Engine
}

// ProvideBasRoleService for role is used in wire
func ProvideBasRoleService(p basrepo.RoleRepo) BasRoleServ {
	return BasRoleServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting role by it's id
func (p *BasRoleServ) FindByID(fix types.FixedCol) (role basmodel.Role, err error) {
	if role, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1043183", "can't fetch the role", fix.CompanyID, fix.NodeID, fix.ID)
		return
	}

	return
}

// List of roles, it support pagination and search and return back count
func (p *BasRoleServ) List(params param.Param) (roles []basmodel.Role,
	count int64, err error) {

	if params.CompanyID != 0 {
		params.PreCondition = fmt.Sprintf(" bas_roles.company_id = '%v' ", params.CompanyID)
	}

	if roles, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in roles list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in roles count")
	}

	return
}

// Create a role
func (p *BasRoleServ) Create(role basmodel.Role) (createdRole basmodel.Role, err error) {
	return p.TxCreate(p.Engine.DB, role)
}

// TxCreate is used in case of transaction activated
func (p *BasRoleServ) TxCreate(db *gorm.DB, role basmodel.Role) (createdRole basmodel.Role, err error) {
	if err = role.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1098554", "validation failed in creating the role", role)
		return
	}

	if createdRole, err = p.Repo.TxCreate(db, role); err != nil {
		err = corerr.Tick(err, "E1042894", "role not created", role)
		return
	}

	return
}

// Save a role, if it is exist update it, if not create it
func (p *BasRoleServ) Save(role basmodel.Role) (savedRole basmodel.Role, err error) {
	if err = role.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1037119", corerr.ValidationFailed, role)
		return
	}

	if savedRole, err = p.Repo.Save(role); err != nil {
		err = corerr.Tick(err, "E1078742", "role not saved")
		return
	}

	BasAccessResetFullCache()
	return
}

// Delete role, it is soft delete
func (p *BasRoleServ) Delete(fix types.FixedCol) (role basmodel.Role, err error) {
	if role, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E1052861", "role not found for deleting")
		return
	}

	// check if a user exist with this role
	params := param.NewForDelete("bas_users", "role_id", fix.ID)
	userServ := ProvideBasUserService(basrepo.ProvideUserRepo(p.Engine))

	var users []basmodel.User
	if users, err = userServ.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1035950", "related user not fetch for delete an account")
		return
	}

	if len(users) > 0 {
		err = limberr.New("roles is related to a user", "E1047255").
			Message(corerr.VIsConnectedToAVVPleaseDeleteItFirst, dict.R(basterm.Role),
				dict.R(basterm.User), users[0].Username).
			Custom(corerr.ForeignErr).Build()
		return
	}

	if err = p.Repo.Delete(role); err != nil {
		err = corerr.Tick(err, "E1017987", "role not deleted")
		return
	}

	BasAccessResetFullCache()
	return
}

// Excel is used for export excel file
func (p *BasRoleServ) Excel(params param.Param) (roles []basmodel.Role, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", basmodel.RoleTable)

	if roles, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1067385", "cant generate the excel list for roles")
		return
	}

	return
}
