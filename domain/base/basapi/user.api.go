package basapi

import (
	"fmt"
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/pkg/glog"
	"omono/pkg/helper/excel"

	"github.com/syronz/dict"

	"github.com/gin-gonic/gin"
)

// UserAPI for injecting user service
type UserAPI struct {
	Service service.BasUserServ
	Engine  *core.Engine
}

// ProvideUserAPI for user is used in wire
func ProvideUserAPI(c service.BasUserServ) UserAPI {
	return UserAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a user by it's id
func (p *UserAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var user basmodel.User
	var id uint

	if id, err = resp.GetID(c.Param("userID"), "E1090173", basterm.User); err != nil {
		return
	}

	if user, err = p.Service.FindByID(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	user.Password = ""

	resp.Record(base.ViewUser)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.User).
		JSON(user)
}

// FindByUsername is used when we try to find a user with username
func (p *UserAPI) FindByUsername(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	username := c.Param("username")
	var err error

	user, err := p.Service.FindByUsername(username)
	if err != nil {
		resp.Status(http.StatusNotFound).Error(err).JSON()
		return
	}

	user.Password = ""

	resp.Status(http.StatusOK).JSON(user)
}

// List of users
func (p *UserAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basmodel.UserTable, base.Domain)

	if username := c.Query("username"); username != "" {
		params.Filter = fmt.Sprintf("bas_users.username[eq]'%v'", username)
	}

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListUser)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Users).
		JSON(data)
}

// Create user
func (p *UserAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var user, createdUser basmodel.User
	var err error

	if err = resp.Bind(&user, "E1082301", base.Domain, basterm.User); err != nil {
		return
	}

	if createdUser, err = p.Service.Create(user); err != nil {
		resp.Error(err).JSON()
		return
	}

	user.Password = ""

	resp.RecordCreate(base.CreateUser, user)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, dict.R(basterm.User)).
		JSON(createdUser)
}

// Update user
func (p *UserAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error

	var user, userBefore, userUpdated basmodel.User
	var id uint

	if id, err = resp.GetID(c.Param("userID"), "E1097541", basterm.User); err != nil {
		return
	}

	if err = resp.Bind(&user, "E1065844", base.Domain, basterm.User); err != nil {
		return
	}

	user.ID = id
	user.CreatedAt = userBefore.CreatedAt
	if userUpdated, userBefore, err = p.Service.Save(user); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.UpdateUser, userBefore, userUpdated)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, dict.R(basterm.User)).
		JSON(userUpdated)
}

// Delete user
func (p *UserAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var user basmodel.User
	var id uint

	if id, err = resp.GetID(c.Param("userID"), "E1046157", basterm.User); err != nil {
		return
	}

	if user, err = p.Service.Delete(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.DeleteUser, user)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.User).
		JSON()
}

// Excel generate excel files based on search
func (p *UserAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Users, base.Domain)
	var err error

	users, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("node").
		AddSheet("Nodes").
		AddSheet("Summary").
		Active("Nodes").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("A", "C", 17).
		SetColWidth("D", "F", 15.3).
		SetColWidth("H", "H", 20).
		SetColWidth("N", "O", 20).
		Active("Summary").
		Active("Nodes").
		WriteHeader("ID", "Company ID", "Node ID", "Username", "Role", "Lang", "Email")

	for i, v := range users {
		column := &[]interface{}{
			v.ID,
			v.Username,
			v.Role,
			v.Lang,
			v.Email,
		}
		err = ex.File.SetSheetRow(ex.ActiveSheet, fmt.Sprint("A", i+2), column)
		glog.CheckError(err, "Error in writing to the excel in user")
	}

	ex.Sheets[ex.ActiveSheet].Row = len(users) + 1

	ex.AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ExcelUser)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
