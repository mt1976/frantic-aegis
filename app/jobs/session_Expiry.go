package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/logger"
	"github.com/mt1976/frantic-core/timing"
)

type sessionExpiry struct {
}

func (p sessionExpiry) Run() error {
	jobSessionExpiry()
	NextRun(p)
	return nil
}

func (p sessionExpiry) Service() func() {
	return func() {
		p.Run()
	}
}

func (p sessionExpiry) Schedule() string {
	// Every 30 seconds
	return "0/30 * * ? * * *"
}

func (p sessionExpiry) Name() string {
	return "Session Expiry"
}

func jobSessionExpiry() {
	// Do something every day at midnight
	name := "Session"
	j := timing.Start(strings.ToUpper(name), "Expiry", "Monitor")

	sessionLifespan := cfg.Security.SessionExpiry
	if sessionLifespan == 0 {
		logger.SecurityLogger.Printf("[%v] NO SESSION TIMEOUT, Session Life=[%v]", strings.ToUpper(name), sessionLifespan)
		return
	}
	logger.ServiceLogger.Printf("[%v] Session Life=[%v]", strings.ToUpper(name), sessionLifespan)
	count := 0

	// Get all the sessions
	// For each session, check the expiry date
	// If the expiry date is less than now, then delete the session

	sessions, err := sessionStore.GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
		return
	}

	noSessions := len(sessions)
	logger.SecurityLogger.Printf("[%v] Sessions=[%v]", strings.ToUpper(name), noSessions)

	for _, s := range sessions {
		if s.Expiry.Before(time.Now()) {
			count++
			logger.SecurityLogger.Printf("[%v] Session=[%v] Expired=[%v]", strings.ToUpper(name), s.ID, s.Expiry)
			err := s.Delete(context.TODO(), "Session Expired")
			if err != nil {
				logger.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
			}
		}
	}

	j.Stop(count)
}
