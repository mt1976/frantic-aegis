package jobs

import (
	"time"

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/mt1976/frantic-core/date"
	"github.com/mt1976/frantic-core/logger"
)

func StartOfDay(t time.Time) time.Time {
	// Purpose: To remove the time from a date
	return date.StartOfDay(t)
}

func BeforeOrEqualTo(t1, t2 time.Time) bool {
	return date.IsBeforeOrEqualTo(t1, t2)
}

func AfterOrEqualTo(t1, t2 time.Time) bool {
	return date.IsAfterOrEqualTo(t1, t2)
}

func NextRun(j Job) string {
	// Purpose: To determine the next run time of a job
	//bkHuman1, _ := crondescriptor.NewCronDescriptor(j.Schedule())
	//nr, _ := bkHuman1.GetDescription(crondescriptor.Full)
	//nr := support.NextRun(j.Schedule())
	logger.ServiceLogger.Printf("[%v] [%v] NextRun=[%v]", domain, j.Name(), getFreqHuman(j.Schedule()))
	return ""
}

func announceJob(j Job, action string) {

	tName, _ := trnsl8.Get(j.Name())
	tAction, _ := trnsl8.Get(action)

	logger.ServiceLogger.Printf("[%v] [%v] %v", domain, tName.String(), tAction.String())
}

func getFreqHuman(freq string) string {
	bkHuman1, _ := crondescriptor.NewCronDescriptor(freq)
	bkHuman, _ := bkHuman1.GetDescription(crondescriptor.Full)
	return *bkHuman
}
