package insertdata

import (
	"omono/cmd/restapi/insertdata/table"
	"omono/internal/core"
)

// Insert is used for add static rows to database
func Insert(engine *core.Engine) {

	if engine.Envs.ToBool(core.AutoMigrate) {
		table.InsertCities(engine)
		table.InsertRoles(engine)
		table.InsertAccounts(engine)

		table.InsertUsers(engine)
		table.InsertSettings(engine)

	}

}
