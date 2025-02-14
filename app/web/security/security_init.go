package security

import (
	"github.com/mt1976/frantic-core/common"
	"github.com/mt1976/frantic-core/logger"
)

var serverPort string
var serverHost string
var serverProtocol string

var sessionUserIDKey string
var sessionKey string
var sessionUserCodeKey string
var sessionTokenKey string
var sessionExpiryKey string
var appModeDev = false
var trnsServerHost string
var trnsServerPort int
var trnsServerProtocol string
var trnsLocale string
var appName string
var cfg *common.Settings
var msgTypeKey string
var msgTitleKey string
var msgContentKey string
var msgActionKey string

func init() {

	logger.TranslationLogger.Println("Initialised")
	cfg = common.Get()
	trnsServerProtocol = cfg.GetTranslationServerProtocol()
	trnsServerHost = cfg.GetTranslationServerHost()
	trnsServerPort = cfg.GetTranslationServerPort()
	trnsLocale = cfg.GetTranslationLocale()
	appName = cfg.GetApplicationName()

	// dao.Initialise()
	sessionKey = cfg.GetSecuritySessionKeyName()
	sessionUserIDKey = cfg.GetSecuritySessionUserIDKey()
	sessionUserCodeKey = cfg.GetSecuritySessionUserCodeKey()
	sessionTokenKey = cfg.GetSecuritySessionTokenKey()

	appModeDev = cfg.IsApplicationMode(common.MODE_DEVELOPMENT)
	if appModeDev {
		logger.SecurityLogger.Printf("sessionKey: %v\n", sessionKey)
		logger.SecurityLogger.Printf("sessionUserIDKey: %v\n", sessionUserIDKey)
		logger.SecurityLogger.Printf("sessionUserCodeKey: %v\n", sessionUserCodeKey)
		logger.SecurityLogger.Printf("sessionTokenKey: %v\n", sessionTokenKey)
	}

	msgTypeKey = cfg.GetMessageTypeKey()
	msgTitleKey = cfg.GetMessageTitleKey()
	msgContentKey = cfg.GetMessageContentKey()
	msgActionKey = cfg.GetMessageActionKey()

	serverPort = cfg.GetServerPortAsString()
	serverHost = cfg.GetServerHost()
	serverProtocol = cfg.GetServerProtocol()

}
