package insertdata

import (
	"omono/cmd/testinsertion/insertdata/table"
	"omono/internal/core"
)

// Insert is used for add static rows to database
func Insert(engine *core.Engine) {

	if engine.Envs.ToBool(core.AutoMigrate) {
		table.InsertSettings(engine)
		table.InsertRoles(engine)
		table.InsertAccounts(engine)
		table.InsertUsers(engine)
		table.InsertCurrencies(engine)
		// table.InsertTransactions(engine)
		table.InsertJournals(engine)
		table.InsertPhones(engine)
	}

}
