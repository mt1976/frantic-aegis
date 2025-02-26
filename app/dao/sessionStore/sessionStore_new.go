package sessionStore

// Data Access Object Session
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: Update the New function to implement the creation of a new domain entity
//TODO: Create any new functions required to support the domain entity

import (
	"context"
	"fmt"
	"time"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func New(ctx context.Context, userKey, userCode string) (Session_Store, error) {

	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

	//logHandler.InfoLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", userKey))

	sessionID := idHelpers.GetUUID()

	// Create a new struct
	record := Session_Store{}
	record.Key = idHelpers.Encode(sessionID)
	record.Raw = sessionID
	record.UserKey = userKey
	record.UserCode = userCode
	record.SessionID = record.Key

	record.Expiry = time.Now().Add(time.Minute * time.Duration(sessionExpiry))

	// Save the status instance to the database
	writeErr := record.insertOrUpdate(ctx, fmt.Sprintf("New %v created %v", domain, userKey), actions.CREATE.GetCode(), audit.CREATE, "Create")
	if writeErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOCreateError(domain, record.ID, writeErr))
		//	panic(writeErr)
	}

	//logHandler.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.CREATE, strings.ToUpper(domain), record.ID, fmt.Sprintf("New %v: %v", domain, field1))
	clock.Stop(1)
	return record, nil
}
