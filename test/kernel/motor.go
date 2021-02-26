package kernel

import (
	"omono/internal/core"
	"omono/internal/corstartoff"
	"omono/pkg/dict"
	"omono/pkg/glog"
)

// StartMotor for generating engine special for TDD
func StartMotor(printQueries bool, debugLevel bool) *core.Engine {
	engine := LoadTestEnv()

	if debugLevel {
		engine.Envs[core.ServerLogLevel] = "trace"
	}

	glog.Init(engine.Envs[core.ServerLogFormat],
		engine.Envs[core.ServerLogOutput],
		engine.Envs[core.ServerLogLevel],
		engine.Envs.ToBool(core.ServerLogJSONIndent),
		true)

	dict.Init(engine.Envs[core.TermsPath], engine.Envs.ToBool(core.TranslateInBackend))

	corstartoff.ConnectDB(engine, printQueries)
	corstartoff.ConnectActivityDB(engine)

	// logparam.ServerLog(engine)
	// corstartoff.LoadTerms(engine)
	// logparam.ServerLog(engine)
	// corstartoff.ConnectDB(engine, printQueries)
	// corstartoff.ConnectActivityDB(engine)

	return engine
}
