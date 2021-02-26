package eacapi

import (
	"net/http"
	"omono/domain/eaccounting"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacterm"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/response"

	"github.com/gin-gonic/gin"
)

// IncomeStatementAPI for injecting IncomeStatement service
type IncomeStatementAPI struct {
	Service service.EacBalanceSheetServ
	Engine  *core.Engine
}

// ProvideIncomeStatementAPI is used in wire
func ProvideIncomeStatementAPI(c service.EacBalanceSheetServ) IncomeStatementAPI {
	return IncomeStatementAPI{Service: c, Engine: c.Engine}
}

//BalanceSheet will generate balance sheet report
func (p *IncomeStatementAPI) BalanceSheet(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, eacterm.Transactions, eaccounting.Domain)
	var err error
	var balanceSheet []eacmodel.BalanceSheet
	var level string
	if params.CompanyID, err = resp.GetCompanyID("E1440507"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	level = c.Param("level")

	if balanceSheet, err = p.Service.BalanceSheet(params, level); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(eaccounting.ViewTransaction)
	resp.Status(http.StatusOK).
		JSON(balanceSheet)
}
