package subapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"omono/domain/base/basterm"
	"omono/domain/service"
	"omono/domain/subscriber"
	"omono/domain/subscriber/submodel"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/pkg/helper/excel"
)

// AccountAPI for injecting account service
type AccountAPI struct {
	Service service.SubAccountServ
	Engine  *core.Engine
}

// ProvideAccountAPI for account is used in wire
func ProvideAccountAPI(c service.SubAccountServ) AccountAPI {
	return AccountAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a account by it's id
func (p *AccountAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var err error
	var account submodel.Account
	var id uint

	if id, err = resp.GetID(c.Param("accountID"), "E1070061", basterm.Account); err != nil {
		return
	}

	if account, err = p.Service.FindByID(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.ViewAccount)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Account).
		JSON(account)
}

// List of accounts
func (p *AccountAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, submodel.AccountTable, subscriber.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.ListAccount)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Accounts).
		JSON(data)
}

// Create account
func (p *AccountAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var account, createdAccount submodel.Account
	var err error

	if err = resp.Bind(&account, "E1057541", subscriber.Domain, basterm.Account); err != nil {
		return
	}

	if createdAccount, err = p.Service.Create(account); err != nil {
		resp.Error(err).JSON()
		return
	}

	// resp.RecordCreate(subscriber.CreateAccount, account)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.Account).
		JSON(createdAccount)
}

// Update account
func (p *AccountAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var err error

	var account, accountBefore, accountUpdated submodel.Account
	var id uint

	if id, err = resp.GetID(c.Param("accountID"), "E1076703", basterm.Account); err != nil {
		return
	}

	if err = resp.Bind(&account, "E1086162", subscriber.Domain, basterm.Account); err != nil {
		return
	}

	account.ID = id
	account.CreatedAt = accountBefore.CreatedAt
	if accountUpdated, accountBefore, err = p.Service.Save(account); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.UpdateAccount, accountBefore, account)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, basterm.Account).
		JSON(accountUpdated)
}

// Delete account
func (p *AccountAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var err error
	var account submodel.Account
	var id uint

	if id, err = resp.GetID(c.Param("accountID"), "E1092196", basterm.Account); err != nil {
		return
	}

	if account, err = p.Service.Delete(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.DeleteAccount, account)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.Account).
		JSON()
}

// Excel generate excel files based on search
func (p *AccountAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Accounts, subscriber.Domain)
	var err error

	accounts, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("account")
	ex.AddSheet("Accounts").
		AddSheet("Summary").
		Active("Accounts").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "G", 15.3).
		SetColWidth("H", "H", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Accounts").
		WriteHeader("ID", "Name", "Code", "Type", "Status", "Updated At").
		SetSheetFields("ID", "Name", "Code", "Type", "Status", "UpdatedAt").
		WriteData(accounts).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.ExcelAccount)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
