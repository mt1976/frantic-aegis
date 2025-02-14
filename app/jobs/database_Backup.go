package jobs

import (
	"strings"
	"time"

	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/date"
	"github.com/mt1976/frantic-core/io"
	"github.com/mt1976/frantic-core/logger"
	"github.com/mt1976/frantic-core/paths"
	"github.com/mt1976/frantic-core/timing"
)

type databaseBackup struct {
}

func (job databaseBackup) Run() error {
	//name := "Database Backup"
	announceJob(job, "Started")
	j := timing.Start(job.Name(), "Backup", "All")

	dateTime := time.Now().Format(date.Format.BackupFolder)
	logger.ServiceLogger.Printf("[%v] [BACKUP] Date=[%v]", domain, dateTime)

	destPath := paths.Backups().String() + paths.Seperator() + dateTime
	fullBackupPath := paths.Application().String() + destPath

	//create a folder
	err := io.MkDir(fullBackupPath)
	if err != nil {
		logger.ErrorLogger.Printf("[%v] [%v] Error=[%v]", domain, strings.ToUpper(job.Name()), err.Error())
	}

	// texts.Backup(destPath)
	// status.Backup(destPath)
	// users.Backup(destPath)
	// settings.Backup(destPath)
	//hosts.Backup(destPath)
	database.Backup(destPath)

	j.Stop(6)
	//b.lastRun = time.Now()

	announceJob(job, "Completed")
	NextRun(job)
	return nil
}

func (job databaseBackup) Service() func() {
	return func() {
		job.Run()
	}
}

func (job databaseBackup) Schedule() string {
	return "55 11 * * *"
}

func (job databaseBackup) Name() string {
	return "Scheduled Database Backup"
}
