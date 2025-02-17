package sessionStore

import (
	"context"
	"fmt"
	"strings"
	"time"

	appError "github.com/mt1976/frantic-core/commonErrors"
	audit "github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/database"
	lookup "github.com/mt1976/frantic-core/dao/lookup"
	id "github.com/mt1976/frantic-core/idHelpers"
	logger "github.com/mt1976/frantic-core/logHandler"
	stopwatch "github.com/mt1976/frantic-core/timing"
)

func (u *Aegis_SessionStore) prepare() (Aegis_SessionStore, error) {
	//os.Exit(0)
	//logger.ErrorLogger.Printf("ACT: VAL Validate")

	return *u, nil
}

func (u *Aegis_SessionStore) calculate() error {
	// Calculate the duration in days between the start and end dates
	return nil
}

func (u *Aegis_SessionStore) isDuplicateOf(id string) (Aegis_SessionStore, error) {

	//logger.InfoLogger.Printf("CHK: CheckUniqueCode %v", name)

	//TODO: Could be replaced with a simple read...

	// Get all status
	activityList, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return Aegis_SessionStore{}, err
	}

	// range through status list, if status code is found and deletedby is empty then return error
	for _, a := range activityList {
		//s.Dump(!,strings.ToUpper(code) + "-uchk-" + s.Code)
		testValue := strings.ToUpper(id)
		checkValue := strings.ToUpper(a.ID)
		//logger.InfoLogger.Printf("CHK: TestValue:[%v] CheckValue:[%v]", testValue, checkValue)
		//logger.InfoLogger.Printf("CHK: Code:[%v] s.Code:[%v] s.Audit.DeletedBy:[%v]", testCode, s.Code, s.Audit.DeletedBy)
		if checkValue == testValue && a.Audit.DeletedBy == "" {
			logger.InfoLogger.Printf("[%v] DUPLICATE %v already in use", strings.ToUpper(domain), u.ID)
			return a, appError.ErrorDuplicate
		}
	}

	//logger.InfoLogger.Printf("CHK: %v is unique", strings.ToUpper(name))

	// Return nil if the code is unique

	return Aegis_SessionStore{}, nil
}

func BuildLookup() (lookup.Lookup, error) {
	clock := stopwatch.Start("Sessions", "BuildLookup", "")
	//logger.InfoLogger.Printf("BuildLookup")

	// Get all status
	Activities, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		clock.Stop(0)
		return lookup.Lookup{}, err
	}

	// Create a new Lookup
	var rtnList lookup.Lookup
	rtnList.Data = make([]lookup.LookupData, 0)

	// range through Behaviour list, if status code is found and deletedby is empty then return error
	for _, a := range Activities {
		rtnList.Data = append(rtnList.Data, lookup.LookupData{Key: a.ID, Value: a.ID})
	}
	clock.Stop(len(rtnList.Data))
	return rtnList, nil
}

func GetByUserID(userID int) []Aegis_SessionStore {
	clock := stopwatch.Start("Sessions", "GetByUserID", "")
	var rtnList []Aegis_SessionStore
	// Get all status
	activityList, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		clock.Stop(0)
		return rtnList
	}

	// range through status list, if status code is found and deletedby is empty then return error
	for _, a := range activityList {
		if a.UserID == userID {
			rtnList = append(rtnList, a)
		}
	}

	clock.Stop(len(rtnList))

	return rtnList
}

func New(userID int) (Aegis_SessionStore, error) {

	clock := stopwatch.Start("Sessions", "New", "Create")
	// Create a new struct
	s := Aegis_SessionStore{UserID: userID}
	sessionID := id.GetUUID()
	s.Raw = sessionID
	s.ID = id.Encode(sessionID)
	s.Expiry = time.Now().Add(time.Minute * time.Duration(sessionExpiry))
	s.Dump("NEW" + strings.ToUpper(domain))

	// Record the create action in the audit data
	_ = s.Audit.Action(context.TODO(), audit.CREATE.WithMessage(fmt.Sprintf("New Session for: %v", userID)))

	// Save the status instance to the database
	err := database.Create(&s)
	if err != nil {
		// Log and panic if there is an error creating the status instance
		logger.ErrorLogger.Printf("[%v] Creating Session=[%v] %e", strings.ToUpper(domain), s.ID, err)
		panic(err)
	}

	logger.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.CREATE, strings.ToUpper(domain), s.ID, fmt.Sprintf("New Session: %v", userID))
	clock.Stop(1)
	return s, nil
}

func Initialise() error {

	clock := stopwatch.Start("Sessions", "Initialise", "")

	// Delete all active session tokens
	tokens, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		clock.Stop(0)
		return err
	}

	noTokens := len(tokens)

	for _, t := range tokens {
		_ = DeleteByID(context.TODO(), t.ID, "Initialise")
	}
	clock.Stop(noTokens)
	return nil
}
