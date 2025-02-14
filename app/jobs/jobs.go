package jobs

import (
	"github.com/mt1976/frantic-core/logger"
)

var Template Job = template{} // This is a template for other jobs.

// var DatabaseBackup Job = databaseBackup{}
// var DatabasePrune Job = databasePrune{}
var SessionExpiry Job = sessionExpiry{}

func Start() {
	// Start the job
	logger.ServiceLogger.Printf("[%v] Queue - Starting", domain)
	// Database Backup
	//	Schedule(DatabaseBackup)
	// Prune the archive of backups
	//	Schedule(DatabasePrune)
	// Check the status of the hosts
	// Start all the background jobs
	Schedule(SessionExpiry)
	scheduledTasks.Start()
	logger.ServiceLogger.Printf("[%v] Queue - Started", domain)
}

func Schedule(j Job) {
	// Start the job
	scheduledTasks.AddFunc(j.Schedule(), j.Service())
	announceJob(j, "Scheduled")
	NextRun(j)
}
