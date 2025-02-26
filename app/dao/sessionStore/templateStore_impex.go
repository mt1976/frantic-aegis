package sessionStore

// Data Access Object Session
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: Implement the importProcessor function to process the domain entity

import (
	"context"

	"github.com/mt1976/frantic-core/logHandler"
)

// importProcessor is a helper function to create a new entry instance and save it to the database
// It should be customised to suit the specific requirements of the entryination table/DAO.
func importProcessor(inOriginal **Session_Store) (string, error) {

	importedData := **inOriginal

	_, err := New(context.TODO(), importedData.UserKey, importedData.UserCode)
	if err != nil {
		logHandler.ImportLogger.Panicf("Error importing %v: %v", domain, err.Error())
		return importedData.UserCode, err
	}

	return importedData.UserCode, nil
}
