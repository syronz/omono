package eacrepo

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/base/message/basterm"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacterm"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/param"
)

// IncomeStatementRepo for injecting engine
type IncomeStatementRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideIncomeStatementRepo is used in wire and initiate the Cols
func ProvideIncomeStatementRepo(engine *core.Engine) IncomeStatementRepo {
	return IncomeStatementRepo{
		Engine: engine,
	}
}

//GetAllAccountBalances will fetch all the balances from balance account
func (p *IncomeStatementRepo) GetAllAccountBalances(params param.Param) (balances []eacmodel.Balance, err error) {
	err = p.Engine.ReadDB.Table(eacmodel.BalanceTable).
		Where("company_id = ?", params.CompanyID).Find(&balances).Error

	err = p.dbError(err, "E1481050", eacmodel.Balance{}, corterm.List)

	return

}

// dbError is an internal method for generate proper dataeace error
func (p *IncomeStatementRepo) dbError(err error, code string, balance eacmodel.Balance, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, balance.AccountID, eacterm.BalanceSheet)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(eacterm.Transaction), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
