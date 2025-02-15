package sessionStore

import (
	"github.com/mt1976/frantic-core/common"
	"github.com/mt1976/frantic-core/logger"
)

var domain = "SESSION" // table name in the database
var sessionExpiry = 20 // default to 20 mins

func init() {
	logger.InfoLogger.Printf("Initialising %v", domain)
	cfg := common.Get()
	sessionExpiry = cfg.Security.SessionExpiry
	if sessionExpiry == 0 {
		logger.InfoLogger.Printf("[%v] NO SESSION TIMEOUT, Session Life=[%v]", domain, sessionExpiry)
		return
	}
	logger.InfoLogger.Printf("Initialised %v", domain)
}
