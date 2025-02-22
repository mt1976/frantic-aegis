package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

type SessionExpiryJob struct {
	databaseAccessors []func() ([]*database.DB, error)
}

// AddFunction implements jobs.Job.
func (job *SessionExpiryJob) AddDatabaseAccessFunctions(fn func() ([]*database.DB, error)) {
	logHandler.ServiceLogger.Printf("[%v] [%v] Adding Function", domain, job.Name())
	job.databaseAccessors = append(job.databaseAccessors, fn)
	logHandler.ServiceLogger.Printf("[%v] [%v] Function Added - No Funcs=(%v)", domain, job.Name(), len(job.databaseAccessors))
}

// Description implements jobs.Job.
func (j *SessionExpiryJob) Description() string {
	return "Session Expiry Process, Deletes Expired Sessions"
}

func (j *SessionExpiryJob) Run() error {
	clock := timing.Start(jobs.CodedName(j), actions.PROCESS.GetCode(), j.Description())
	jobs.PreRun(j)
	count := 0

	for _, f := range j.databaseAccessors {
		dbList, err := f()
		if err != nil {
			logHandler.ErrorLogger.Printf("[%v] Error=[%v]", j.Name(), err.Error())
			continue
		}
		for _, db := range dbList {
			JobSessionExpiry(j, *db)
			count++
		}
	}

	jobs.PostRun(j)
	clock.Stop(count)
	return nil
}

func (j *SessionExpiryJob) Service() func() {
	return func() {
		err := j.Run()
		if err != nil {
			logHandler.ErrorLogger.Printf("[%v] Error=[%v]", j.Name(), err.Error())
			panic(err)
		}
	}
}

func (j *SessionExpiryJob) Schedule() string {
	// Every 30 seconds
	return "0/30 * * * * *"
}

func (j *SessionExpiryJob) Name() string {
	returnValue, _ := translationServiceRequest.Get("Session Expiry")
	return returnValue.String()
}

func JobSessionExpiry(j jobs.Job, db database.DB) {
	// Do something every day at midnight
	name := jobs.CodedName(j)
	clock := timing.Start(name, actions.START.GetCode(), j.Description())

	sessionLifespan := cfg.Security.SessionExpiry
	if sessionLifespan == 0 {
		logHandler.ServiceLogger.Printf("[%v] NO SESSION TIMEOUT, Session Life=[%v]", strings.ToUpper(name), sessionLifespan)
		return
	}
	logHandler.ServiceLogger.Printf("[%v] Session Life=[%v]", strings.ToUpper(name), sessionLifespan)
	count := 0

	// Get all the sessions
	// For each session, check the expiry date
	// If the expiry date is less than now, then delete the session
	var sessions []sessionStore.Session_Store

	err := db.GetAll(&sessions)
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
			logHandler.ServiceLogger.Printf("[%v] NOK (%v/%v) Session=[%v] EXPIRED=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry)
			logHandler.SecurityLogger.Printf("[%v] NOK (%v/%v) Session=[%v] EXPIRED=[%v] User=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry, s.UserID)
			err := sessionStore.Delete(context.TODO(), s.ID, "Session Expired")
			if err != nil {
				logHandler.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(name), err.Error())
			}
		} else {
			logHandler.ServiceLogger.Printf("[%v]  OK (%v/%v) Session=[%v] Expires=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry)
			logHandler.SecurityLogger.Printf("[%v]  OK (%v/%v) Session=[%v] Expires=[%v] User=[%v]", strings.ToUpper(name), x+1, noSessions, s.ID, s.Expiry, s.UserID)
		}
	}
	clock.Stop(count)
}
