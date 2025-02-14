package jobs

import (
	"github.com/mt1976/frantic-core/common"
	cron3 "github.com/robfig/cron/v3"
)

var scheduledTasks *cron3.Cron
var cfg *common.Settings
var domain = "JOBS"
var appName string

func init() {
	scheduledTasks = cron3.New()
	cfg = common.Get()
	appName = cfg.GetApplicationName()
}
