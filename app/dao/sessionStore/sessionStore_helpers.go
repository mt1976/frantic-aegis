package sessionStore

// Data Access Object Session - Sessionstore
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"
	"fmt"
	"strings"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func (record *Session_Store) prepare() (Session_Store, error) {
	//os.Exit(0)
	//logger.ErrorLogger.Printf("ACT: VAL Validate")
	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	return *record, nil
}

func (record *Session_Store) calculate() error {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	// Calculate the duration in days between the start and end dates
	return nil
}

func (record *Session_Store) isDuplicateOf(id string) (Session_Store, error) {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	//TODO: Could be replaced with a simple read...

	// Get all status
	recordList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("Getting all %v failed %v", domain, err.Error())
		return Session_Store{}, err
	}

	// range through status list, if status code is found and deletedby is empty then return error
	for _, checkRecord := range recordList {
		//s.Dump(!,strings.ToUpper(code) + "-uchk-" + s.Code)
		testValue := strings.ToUpper(id)
		checkValue := strings.ToUpper(checkRecord.Key)
		//logger.InfoLogger.Printf("CHK: TestValue:[%v] CheckValue:[%v]", testValue, checkValue)
		//logger.InfoLogger.Printf("CHK: Code:[%v] s.Code:[%v] s.Audit.DeletedBy:[%v]", testCode, s.Code, s.Audit.DeletedBy)
		if checkValue == testValue && checkRecord.Audit.DeletedBy == "" {
			logHandler.WarningLogger.Printf("Duplicate %v, %v already in use", strings.ToUpper(domain), record.ID)
			return checkRecord, commonErrors.ErrorDuplicate
		}
	}

	return Session_Store{}, nil
}

func ClearDown(ctx context.Context) error {
	logHandler.EventLogger.Printf("Clearing %v", domain)

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.CLEAR.GetCode(), "INITIALISE")

	// Delete all active session recordList
	recordList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, err).Error())
		clock.Stop(0)
		return commonErrors.WrapDAOInitialisationError(domain, err)
	}

	noRecords := len(recordList)
	count := 0

	for thisRecord, record := range recordList {
		logHandler.InfoLogger.Printf("Deleting %v (%v/%v) %v", domain, thisRecord, noRecords, record.Key)
		delErr := Delete(ctx, record.ID, fmt.Sprintf("Clearing %v %v @ initialisation ", domain, record.ID))
		if delErr != nil {
			logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, delErr).Error())
			continue
		}
		count++
	}

	clock.Stop(count)

	return nil
}

func GetAllByUserID(userID int) []Session_Store {

	rtnList, err := GetAllWhere(FIELD_UserID, userID)
	if err != nil {
		logHandler.ErrorLogger.Printf("ERROR Getting all status: %v", err)
	}

	return rtnList
}

func GetAllByUserCode(userCode string) []Session_Store {
	rtnList, err := GetAllWhere(FIELD_UserCode, userCode)
	if err != nil {
		logHandler.ErrorLogger.Printf("ERROR Getting all status: %v", err)
	}
	return rtnList
}
