package mappings

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
)

type SessionParams struct {
	SessionId   user.LoginSessionId `gorm:"column:id"`
	RedirectURL string              `gorm:"column:redirect_url"`
}

type LoginRequestParams struct {
	LoginRequest user.LoginRequestId
	AuthTime     time.Time
	IP           string
	UserAgent    string
	ExpiresAt    time.Time
}
type UserLoginParams struct {
	UserId       user.UserId `gorm:"column:id"`
	UserName     string      `gorm:"column:login"`
	UserPassword string      `gorm:"column:hashed_password"`
}
