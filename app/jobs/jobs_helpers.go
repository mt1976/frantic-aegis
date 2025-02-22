package jobs

import (
	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/dao/database"
)

func GetDB() func() (*database.DB, error) {
	//logHandler.InfoLogger.Println("GETDB")
	return sessionStore.GetDB()
}
