package dao

import (
	"strings"

	storm "github.com/asdine/storm/v3"
	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/common"
	"github.com/mt1976/frantic-core/logger"
)

var name = "DAO"
var Version = 1
var DB *storm.DB
var tableName = "database"

func Initialise(cfg *common.Settings) error {
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
	return nil
}
