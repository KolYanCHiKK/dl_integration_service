package app

import "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/app"

const (
	SSLModeDisable    app.SSLMode = "disable"
	SSLModeAllow      app.SSLMode = "allow"
	SSLModePrefer     app.SSLMode = "prefer"
	SSLModeRequire    app.SSLMode = "require"
	SSLModeVerifyCA   app.SSLMode = "verify-ca"
	SSLModeVerifyFull app.SSLMode = "verify-full"
)
