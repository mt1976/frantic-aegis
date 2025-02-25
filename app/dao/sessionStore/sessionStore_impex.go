package sessionStore

import (
	"context"
	"strconv"

	"github.com/mt1976/frantic-core/logHandler"
)

// importProcessor is a helper function to create a new entry instance and save it to the database
// It should be customised to suit the specific requirements of the entryination table/DAO.
func importProcessor(inOriginal **Session_Store) (string, error) {

	importedData := **inOriginal

	//	logHandler.ImportLogger.Printf("Importing %v [%v] [%v]", domain, original.Raw, original.Field1)

	//logger.InfoLogger.Printf("ACT: NEW New %v %v %v", tableName, name, entryination)
	// u := Behaviour_Store{}
	// u.Key = idHelpers.Encode(strings.ToUpper(original.Raw))
	// u.Raw = original.Raw
	// u.Text = original.Text
	// // u.Action = original.Action
	// u.Domain = original.Domain

	// importAction := actions.New(original.Action.Name)
	// bh, err := Declare(importAction, domains.Special(original.Domain), original.Text)
	// if err != nil {
	// 	logHandler.ImportLogger.Panicf("Error importing template: %v", err.Error())
	// }

	// Return the created entry and nil error
	//logHandler.ImportLogger.Printf("Imported %v [%+v]", domain, importedData)

	stringField1 := strconv.Itoa(importedData.ID)

	_, err := New(context.TODO(), importedData.Key, importedData.UserCode)
	if err != nil {
		logHandler.ImportLogger.Panicf("Error importing %v: %v", domain, err.Error())
		return stringField1, err
	}

	return stringField1, nil
}
