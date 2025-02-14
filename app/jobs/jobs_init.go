package jobs

import (
	"github.com/mt1976/frantic-core/common"
	"github.com/mt1976/frantic-core/logger"
	trnsl8r "github.com/mt1976/trnsl8r_connect"
	cron3 "github.com/robfig/cron/v3"
)

var scheduledTasks *cron3.Cron
var cfg *common.Settings
var domain = "JOBS"
var appName string
var trnsl8 trnsl8r.Request

func init() {
	scheduledTasks = cron3.New()
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
