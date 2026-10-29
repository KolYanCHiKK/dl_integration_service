package interfaces

import "time"

type LoginParameters interface {
	GetAuthTime() time.Time
	GetIp() string
	GetUserAgent() string
	GetExpiresAt() time.Time
}
