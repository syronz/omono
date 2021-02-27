package main

import (
	"omono/cmd/restapi/insertdata"
	"omono/cmd/restapi/server"
	"omono/cmd/restapi/startoff"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/corstartoff"
	"omono/pkg/glog"

	"github.com/syronz/dict"
)

func main() {
	engine := startoff.LoadEnvs()

	// set glog as a global variable for logging the errors and debug
	glog.Init(engine.Envs[core.ServerLogFormat],
		engine.Envs[core.ServerLogOutput],
		engine.Envs[core.ServerLogLevel],
		engine.Envs.ToBool(core.ServerLogJSONIndent),
		true)

	// load terms
	dict.Init(engine.Envs[core.TermsPath], engine.Envs.ToBool(core.TranslateInBackend))

	// connect the database
	corstartoff.ConnectDB(engine, false)
	corstartoff.ConnectActivityDB(engine)

	// migrate the database
	startoff.Migrate(engine)

	// insert basic data
	insertdata.Insert(engine)

	// ActivityWatcher is use a channel for checking all activities for recording
	engine.ActivityCh = make(chan basmodel.Activity, 1)
	activityRepo := basrepo.ProvideActivityRepo(engine)
	basActivityServ := service.ProvideBasActivityService(activityRepo)
	go basActivityServ.ActivityWatcher()

	// load setting
	corstartoff.LoadSetting(engine)

	// start the API
	server.Start(engine)
}
