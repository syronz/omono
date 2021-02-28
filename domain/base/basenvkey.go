package base

import "omono/internal/types"

// types for base environment keys
const (
	PasswordSalt          types.Envkey = "PASSWORD_SALT"
	JWTSecretKey          types.Envkey = "JWT_SECRET_KEY"
	JWTExpiration         types.Envkey = "JWT_EXPIRATION"
	RecordRead            types.Envkey = "RECORD_READ"
	RecordWrite           types.Envkey = "RECORD_WRITE"
	ActivityFileCounter   types.Envkey = "ACTIVITY_FILE_COUNTER"
	ActivityTickTimer     types.Envkey = "ACTIVITY_TICK_TIMER"
	AdminUsername         types.Envkey = "ADMIN_USERNAME"
	AdminPassword         types.Envkey = "ADMIN_PASSWORD"
	MaxHourTemporaryToken types.Envkey = "MAX_HOUR_TEMPORARY_TOKEN"
)
