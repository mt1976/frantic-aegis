package sessionStore

// Data Access Object Session
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/dao/lookup"
	"github.com/mt1976/frantic-core/ioHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/paths"
	"github.com/mt1976/frantic-core/timing"
)

func Count() (int, error) {
	logHandler.DatabaseLogger.Printf("Count %v", domain)
	return activeDB.Count(&Session_Store{})
}

func CountWhere(field string, value any) (int, error) {
	logHandler.DatabaseLogger.Printf("Count %v where (%v=%v)", domain, field, value)
	clock := timing.Start(domain, actions.COUNT.GetCode(), fmt.Sprintf("%v=%v", field, value))
	list, err := GetAllWhere(field, value)
	if err != nil {
		clock.Stop(0)
		return 0, err
	}
	clock.Stop(len(list))
	return len(list), nil
}

func GetById(id any) (Session_Store, error) {
	return GetBy(FIELD_ID, id)
}

func GetByKey(key any) (Session_Store, error) {
	return GetBy(FIELD_SessionID, key)
}

func GetBy(field string, value any) (Session_Store, error) {

	clock := timing.Start(domain, actions.GET.GetCode(), fmt.Sprintf("%v=%v", field, value))

	dao.CheckDAOReadyState(domain, audit.GET, initialised) // Check the DAO has been initialised, Mandatory.

	if err := dao.IsValidFieldInStruct(field, Session_Store{}); err != nil {
		return Session_Store{}, err
	}

	record := Session_Store{}
	logHandler.DatabaseLogger.Printf("Get %v where (%v=%v)", domain, field, value)

	if err := activeDB.Retrieve(field, value, &record); err != nil {
		clock.Stop(0)
		return Session_Store{}, commonErrors.WrapRecordNotFoundError(domain, field, fmt.Sprintf("%v", value))
	}

	if err := record.PostGet(); err != nil {
		clock.Stop(0)
		return Session_Store{}, commonErrors.WrapDAOReadError(domain, field, value, err)
	}

	clock.Stop(1)
	return record, nil
}

func GetAll() ([]Session_Store, error) {

	dao.CheckDAOReadyState(domain, audit.GET, initialised) // Check the DAO has been initialised, Mandatory.

	recordList := []Session_Store{}

	clock := timing.Start(domain, actions.GETALL.GetCode(), "ALL")

	if errG := activeDB.GetAll(&recordList); errG != nil {
		clock.Stop(0)
		return []Session_Store{}, commonErrors.WrapNotFoundError(domain, errG)
	}

	if _, errPost := PostGet(&recordList); errPost != nil {
		clock.Stop(0)
		return nil, errPost
	}

	clock.Stop(len(recordList))

	return recordList, nil
}

func GetAllWhere(field string, value any) ([]Session_Store, error) {

	dao.CheckDAOReadyState(domain, audit.GET, initialised) // Check the DAO has been initialised, Mandatory.

	recordList := []Session_Store{}
	resultList := []Session_Store{}

	clock := timing.Start(domain, actions.GETALL.GetCode(), fmt.Sprintf("%v=%v", field, value))

	if err := dao.IsValidFieldInStruct(field, Session_Store{}); err != nil {
		return recordList, err
	}

	recordList, err := GetAll()
	if err != nil {
		return []Session_Store{}, err
	}
	count := 0

	for _, record := range recordList {
		if reflect.ValueOf(record).FieldByName(field).Interface() == value {
			count++
			resultList = append(resultList, record)
		}
	}

	if _, errPost := PostGet(&resultList); errPost != nil {
		clock.Stop(0)
		return nil, errPost
	}

	clock.Stop(len(resultList))

	return resultList, nil
}

func (record *Session_Store) Update(ctx context.Context, note string) error {

	dao.CheckDAOReadyState(domain, audit.UPDATE, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.UPDATE.GetCode(), fmt.Sprintf("%v", record.ID))

	if err := record.Validate(); err != nil {
		clock.Stop(0)
		return err
	}

	if calculationError := record.calculate(); calculationError != nil {
		rtnErr := commonErrors.WrapDAOCaclulationError(domain, calculationError)
		logHandler.ErrorLogger.Print(rtnErr.Error())
		clock.Stop(0)
		return rtnErr
	}

	if _, validationError := record.prepare(); validationError != nil {
		valErr := commonErrors.WrapDAOValidationError(domain, validationError)
		logHandler.ErrorLogger.Print(valErr.Error())
		clock.Stop(0)
		return valErr
	}

	auditErr := record.Audit.Action(ctx, audit.UPDATE.WithMessage(note))
	if auditErr != nil {
		audErr := commonErrors.WrapDAOUpdateAuditError(domain, record.ID, auditErr)
		logHandler.ErrorLogger.Print(audErr.Error())
		clock.Stop(0)
		return audErr
	}

	if err := activeDB.Update(record); err != nil {
		updErr := commonErrors.WrapDAOUpdateError(domain, err)
		logHandler.ErrorLogger.Panic(updErr.Error(), err)
		clock.Stop(0)
		return updErr
	}

	//logHandler.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.UPDATE, strings.ToUpper(domain), record.ID, note)
	clock.Stop(1)

	return nil
}

