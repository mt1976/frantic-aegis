package sessionStore

// Data Access Object Session
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"time"

	"github.com/mt1976/frantic-core/dao/audit"
)

// Session_Store represents a Session_Store entity.
type Session_Store struct {
	ID        int         `storm:"id,increment=100000"` // primary key with auto increment
	Key       string      `storm:"unique"`              // key, not used
	Raw       string      `storm:"index"`               // raw ID before encoding
	Audit     audit.Audit `csv:"-"`                     // audit data
	SessionID string      `storm:"index"`               // session key
	UserKey   string      `storm:"index"`               // user key
	UserCode  string      `storm:"index"`               // user code
	Expiry    time.Time   // expiry time

}

// Define the field set as names
var (
	FIELD_ID        = "ID"
	FIELD_Key       = "Key"
	FIELD_SessionID = "SessionID"
	FIELD_Raw       = "Raw"
	FIELD_UserID    = "UserID"
	FIELD_UserCode  = "UserCode"
	FIELD_Expiry    = "Expiry"
	FIELD_Audit     = "Audit"
)

var domain = "Session"
