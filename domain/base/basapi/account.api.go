package basapi

import (
	"fmt"
	"net/http"
	"omono/cmd/restapi/enum/settingfields"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/enum/accounttype"
	"omono/domain/base/message/basterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"

	"github.com/gin-gonic/gin"
)

// AccountAPI for injecting account service
type AccountAPI struct {
	Service service.BasAccountServ
	Engine  *core.Engine
}

// ProvideAccountAPI for account is used in wire
func ProvideAccountAPI(c service.BasAccountServ) AccountAPI {
	return AccountAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a account by it's id
func (p *AccountAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var account basmodel.Account
	var fix types.FixedNode

	if fix, err = resp.GetFixedNode(c.Param("accountID"), "E1070061", basterm.Account); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if account, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewAccount)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Account).
		JSON(account)
}

// List of accounts
func (p *AccountAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basmodel.AccountTable, base.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1056290"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListAccount)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Accounts).
		JSON(data)
}

// Create account
func (p *AccountAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var account, createdAccount basmodel.Account
	var err error

	if account.CompanyID, account.NodeID, err = resp.GetCompanyNode("E1057239", base.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if account.CompanyID, err = resp.GetCompanyID("E1085677"); err != nil {
		return
	}

	if !resp.CheckRange(account.CompanyID) {
		return
	}

	if err = resp.Bind(&account, "E1057541", base.Domain, basterm.Account); err != nil {
		return
	}

	if createdAccount, err = p.Service.Create(account); err != nil {
		resp.Error(err).JSON()
		return
	}

	// resp.RecordCreate(base.CreateAccount, account)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.Account).
		JSON(createdAccount)
}

// Update account
func (p *AccountAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error

	var account, accountBefore, accountUpdated basmodel.Account
	var fix types.FixedNode

	if fix, err = resp.GetFixedNode(c.Param("accountID"), "E1076703", basterm.Account); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&account, "E1086162", base.Domain, basterm.Account); err != nil {
		return
	}

	if accountBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	account.ID = fix.ID
	account.CompanyID = fix.CompanyID
	account.NodeID = fix.NodeID
	account.CreatedAt = accountBefore.CreatedAt
	if accountUpdated, err = p.Service.Save(account); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.UpdateAccount, accountBefore, account)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, basterm.Account).
		JSON(accountUpdated)
}

// Delete account
func (p *AccountAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var account basmodel.Account
	var fix types.FixedNode

	if fix, err = resp.GetFixedNode(c.Param("accountID"), "E1092196", basterm.Account); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if account, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.DeleteAccount, account)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.Account).
		JSON()
}

// Excel generate excel files based on search
func (p *AccountAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Accounts, base.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E1066535"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

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
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Code", "Type", "Status", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "Code", "Type", "Status", "UpdatedAt").
		WriteData(accounts).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ExcelAccount)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}

// ChartOfAccount is used cached chart of account for getting the last status of chart of accounts
func (p *AccountAPI) ChartOfAccount(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basmodel.AccountTable, base.Domain)

	data := make(map[string]interface{})
	var err error

	params.Select = "bas_accounts.id,bas_accounts.parent_id,bas_accounts.code,bas_accounts.name_ar,bas_accounts.name_en,bas_accounts.name_ku,bas_accounts.type"
	params.PreCondition = fmt.Sprintf("bas_accounts.type != '%v'", accounttype.Customer)
	if params.CompanyID, err = resp.GetCompanyID("E1080813"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	refresh := c.Query("refresh")
	if refresh == "true" {
		if data["list"], err = p.Service.ChartOfAccountRefresh(params); err != nil {
			resp.Error(err).JSON()
			return
		}
	} else {
		if data["list"], err = p.Service.ChartOfAccount(params); err != nil {
			resp.Error(err).JSON()
			return
		}
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Accounts).
		JSON(data)
}

// GetCashAccount is used for returning the id of default cash account
func (p *AccountAPI) GetCashAccount(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var account basmodel.Account

	fix := types.FixedNode{
		CompanyID: 1001,
		NodeID:    101,
		ID:        p.Engine.Setting[settingfields.CashAccountID].ToRowID(),
	}

	if account, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Account).
		JSON(account)
}

// SearchLeafs is used for finding the accounts ready for transactions
func (p *AccountAPI) SearchLeafs(c *gin.Context) {
	resp := response.New(p.Engine, c, base.Domain)
	var err error
	var accounts []basmodel.Account

	search := c.Query("search")

	lang := core.GetLang(c, p.Engine)

	if accounts, err = p.Service.SearchLeafs(search, lang); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Account).
		JSON(accounts)
}
