package consts

import (
	"math"
)

// constants which used inside the app
const (
	MinimumPasswordChar = 8

	// TemporaryTokenDuration = 100 * 100000 //in seconds
	TemporaryTokenDuration = 10

	MaxRowsCount = 1 << 62

	// MinFloat64 = k
	MinFloat64 = -1 * math.MaxFloat64

	DefaultLimit    = 100
	ShowDeletedRows = false

	DateLayout         = "2006-01-02"
	DateTimeLayout     = "2006-01-02 15:04:05"
	DateTimeLayoutZone = "2006-01-02 15:04:05 -0700"

	UserSuperAdminID   = 79
	UserResultViewerID = 12

	// it is used in chart of accounts after this numbers show more button emerge
	MaxChildrenForChartOfAccounts = 20
)
