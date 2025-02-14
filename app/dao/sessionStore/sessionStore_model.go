package sessionStore

import (
	"time"

	audit "github.com/mt1976/frantic-core/dao/audit"
)

// SessionStore represents a SessionStore entity.
type SessionStore struct {
	ID     string      `storm:"id"` // primary key with auto increment
	Raw    string      // raw ID before encoding
	UserID int         `storm:"index"` // user key
	Expiry time.Time   // expiry time
	Audit  audit.Audit // audit data
}
