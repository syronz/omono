package main

import (
	"flag"
	"github.com/syronz/dict"
	"omono/cmd/restapi/insertdata"
	"omono/cmd/restapi/server"
	"omono/cmd/restapi/startoff"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/eaccounting/eacrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/corstartoff"
	"omono/pkg/glog"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {

	engine := startoff.LoadEnvs()

	glog.Init(engine.Envs[core.ServerLogFormat],
		engine.Envs[core.ServerLogOutput],
		engine.Envs[core.ServerLogLevel],
		engine.Envs.ToBool(core.ServerLogJSONIndent),
		true)

	dict.Init(engine.Envs[core.TermsPath], engine.Envs.ToBool(core.TranslateInBackend))

	corstartoff.ConnectDB(engine, false)
	corstartoff.ConnectActivityDB(engine)
	engine.ActivityCh = make(chan basmodel.Activity, 1)

	startoff.Migrate(engine)

	insertdata.Insert(engine)

	activityRepo := basrepo.ProvideActivityRepo(engine)
	basActivityServ := service.ProvideBasActivityService(activityRepo)
	//ActivityWatcher is use a channel for checking all activities for recording
	go basActivityServ.ActivityWatcher()

	engine.TransactionCh = make(chan eacmodel.TransactionCh, 1)
	phoneServ := service.ProvideBasPhoneService(basrepo.ProvidePhoneRepo(engine))
	accountServ := service.ProvideBasAccountService(basrepo.ProvideAccountRepo(engine), phoneServ)
	currencyServ := service.ProvideEacCurrencyService(eacrepo.ProvideCurrencyRepo(engine))
	slotServ := service.ProvideEacSlotService(eacrepo.ProvideSlotRepo(engine), currencyServ, accountServ)
	transactionRepo := eacrepo.ProvideTransactionRepo(engine)
	transactionServ := service.ProvideEacTransactionService(transactionRepo, slotServ)
	go transactionServ.JournalEntryWatcher()

	/*
		//init of views
		view.InitViewReports(engine)
		view.InitDasboardViews(engine)
		//init of procedures
		procedure.InitDashboardProcedure(engine)
		procedure.InitReportProcedure(engine)
		//init of events
		event.InitdashboardEvent(engine)
		event.InitreportEvent(engine)
	*/

	corstartoff.LoadSetting(engine)

	server.Start(engine)

}
