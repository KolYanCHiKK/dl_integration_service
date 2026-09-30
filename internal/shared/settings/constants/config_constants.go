package constants

import "github.com/KolYanCHiKK/dl_integration_service/internal/shared/settings/types"

const (
	Disable    types.SSLMode = "disable"
	Allow      types.SSLMode = "allow"
	Prefer     types.SSLMode = "prefer"
	Require    types.SSLMode = "require"
	VerifyCa   types.SSLMode = "verify-ca"
	VerifyFull types.SSLMode = "verify-full"
)
