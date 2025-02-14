package sessionStore

import (
	"fmt"
	"strings"
	"time"

	appError "github.com/mt1976/frantic-core/commonErrors"
	audit "github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/database"
	lookup "github.com/mt1976/frantic-core/dao/lookup"
	"github.com/mt1976/frantic-core/id"
	"github.com/mt1976/frantic-core/logger"
	stopwatch "github.com/mt1976/frantic-core/timing"
)

func (u *SessionStore) prepare() (SessionStore, error) {
	//os.Exit(0)
	//logger.ErrorLogger.Printf("ACT: VAL Validate")

	return *u, nil
}

func (u *SessionStore) calculate() error {
	// Calculate the duration in days between the start and end dates
	return nil
}

func (u *SessionStore) isDuplicateOf(id string) (SessionStore, error) {

	//logger.InfoLogger.Printf("CHK: CheckUniqueCode %v", name)

	//TODO: Could be replaced with a simple read...

	// Get all status
	activityList, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return SessionStore{}, err
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

	return SessionStore{}, nil
}

func BuildLookup() (lookup.Lookup, error) {

	//logger.InfoLogger.Printf("BuildLookup")

	// Get all status
	Activities, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return lookup.Lookup{}, err
	}

	// Create a new Lookup
	var rtnList lookup.Lookup
	rtnList.Data = make([]lookup.LookupData, 0)

	// range through Behaviour list, if status code is found and deletedby is empty then return error
	for _, a := range Activities {
		rtnList.Data = append(rtnList.Data, lookup.LookupData{Key: a.ID, Value: a.ID})
	}

	return rtnList, nil
}

func GetByUserID(userID int) []SessionStore {
	var rtnList []SessionStore
	// Get all status
	activityList, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return rtnList
	}

	// range through status list, if status code is found and deletedby is empty then return error
	for _, a := range activityList {
		if a.UserID == userID {
			rtnList = append(rtnList, a)
		}
	}

	return rtnList
}

func New(userID int) (SessionStore, error) {
	// Create a new struct
	u := SessionStore{UserID: userID}
	sessionID := id.GetUUID()
	u.Raw = sessionID
	u.ID = id.Encode(sessionID)
	u.Expiry = time.Now().Add(time.Minute * time.Duration(sessionExpiry))
	u.Dump("NEW")

	// Record the create action in the audit data
	_ = u.Audit.Action(nil, audit.CREATE.WithMessage(fmt.Sprintf("New Session for: %v", userID)))

	// Save the status instance to the database
	err := database.Create(&u)
	if err != nil {
		// Log and panic if there is an error creating the status instance
		logger.ErrorLogger.Printf("[%v] Creating Session=[%v] %e", strings.ToUpper(domain), u.ID, err)
		panic(err)
	}

	logger.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.CREATE, strings.ToUpper(domain), u.ID, fmt.Sprintf("New Session: %v", userID))

	return u, nil
}

func Initialise() error {

	timeing := stopwatch.Start("SessionTokens", "Initialise", "")

	// Delete all active session tokens
	tokens, err := GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return err
	}

	noTokens := len(tokens)

	for _, t := range tokens {
		_ = DeleteByID(nil, t.ID, "Initialise")
	}
	timeing.Stop(noTokens)
	return nil
}
