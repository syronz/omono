package basapi

import (
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/message/basterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// RoleAPI for injecting role service
type RoleAPI struct {
	Service service.BasRoleServ
	Engine  *core.Engine
}

// ProvideRoleAPI for role is used in wire
func ProvideRoleAPI(c service.BasRoleServ) RoleAPI {
	return RoleAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a role by it's id
func (p *RoleAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var role basmodel.Role
	var id uint

	if id, err = resp.GetID(c.Param("roleID"), "E1053982", basterm.Role); err != nil {
		return
	}

	if role, err = p.Service.FindByID(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewRole)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Role).
		JSON(role)
}

// List of roles
func (p *RoleAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basmodel.RoleTable, base.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1097829"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListRole)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Roles).
		JSON(data)
}

// Create role
func (p *RoleAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var role, createdRole basmodel.Role
	var err error

	if err = resp.Bind(&role, "E1088259", base.Domain, basterm.Role); err != nil {
		return
	}

	if createdRole, err = p.Service.Create(role); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(base.CreateRole, role)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.Role).
		JSON(createdRole)
}

// Update role
func (p *RoleAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error

	var role, roleBefore, roleUpdated basmodel.Role
	var id uint

	if id, err = resp.GetID(c.Param("roleID"), "E1082097", basterm.Role); err != nil {
		return
	}

	if err = resp.Bind(&role, "E1076117", base.Domain, basterm.Role); err != nil {
		return
	}

	if roleBefore, err = p.Service.FindByID(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	role.ID = id
	role.CreatedAt = roleBefore.CreatedAt
	if roleUpdated, err = p.Service.Save(role); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.UpdateRole, roleBefore, role)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, basterm.Role).
		JSON(roleUpdated)
}

// Delete role
func (p *RoleAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var role basmodel.Role
	var id uint

	if id, err = resp.GetID(c.Param("roleID"), "E1088446", basterm.Role); err != nil {
		return
	}

	if role, err = p.Service.Delete(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.DeleteRole, role)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.Role).
		JSON()
}

// Excel generate excel files based on search
func (p *RoleAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Roles, base.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1013408"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	roles, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("role")
	ex.AddSheet("Roles").
		AddSheet("Summary").
		Active("Roles").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "D", 15.3).
		SetColWidth("E", "E", 80).
		SetColWidth("F", "G", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Roles").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Resources", "Description", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Resources", "Description", "UpdatedAt").
		WriteData(roles).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ExcelRole)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