func Delete(ctx context.Context, id int, note string) error {
	return DeleteBy(ctx, FIELD_ID, id, note)
}

func DeleteByKey(ctx context.Context, key string, note string) error {
	return DeleteBy(ctx, FIELD_SessionID, key, note)
}
func DeleteBy(ctx context.Context, field string, value any, note string) error {

	dao.CheckDAOReadyState(domain, audit.DELETE, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.DELETE.GetCode(), fmt.Sprintf("%v=%v", field, value))

	if err := dao.IsValidFieldInStruct(field, Session_Store{}); err != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAODeleteError(domain, field, value, err).Error())
		clock.Stop(0)
		return commonErrors.WrapDAODeleteError(domain, field, value, err)
	}

	record, err := GetBy(field, value)

	if err != nil {
		getErr := commonErrors.WrapDAODeleteError(domain, field, value, err)
		logHandler.ErrorLogger.Panic(getErr.Error(), err)
		clock.Stop(0)
		return getErr
	}

	auditErr := record.Audit.Action(ctx, audit.DELETE.WithMessage(note))
	if auditErr != nil {
		audErr := commonErrors.WrapDAOUpdateAuditError(domain, value, auditErr)
		logHandler.ErrorLogger.Print(audErr.Error())
		clock.Stop(0)
		return audErr
	}

	record.Export(audit.DELETE.Description())

	if err := activeDB.Delete(&record); err != nil {
		delErr := commonErrors.WrapDAODeleteError(domain, field, value, err)
		logHandler.ErrorLogger.Panic(delErr.Error())
		clock.Stop(0)
		return delErr
	}

	//logHandler.AuditLogger.Printf("%v %v (%v=%v) %v", audit.DELETE.Description(), domain, field, value, note)

	clock.Stop(1)

	return nil
}

func (record *Session_Store) Spew() {
	logHandler.InfoLogger.Printf(" [%v] Record=[%+v]", strings.ToUpper(domain), record)
}

func (record *Session_Store) Validate() error {
	return nil
}

func PostGet(recordList *[]Session_Store) ([]Session_Store, error) {
	clock := timing.Start(domain, actions.PROCESS.GetCode(), "POSTGET")
	returnList := []Session_Store{}
	for _, record := range *recordList {
		if err := record.PostGet(); err != nil {
			return nil, err
		}
		returnList = append(returnList, record)
	}
	clock.Stop(len(returnList))
	return returnList, nil
}

func (s *Session_Store) PostGet() error {
	clock := timing.Start(domain, actions.PROCESS.GetCode(), fmt.Sprintf("%v", s.ID))
	clock.Stop(1)
	return nil
}

func Export(message string) {

	dao.CheckDAOReadyState(domain, audit.EXPORT, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.EXPORT.GetCode(), "ALL")
	recordList, _ := GetAll()
	if len(recordList) == 0 {
		logHandler.WarningLogger.Printf("[%v] %v data not found", strings.ToUpper(domain), domain)
		clock.Stop(0)
		return
	}
	SEP := "!"
	for _, record := range recordList {
		msg := fmt.Sprintf("%v%v%v", audit.EXPORT.Description(), SEP, message)
		if message == "" {
			msg = fmt.Sprintf("%v%v", audit.EXPORT.Description(), SEP)
		}
		record.Export(msg)
	}
	clock.Stop(len(recordList))
}

func (record *Session_Store) Export(name string) {

	ID := reflect.ValueOf(*record).FieldByName(FIELD_ID)

	clock := timing.Start(domain, actions.EXPORT.GetCode(), fmt.Sprintf("%v", ID))

	ioHelpers.Dump(domain, paths.Dumps(), name, fmt.Sprintf("%v", ID), record)

	clock.Stop(1)
}

func GetDefaultLookup() (lookup.Lookup, error) {
	return GetLookup(FIELD_SessionID, FIELD_Raw)
}

func GetLookup(field, value string) (lookup.Lookup, error) {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.LOOKUP.GetCode(), "BUILD")

	// Get all status
	recordList, err := GetAll()
	if err != nil {
		lkpErr := commonErrors.WrapDAOLookupError(domain, field, value, err)
		logHandler.ErrorLogger.Print(lkpErr.Error())
		clock.Stop(0)
		return lookup.Lookup{}, lkpErr
	}

	// Create a new Lookup
	var rtnLookup lookup.Lookup
	rtnLookup.Data = make([]lookup.LookupData, 0)

	// range through Behaviour list, if status code is found and deletedby is empty then return error
	for _, a := range recordList {
		key := reflect.ValueOf(a).FieldByName(field).Interface().(string)
		val := reflect.ValueOf(a).FieldByName(value).Interface().(string)
		rtnLookup.Data = append(rtnLookup.Data, lookup.LookupData{Key: key, Value: val})
	}

	clock.Stop(len(rtnLookup.Data))

	return rtnLookup, nil
}

func Drop() error {
	return activeDB.Drop(Session_Store{})
}

// FetchDatabaseInstances returns the database connection
func FetchDatabaseInstances() func() ([]*database.DB, error) {
	//logHandler.InfoLogger.Println("GETDB")
	return func() ([]*database.DB, error) {
		return []*database.DB{activeDB}, nil
	}
}
