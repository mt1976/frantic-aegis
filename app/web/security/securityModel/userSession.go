package securityModel

import (
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
)

type Session struct {
	SessionID string
	Life      time.Duration
	UserKey   string
	UserCode  string
	Token     sessionStore.Session_Store
}
