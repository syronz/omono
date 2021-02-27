package startoff

import (
	"omono/domain/base"
	"omono/domain/notification"
	"omono/internal/core"
	"omono/internal/types"
	"os"
)

// LoadEnvs load environment from env|JSON file
func LoadEnvs() *core.Engine {
	var engine core.Engine
	var envs types.Envs

	envs = make(types.Envs, 29)

	envs[core.Port] = os.Getenv("OMONO_CORE_PORT")
	envs[core.Addr] = os.Getenv("OMONO_CORE_ADDR")
	envs[core.DatabaseDataWriteDSN] = os.Getenv("OMONO_CORE_DATABASE_DATA_WRITE_DSN")
	envs[core.DatabaseDataReadDSN] = os.Getenv("OMONO_CORE_DATABASE_DATA_READ_DSN")
	envs[core.DatabaseDataType] = os.Getenv("OMONO_CORE_DATABASE_DATA_TYPE")
	envs[core.DatabaseDataLog] = os.Getenv("OMONO_CORE_DATABASE_DATA_LOG")
	envs[core.DatabaseActivityDSN] = os.Getenv("OMONO_CORE_DATABASE_ACTIVITY_DSN")
	envs[core.DatabaseActivityType] = os.Getenv("OMONO_CORE_DATABASE_ACTIVITY_TYPE")
	envs[core.DatabaseActivityLog] = os.Getenv("OMONO_CORE_DATABASE_ACTIVITY_LOG")
	envs[core.AutoMigrate] = os.Getenv("OMONO_CORE_AUTO_MIGRATE")
	envs[core.ServerLogFormat] = os.Getenv("OMONO_CORE_SERVER_LOG_FORMAT")
	envs[core.ServerLogOutput] = os.Getenv("OMONO_CORE_SERVER_LOG_OUTPUT")
	envs[core.ServerLogLevel] = os.Getenv("OMONO_CORE_SERVER_LOG_LEVEL")
	envs[core.ServerLogJSONIndent] = os.Getenv("OMONO_CORE_SERVER_LOG_JSON_INDENT")
	envs[core.APILogFormat] = os.Getenv("OMONO_CORE_API_LOG_FORMAT")
	envs[core.APILogOutput] = os.Getenv("OMONO_CORE_API_LOG_OUTPUT")
	envs[core.APILogLevel] = os.Getenv("OMONO_CORE_API_LOG_LEVEL")
	envs[core.APILogJSONIndent] = os.Getenv("OMONO_CORE_API_LOG_JSON_INDENT")
	envs[core.TermsPath] = os.Getenv("OMONO_CORE_TERMS_PATH")
	envs[core.DefaultLang] = os.Getenv("OMONO_CORE_DEFAULT_LANGUAGE")
	envs[core.TranslateInBackend] = os.Getenv("OMONO_CORE_TRANSLATE_IN_BACKEND")
	envs[core.ExcelMaxRows] = os.Getenv("OMONO_CORE_EXCEL_MAX_ROWS")
	envs[core.ErrPanel] = os.Getenv("OMONO_CORE_ERR_PANEL")
	envs[core.OriginalError] = os.Getenv("OMONO_CORE_ORIGINAL_ERROR")
	envs[core.GinMode] = os.Getenv("GIN_MODE")
	envs[core.URL] = os.Getenv("OMONO_CORE_URL")

	envs[base.PasswordSalt] = os.Getenv("OMONO_BASE_PASSWORD_SALT")
	envs[base.JWTSecretKey] = os.Getenv("OMONO_BASE_JWT_SECRET_KEY")
	envs[base.JWTExpiration] = os.Getenv("OMONO_BASE_JWT_EXPIRATION")
	envs[base.RecordRead] = os.Getenv("OMONO_BASE_RECORD_READ")
	envs[base.RecordWrite] = os.Getenv("OMONO_BASE_RECORD_WRITE")
	envs[base.ActivityFileCounter] = os.Getenv("OMONO_BASE_ACTIVITY_FILE_COUNTER")
	envs[base.ActivityTickTimer] = os.Getenv("OMONO_BASE_ACTIVITY_TICK_TIMER")
	envs[base.AdminUsername] = os.Getenv("OMONO_BASE_ADMIN_USERNAME")
	envs[base.AdminPassword] = os.Getenv("OMONO_BASE_ADMIN_PASSWORD")
	envs[base.DefaultUsersParentID] = os.Getenv("OMONO_BASE_DEFAULT_USER_PARENT_ID")

	envs[notification.AppURL] = os.Getenv("OMONO_NOTIFICATION_APP_URL")

	engine.Envs = envs

	return &engine
}
