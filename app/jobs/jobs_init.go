package jobs

import (
	"github.com/mt1976/frantic-core/commonConfig"
	logger "github.com/mt1976/frantic-core/logHandler"
	trnsl8r "github.com/mt1976/trnsl8r_connect"
)

var cfg *commonConfig.Settings
var appName string
var translationServiceRequest trnsl8r.Request

func init() {
	cfg = commonConfig.Get()
	appName = cfg.GetApplicationName()
	trnsServerProtocol := cfg.GetTranslationServerProtocol()
	trnsServerHost := cfg.GetTranslationServerHost()
	trnsServerPort := cfg.GetTranslationServerPort()
	trnsLocale := cfg.GetTranslationLocale()
	err := error(nil)
	translationServiceRequest, err = trnsl8r.NewRequest().FromOrigin(appName).WithHost(trnsServerHost).WithPort(trnsServerPort).WithProtocol(trnsServerProtocol).WithLogger(logger.TranslationLogger).WithFilter(trnsl8r.LOCALE, trnsLocale)
	if err != nil {
		panic(err)
	}
}
