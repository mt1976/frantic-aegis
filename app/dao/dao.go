package dao

import (
	"context"

	storm "github.com/asdine/storm/v3"
	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var name = "Session Management"
var Version = 1
var DB *storm.DB

func Initialise(cfg *commonConfig.Settings) error {
	clock := timing.Start(name, actions.INITIALISE.GetCode(), name)
	logHandler.InfoLogger.Printf("Initialising %v...", name)

	sessionStore.Initialise(context.TODO())

	logHandler.InfoLogger.Printf("Initialised %v", name)
	clock.Stop(1)
	return nil
}
