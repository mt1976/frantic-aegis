package sessionStore

import (
	"github.com/mt1976/frantic-core/commonConfig"
	logger "github.com/mt1976/frantic-core/logHandler"
)

var domain = "SESSION" // table name in the database
var sessionExpiry = 20 // default to 20 mins

func init() {
	logger.InfoLogger.Printf("Initialising %v", domain)
	cfg := commonConfig.Get()
	sessionExpiry = cfg.Security.SessionExpiry
	if sessionExpiry == 0 {
		logger.InfoLogger.Printf("[%v] NO SESSION TIMEOUT, Session Life=[%v]", domain, sessionExpiry)
		return
	}
	logger.InfoLogger.Printf("Initialised %v", domain)
}
