package basapi

import (
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/param"
	"omono/internal/response"
	"omono/pkg/helper"
	"strconv"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"

	"github.com/gin-gonic/gin"
)

// AuthAPI for injecting auth service
type AuthAPI struct {
	Service service.BasAuthServ
	Engine  *core.Engine
}

// ProvideAuthAPI for auth used in wire
func ProvideAuthAPI(p service.BasAuthServ) AuthAPI {
	return AuthAPI{Service: p, Engine: p.Engine}
}

// Login auth
func (p *AuthAPI) Login(c *gin.Context) {
	var auth basmodel.Auth
	resp, params := response.NewParam(p.Engine, c, basterm.Users, base.Domain)

	if err := resp.Bind(&auth, "E1053877", base.Domain, basterm.UsernameAndPassword); err != nil {
		return
	}

	user, err := p.Service.Login(auth, params)
	if err != nil {
		resp.Error(err).JSON()
		resp.Record(base.LoginFailed, auth.Username, len(auth.Password))
		return
	}

	tmpUser := user
	tmpUser.Extra = nil

	resp.Record(base.BasLogin, tmpUser)
	resp.Status(http.StatusOK).
		Message(basterm.UserLogedInSuccessfully).
		JSON(user)
}

// Profile returns the user's information
func (p *AuthAPI) Profile(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basmodel.UserTable, base.Domain)

	var user basmodel.User
	var err error
	if user, err = p.Service.Profile(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewProfile)
	resp.Status(http.StatusOK).
		MessageT("profile").
		JSON(user)
}

// Logout will erase the resources from access.Cache
func (p *AuthAPI) Logout(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	params := param.Get(c, p.Engine, basterm.Users)
	p.Service.Logout(params)
	resp.Record(base.BasLogout)
	resp.Status(http.StatusOK).
		Message("user logged out").
		JSON()
}

// TemporaryToken is used for creating temporary access token for download excel and etc
func (p *AuthAPI) TemporaryToken(c *gin.Context) {
	// var auth basmodel.Auth
	resp := response.New(p.Engine, c, base.Domain)

	params := param.Get(c, p.Engine, "users")

	tmpKey, err := p.Service.TemporaryToken(params)
	if err != nil {
		resp.Status(http.StatusInternalServerError).Error(corerr.YouDontHavePermission).JSON()
		return
	}

	resp.Status(http.StatusOK).
		Message(corterm.TemporaryToken).
		JSON(tmpKey)

}

// TemporaryTokenHour is used for creating temporary access token for download excel and etc
func (p *AuthAPI) TemporaryTokenHour(c *gin.Context) {
	// var auth basmodel.Auth
	resp := response.New(p.Engine, c, base.Domain)

	var hour int
	var err error
	if hour, err = strconv.Atoi(c.Param("hour")); err != nil {
		err := limberr.New("hour is not a number", "E4269633").
			Domain(base.Domain).
			Message(corerr.VisNotValid, "hour").
			Custom(corerr.ForbiddenErr).Build()
		resp.Error(err).JSON()
		return
	}

	if p.Engine.Envs.ToInt(base.MaxHourTemporaryToken) < hour {
		err := limberr.New("hour exceeds the limit, choose smaller number", "E4285615").
			Domain(base.Domain).
			Message(corerr.VisNotValid, "hour").
			Custom(corerr.ForbiddenErr).Build()
		resp.Error(err).JSON()
		return
	}

	lang := dict.Lang(c.Param("lang"))

	if ok, _ := helper.Includes(dict.Langs, lang); !ok {
		err := limberr.New("language is not accepted", "E4244811").
			Domain(base.Domain).
			Message(corerr.AcceptedValueForVareV, dict.R(corterm.Language), dict.Langs).
			Custom(corerr.ValidationFailedErr).Build()
		resp.Error(err).JSON()
		return
	}

	tmpKey, err := p.Service.TemporaryTokenHour(hour, lang)
	if err != nil {
		resp.Status(http.StatusInternalServerError).Error(corerr.YouDontHavePermission).JSON()
		return
	}

	resp.Status(http.StatusOK).
		Message(corterm.TemporaryToken).
		JSON(tmpKey)
}

// Register a user
func (p *AuthAPI) Register(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var user, createdUser basmodel.User
	var err error

	if err = resp.Bind(&user, "E1032126", base.Domain, basterm.User); err != nil {

		return
	}

	if createdUser, err = p.Service.Register(user); err != nil {
		resp.Error(err).JSON()
		return
	}

	createdUser.Password = ""

	resp.RecordCreate(base.Register, createdUser)
	resp.Status(http.StatusOK).
		MessageT(basterm.UserRegisteredSuccessfully).
		JSON(user)
}
