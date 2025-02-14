package jobs

import (
	"github.com/mt1976/frantic-core/common"
	"github.com/mt1976/frantic-core/logger"
	trnsl8r "github.com/mt1976/trnsl8r_connect"
)

var cfg *common.Settings
var appName string
var trnsl8 trnsl8r.Request

func init() {
	cfg = common.Get()
	appName = cfg.GetApplicationName()
	trnsServerProtocol := cfg.GetTranslationServerProtocol()
	trnsServerHost := cfg.GetTranslationServerHost()
	trnsServerPort := cfg.GetTranslationServerPort()
	trnsLocale := cfg.GetTranslationLocale()
	err := error(nil)
	trnsl8, err = trnsl8r.NewRequest().FromOrigin(appName).WithHost(trnsServerHost).WithPort(trnsServerPort).WithProtocol(trnsServerProtocol).WithLogger(logger.TranslationLogger).WithFilter(trnsl8r.LOCALE, trnsLocale)
	if err != nil {
		panic(err)
	}
}
