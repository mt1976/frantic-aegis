package security

import (
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
	trnsl8r "github.com/mt1976/trnsl8r_connect"
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
var cfg *commonConfig.Settings
var msgTypeKey string
var msgTitleKey string
var msgContentKey string
var msgActionKey string
var trnsl8 trnsl8r.Request

func init() {
	clock := timing.Start("security", "Initialise", "")
	logHandler.TranslationLogger.Println("Initialised")
	cfg = commonConfig.Get()
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

	appModeDev = cfg.IsApplicationMode(commonConfig.MODE_DEVELOPMENT)
	if appModeDev {
		logHandler.SecurityLogger.Printf("sessionKey: %v\n", sessionKey)
		logHandler.SecurityLogger.Printf("sessionUserIDKey: %v\n", sessionUserIDKey)
		logHandler.SecurityLogger.Printf("sessionUserCodeKey: %v\n", sessionUserCodeKey)
		logHandler.SecurityLogger.Printf("sessionTokenKey: %v\n", sessionTokenKey)
	}

	msgTypeKey = cfg.GetMessageTypeKey()
	msgTitleKey = cfg.GetMessageTitleKey()
	msgContentKey = cfg.GetMessageContentKey()
	msgActionKey = cfg.GetMessageActionKey()

	serverPort = cfg.GetServerPortAsString()
	serverHost = cfg.GetServerHost()
	serverProtocol = cfg.GetServerProtocol()

	database.ConnectToNamedDB("aegis")

	trnsServerProtocol := cfg.GetTranslationServerProtocol()
	trnsServerHost := cfg.GetTranslationServerHost()
	trnsServerPort := cfg.GetTranslationServerPort()
	trnsLocale := cfg.GetTranslationLocale()
	err := error(nil)
	trnsl8, err = trnsl8r.NewRequest().FromOrigin(appName).WithHost(trnsServerHost).WithPort(trnsServerPort).WithProtocol(trnsServerProtocol).WithLogger(logHandler.TranslationLogger).WithFilter(trnsl8r.LOCALE, trnsLocale)
	if err != nil {
		panic(err)
	}
	clock.Stop(1)
}
