package sessionStore

import (
	"context"
	"strings"
	"time"

	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

// Data Access Object Session
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: Implement the validate function to process the domain entity
//TODO: Implement the calculate function to process the domain entity
//TODO: Implement the isDuplicateOf function to process the domain entity
//TODO: Implement the postGetProcessing function to process the domain entity

func (record *Session_Store) upgradeProcessing() error {
	//TODO: Add any upgrade processing here
	//This processing is triggered directly after the record has been retrieved from the database
	return nil
}

func (record *Session_Store) defaultProcessing() error {
	//TODO: Add any default calculations here
	//This processing is triggered directly before the record is saved to the database
	return nil
}

func (record *Session_Store) validationProcessing() error {
	//TODO: Add any record validation here
	//This processing is triggered directly before the record is saved to the database and after the default calculations
	return nil
}

func (h *Session_Store) postGetProcessing() error {
	//TODO: Add any post get processing here
	//This processing is triggered directly after the record has been retrieved from the database and after the upgrade processing
	return nil
}

func (record *Session_Store) preDeleteProcessing() error {
	//TODO: Add any pre delete processing here
	//This processing is triggered directly before the record is deleted from the database
	return nil
}

func cloneProcessing(ctx context.Context, source Session_Store) (Session_Store, error) {
	//TODO: Add any clone processing here
	panic("Not Implemented")
	return Session_Store{}, nil
}

func jobProcessor(j jobs.Job) {
	name := jobs.CodedName(j)
	clock := timing.Start(name, actions.RUN.GetCode(), j.Description())

	sessionLifespan := cfg.GetSecuritySession_ExpiryPeriod()
	if sessionLifespan == 0 {
		logHandler.ServiceLogger.Printf("[%v] NO SESSION TIMEOUT, Session Life=[%v]", strings.ToUpper(name), sessionLifespan)
		return
	}
	logHandler.ServiceLogger.Printf("[%v] Session Life=[%v]", strings.ToUpper(name), sessionLifespan)
	count := 0

	// Get all the sessions
	// For each session, check the expiry date
	// If the expiry date is less than now, then delete the session
	var sessions []Session_Store

	sessions, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
		clock.Stop(0)
		return
	}

	noSessions := len(sessions)
	logHandler.ServiceLogger.Printf("[%v] Sessions=[%v]", strings.ToUpper(name), noSessions)
	if noSessions == 0 {
		logHandler.ServiceLogger.Printf("[%v] No sessions to process", strings.ToUpper(name))
		clock.Stop(0)
		return
	}

	for x, s := range sessions {

		if s.Expiry.Before(time.Now()) {
			count++
			newFunction(name, x, noSessions, s)
			continue
		}

		if s.Expiry.After(time.Now().Add(time.Minute * time.Duration(sessionLifespan))) {
			count++
			newFunction(name, x, noSessions, s)
			continue
		}

		logHandler.ServiceLogger.Printf("[%v]  OK (%v/%v) Session=[%v] Expires=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry)
		logHandler.SecurityLogger.Printf("[%v]  OK (%v/%v) Session=[%v] Expires=[%v] User=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry, s.UserCode)

	}
	clock.Stop(count)
}

func newFunction(name string, x int, noSessions int, s Session_Store) {
	logHandler.ServiceLogger.Printf("[%v] NOK (%v/%v) Session=[%v] EXPIRED=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry)
	logHandler.SecurityLogger.Printf("[%v] NOK (%v/%v) Session=[%v] EXPIRED=[%v] User=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry, s.UserCode)
	err := Delete(context.TODO(), s.ID, "Session Expired")
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
	}
}
