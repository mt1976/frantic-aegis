package security

import (
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
	trnsl8r "github.com/mt1976/trnsl8r_connect"
)

var serverPort string
var serverHost string
var serverProtocol string

var sessionUserKeyKey string
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
	clock := timing.Start(domain, actions.INITIALISE.GetCode(), "")

	cfg = commonConfig.Get()

	trnsServerProtocol = cfg.GetTranslationServer_Protocol()
	trnsServerHost = cfg.GetTranslationServer_Host()
	trnsServerPort = cfg.GetTranslationServer_Port()
	trnsLocale = cfg.GetApplication_Locale()
	appName = cfg.GetApplication_Name()

	// dao.Initialise()
	sessionKey = cfg.GetSecuritySessionKey_Session()
	sessionUserKeyKey = cfg.GetSecuritySessionKey_UserKey()
	sessionUserCodeKey = cfg.GetSecuritySessionKey_UserCode()
	sessionTokenKey = cfg.GetSecuritySessionKey_Token()
	sessionExpiryKey = cfg.GetSecuritySessionKey_ExpiryPeriod()

	appModeDev = cfg.IsApplicationMode(commonConfig.MODE_DEVELOPMENT)
	if appModeDev {
		logHandler.SecurityLogger.Printf("sessionKey: %v\n", sessionKey)
		logHandler.SecurityLogger.Printf("sessionUserIDKey: %v\n", sessionUserKeyKey)
		logHandler.SecurityLogger.Printf("sessionUserCodeKey: %v\n", sessionUserCodeKey)
		logHandler.SecurityLogger.Printf("sessionTokenKey: %v\n", sessionTokenKey)
		logHandler.SecurityLogger.Printf("sessionExpiryKey: %v\n", sessionExpiryKey)
	}

	msgTypeKey = cfg.GetMessageKey_Type()
	msgTitleKey = cfg.GetMessageKey_Title()
	msgContentKey = cfg.GetMessageKey_Content()
	msgActionKey = cfg.GetMessageKey_Action()

	serverPort = cfg.GetServer_PortString()
	serverHost = cfg.GetServer_Host()
	serverProtocol = cfg.GetServer_Protocol()

	database.ConnectToNamedDB("aegis")

	err := error(nil)
	trnsl8, err = trnsl8r.NewRequest().FromOrigin(appName).WithHost(trnsServerHost).WithPort(trnsServerPort).WithProtocol(trnsServerProtocol).WithLogger(logHandler.TranslationLogger).WithFilter(trnsl8r.LOCALE, trnsLocale)
	if err != nil {
		panic(err)
	}

	logHandler.EventLogger.Printf("Initialised %v using '%v'", domain, appName)

	clock.Stop(1)
}
