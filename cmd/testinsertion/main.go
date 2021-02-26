package main

import (
	"flag"
	"omono/cmd/restapi/startoff"
	"omono/cmd/testinsertion/insertdata"
	"omono/internal/core"
	"omono/internal/corstartoff"
	"omono/pkg/dict"
	"omono/pkg/glog"
	"omono/test/kernel"
)

var noReset bool
var logQuery bool

func init() {
	flag.BoolVar(&noReset, "noReset", false, "by default it drop tables before migrate")
	flag.BoolVar(&logQuery, "logQuery", false, "print queries in gorm")
}

func main() {
	flag.Parse()

	engine := kernel.LoadTestEnv()

	glog.Init(engine.Envs[core.ServerLogFormat],
		engine.Envs[core.ServerLogOutput],
		engine.Envs[core.ServerLogLevel],
		engine.Envs.ToBool(core.ServerLogJSONIndent),
		true)

	dict.Init(engine.Envs[core.TermsPath], engine.Envs.ToBool(core.TranslateInBackend))

	corstartoff.ConnectDB(engine, logQuery)
	corstartoff.ConnectActivityDB(engine)
	startoff.Migrate(engine)
	insertdata.Insert(engine)

	if noReset {
		glog.Debug("Data has been migrated successfully (no reset)")
		// fmt.Println("Data has been migrated successfully (no reset)")
	} else {
		glog.Debug("Data has been reset successfully")
		// fmt.Println("Data has been reset successfully")
	}

}
