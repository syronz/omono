package core

import "omono/internal/types"

// list of core environment keys
const (
	Port                 types.Envkey = "CORE_PORT"
	Addr                 types.Envkey = "CORE_ADDR"
	DatabaseDataWriteDSN types.Envkey = "DATABASE_DATA_WRITE_DSN"
	DatabaseDataReadDSN  types.Envkey = "DATABASE_DATA_READ_DSN"
	DatabaseDataType     types.Envkey = "DATABASE_DATA_TYPE"
	DatabaseDataLog      types.Envkey = "DATABASE_DATA_LOG"
	DatabaseActivityDSN  types.Envkey = "DATABASE_ACTIVITY_DSN"
	DatabaseActivityType types.Envkey = "DATABASE_ACTIVITY_TYPE"
	DatabaseActivityLog  types.Envkey = "DATABASE_ACTIVITY_LOG"
	AutoMigrate          types.Envkey = "AUTO_MIGRATE"
	ServerLogFormat      types.Envkey = "SERVER_LOG_FORMAT"
	ServerLogOutput      types.Envkey = "SERVER_LOG_OUTPUT"
	ServerLogLevel       types.Envkey = "SERVER_LOG_LEVEL"
	ServerLogJSONIndent  types.Envkey = "SERVER_LOG_JSON_INDENT"
	APILogFormat         types.Envkey = "API_LOG_FORMAT"
	APILogOutput         types.Envkey = "API_LOG_OUTPUT"
	APILogLevel          types.Envkey = "API_LOG_LEVEL"
	APILogJSONIndent     types.Envkey = "API_LOG_JSON_INDENT"
	TermsPath            types.Envkey = "TERMS_PATH"
	DefaultLang          types.Envkey = "DEFAULT_LANGUAGE"
	TranslateInBackend   types.Envkey = "TRANSLATE_IN_BACKEND"
	ExcelMaxRows         types.Envkey = "EXCEL_MAX_ROWS"
	ErrPanel             types.Envkey = "ERR_PANEL"
	OriginalError        types.Envkey = "ORIGINAL_ERROR"
	GinMode              types.Envkey = "GIN_MODE"
	URL                  types.Envkey = "URL"
)
