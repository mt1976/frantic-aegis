package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/logger"
	"github.com/mt1976/frantic-core/timing"
)

type SessionExpiryJob struct {
}

func (p SessionExpiryJob) Run() error {
	JobSessionExpiry()
	return nil
}

func (p SessionExpiryJob) Service() func() {
	return func() {
		p.Run()
	}
}

func (p SessionExpiryJob) Schedule() string {
	// Every 30 seconds
	return "0/30 * * * * *"
}

func (p SessionExpiryJob) Name() string {
	return "Session Expiry"
}

func JobSessionExpiry() {
	// Do something every day at midnight
	name := "Session"
	j := timing.Start(strings.ToUpper(name), "Expiry", "Monitor")

	sessionLifespan := cfg.Security.SessionExpiry
	if sessionLifespan == 0 {
		logger.ServiceLogger.Printf("[%v] NO SESSION TIMEOUT, Session Life=[%v]", strings.ToUpper(name), sessionLifespan)
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
	logger.ServiceLogger.Printf("[%v] Sessions=[%v]", strings.ToUpper(name), noSessions)
	if noSessions == 0 {
		logger.ServiceLogger.Printf("[%v] No sessions to process", strings.ToUpper(name))
		j.Stop(0)
		return
	}

	for x, s := range sessions {

		if s.Expiry.Before(time.Now()) {
			count++
			logger.ServiceLogger.Printf("[%v] NOK (%v/%v) Session=[%v] EXPIRED=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry)
			logger.SecurityLogger.Printf("[%v] NOK (%v/%v) Session=[%v] EXPIRED=[%v] User=[%v] Code=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry, s.UserID)
			err := s.Delete(context.TODO(), "Session Expired")
			if err != nil {
				logger.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
			}
		} else {
			logger.ServiceLogger.Printf("[%v]  OK (%v/%v) Session=[%v] Expires=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry)
			logger.SecurityLogger.Printf("[%v]  OK (%v/%v) Session=[%v] Expires=[%v] User=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry, s.UserID)
		}
	}
	j.Stop(count)
}
