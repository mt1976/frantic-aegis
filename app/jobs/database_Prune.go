package jobs

import (
	"fmt"
	"strings"
	"time"

	"github.com/mt1976/frantic-core/application"
	dates "github.com/mt1976/frantic-core/date"
	"github.com/mt1976/frantic-core/io"
	"github.com/mt1976/frantic-core/logger"
	"github.com/mt1976/frantic-core/notifications"
	"github.com/mt1976/frantic-core/paths"
	"github.com/mt1976/frantic-core/timing"
)

type databasePrune struct {
}

func (p databasePrune) Run() error {
	jobPruneBackups()
	NextRun(p)
	return nil
}

func (p databasePrune) Service() func() {
	return func() {
		p.Run()
	}
}

func (p databasePrune) Schedule() string {
	return "25 0 * * *"
}

func (p databasePrune) Name() string {
	return "Scheduled Database Maintenance - Prune Old Backups"
}

func jobPruneBackups() {
	// Do something every day at midnight
	name := "Prune"
	j := timing.Start(strings.ToUpper(name), "Prune", "Old Backups")
	// Get Settings

	// set, err := settingStore.Get()
	// if err != nil {
	// 	logger.WarningLogger.Printf("[%v] [%v] Error=[%v]", domain, strings.ToUpper(name), err.Error())
	// 	return
	// }

	retainBackupDays := cfg.GetMaxHistoryEntries()
	logger.ServiceLogger.Printf("[%v] RetainBackupDays=[%v]", strings.ToUpper(name), retainBackupDays)
	today := StartOfDay(time.Now())

	// get today's date
	DMY := dates.Format.DMY
	todayStr := today.Format(DMY)
	logger.ServiceLogger.Printf("[%v] Today=[%v]", strings.ToUpper(name), todayStr)
	deleteBeforeDate := today.AddDate(0, 0, -retainBackupDays)
	deleteBeforeDateStr := deleteBeforeDate.Format(DMY)
	logger.ServiceLogger.Printf("[%v] DeleteBeforeDate=[%v]", strings.ToUpper(name), deleteBeforeDateStr)

	// Get Backups path
	path := paths.Backups().String()
	logger.ServiceLogger.Printf("[%v] Path=[%v]", strings.ToUpper(name), path)
	full := paths.Application().String()
	logger.ServiceLogger.Printf("[%v] AppPath=[%v]", strings.ToUpper(name), full)
	backupPath := full + path
	logger.ServiceLogger.Printf("[%v] BackupPath=[%v]", strings.ToUpper(name), backupPath)

	// Get all folders in the backup directory
	folders, err := io.Dir(backupPath)
	if err != nil {
		logger.WarningLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
		return
	}
	logger.ServiceLogger.Printf("[%v] No Folders=[%v]", strings.ToUpper(name), len(folders))
	count := 0
	// For each folder check if it is before the deleteBeforeDate
	for _, folder := range folders {
		// Get the date from the folder strings.ToUpper(name)
		backupDate, err := getDateFromBackupFolderName(folder)
		if err != nil {
			logger.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
			return
		}
		// Check if the backupDate is before the deleteBeforeDate
		if backupDate.Before(deleteBeforeDate) {
			// Delete the folder
			logger.ServiceLogger.Printf("[%v] Deleting=[%v] FolderDate=[%v] DeleteDate=[%v]", strings.ToUpper(name), folder, backupDate.Format(DMY), deleteBeforeDateStr)
			count++
			err := io.DeleteFolder(backupPath + folder)
			if err != nil {
				logger.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
				return
			}
			msg := "Backup Pruned Folder=[%v] On=[%v]"
			msg = fmt.Sprintf(msg, folder, application.HostName())
			title := "%v (%v) Backup Pruned"
			title = fmt.Sprintf(title, appName, application.HostName())
			err = notifications.Send(msg, title, 0)
			if err != nil {
				logger.WarningLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
				return
			}
		}
	}

	j.Stop(count)
}

func getDateFromBackupFolderName(folder string) (date time.Time, err error) {
	// Get the date from the folder strings.ToUpper(name)
	date, err = time.Parse(dates.Format.BackupFolder, folder)
	if err != nil {
		logger.ErrorLogger.Printf("[%v] [%v] Error=[%v]", domain, "BACKUP", err.Error())
		return
	}
	return
}
