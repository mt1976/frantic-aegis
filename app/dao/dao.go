package dao

import (
	"strings"

	storm "github.com/asdine/storm/v3"
	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/commonConfig"
	logger "github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var name = "DAO"
var Version = 1
var DB *storm.DB
var tableName = "database"

func Initialise(cfg *commonConfig.Settings) error {
	clock := timing.Start(name, "Initialise", "")
	logger.InfoLogger.Printf("[%v] Initialising...", strings.ToUpper(name))

	// Preload the status store
	logger.InfoBanner(name, "Status", "Importing")
	// err := statusStore.ImportCSV()
	// if err != nil {
	// 	logger.ErrorLogger.Fatal(err.Error())
	// }

	sessionStore.Initialise()

	//routes.Initialise(cfg)

	logger.InfoLogger.Printf("[%v] Initialised", strings.ToUpper(name))
	clock.Stop(1)
	return nil
}
