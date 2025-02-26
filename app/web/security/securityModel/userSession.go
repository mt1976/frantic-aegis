package securityModel

import (
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
)

type Session struct {
	SessionID    string
	Expiry       time.Time
	UserKey      string
	UserCode     string
	SessionToken sessionStore.Session_Store
}
