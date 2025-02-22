package sessionStore

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/paths"
)

// SessionImportData is a struct to hold the data from the CSV file
// it is used to import the data into the database
// The struct tags are used to map the fields to the CSV columns
// this struct should be customised to suit the specific requirements of the entryination table/DAO.
type SessionImportData struct {
	Original string `csv:"original"`
	Message  string `csv:"message"`
}

var COMMA = '|'

func ExportCSV() error {
	logHandler.ExportLogger.Printf("Exporting %v", domain)

	Initialise(context.TODO())

	exportFile := openImportFile("export", logHandler.ExportLogger)
	defer exportFile.Close()

	exports, err := GetAll()
	if err != nil {
		logHandler.ExportLogger.Panicf("Error Getting all texts: %v", err.Error())
	}

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = COMMA // Use tab-delimited format
		writer.UseCRLF = true
		return gocsv.NewSafeCSVWriter(writer)
	})

	_, err = gocsv.MarshalString(exports) // Get all texts as CSV string
	if err != nil {
		logHandler.ExportLogger.Panicf("Error exporting texts: %v", err.Error())
	}
	err = gocsv.MarshalFile(&exports, exportFile) // Get all texts as CSV string
	if err != nil {
		logHandler.ExportLogger.Panicf("Error exporting texts: %v", err.Error())
	}

	msg := fmt.Sprintf("# Generated (%v) %vs at %v on %v", len(exports), domain, time.Now().Format("15:04:05"), time.Now().Format("2006-01-02"))
	exportFile.WriteString(msg)

	exportFile.Close()

	logHandler.ExportLogger.Printf("Exported (%v) %vs", len(exports), domain)
	return nil
}

func openImportFile(in string, useLog *log.Logger) *os.File {
	defaultPath := paths.Defaults()
	SessionDataFileName := strings.ToLower(domain) + "s.csv"
	fileName := fmt.Sprintf("%s%s/%s", paths.Application().String(), defaultPath, SessionDataFileName)

	// fmt.Printf("exportPath: %v\n", exportPath)
	// fmt.Printf("textsFile: %v\n", textsFileName)

	dataFileHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Panicf("Error opening file: %v", err)
		panic(err)
	}
	//fmt.Printf("textsFile.Name(): %v\n", textsFile.Name())
	useLog.Printf("Import/Export=[%v] File=[%v]", in, dataFileHandle.Name())
	return dataFileHandle
}

func ImportCSV() error {
	logHandler.ImportLogger.Printf("Importing %v", domain)

	Initialise(context.TODO())

	csvFile := openImportFile("import", logHandler.ImportLogger)
	defer csvFile.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)    // Allows use pipe as delimiter
		r.Comma = COMMA           // Use tab-delimited format
		r.Comment = '#'           // Ignore comment lines
		r.TrimLeadingSpace = true // Trim leading space
		return r                  // Allows use pipe as delimiter
	})

	insertEntriesList := []*SessionImportData{}

	if err := gocsv.UnmarshalFile(csvFile, &insertEntriesList); err != nil { // Load clients from file
		logHandler.ImportLogger.Printf("Importing %v: %v - No Content, nothing to import.", domain, err.Error())
		csvFile.Close()
		return nil
	}

	if _, err := csvFile.Seek(0, 0); err != nil { // Go to the start of the file
		logHandler.ImportLogger.Printf("Importing %v: %v", domain, err.Error())
		panic(err)
	}
	totalImportEntries := len(insertEntriesList)
	for thisPos, insertEntry := range insertEntriesList {
		//fmt.Printf("% 3v)[%v][%v][%v]\n", i, textEntry.Original, textEntry.Message, textEntry)
		logHandler.ImportLogger.Printf("Importing %v (%v/%v) [%v]", domain, thisPos+1, totalImportEntries, insertEntry.Message)
		// Check if the entry already exists, this is a simple check to avoid duplicates
		// it may be necessary to check for duplicates in a more sophisticated way
		// depending on the requirements of the entryination table/DAO.
		existingEntry, _ := GetBy(FIELD_SessionID, idHelpers.Encode(insertEntry.Original))
		if existingEntry.Key != "" {
			//logger.InfoLogger.Printf("Text already exists: [%v]", textEntry.Message)
			continue
		}
		// the load function is a helper function to create a new entry instance and save it to the database
		// the parameters should be customised to suit the specific requirements of the entryination table/DAO.
		_, err := load(insertEntry.Original, insertEntry.Message)
		if err != nil {
			logHandler.ImportLogger.Panicf("Error importing %v: %v", domain, err.Error())
		}
		logHandler.ImportLogger.Printf("Imported %v [%v] [%v]", domain, insertEntry.Original, insertEntry.Message)
	}

	logHandler.ImportLogger.Printf("Imported (%v) %v", len(insertEntriesList), domain)
	csvFile.Close()
	return nil
}

// load is a helper function to create a new entry instance and save it to the database
// It should be customised to suit the specific requirements of the entryination table/DAO.
func load(original, message string) (Session_Store, error) {

	//logger.InfoLogger.Printf("ACT: NEW New %v %v %v", tableName, name, entryination)
	u := Session_Store{}
	u.Key = idHelpers.Encode(strings.ToUpper(original))
	// u.Message = message
	// u.Original = message
	// Add basic attributes

	// Record the create action in the audit data
	_ = u.Audit.Action(context.TODO(), audit.IMPORT.WithMessage(fmt.Sprintf("Imported text [%v]", message)))

	dupe, err := u.isDuplicateOf(u.Key)
	// Log the entry instance before the creation
	if u.Validate() == commonErrors.ErrorDuplicate {
		// This is OK, do nothing as this is a duplicate record
		// we ignore duplicate entryinations.
		logHandler.ImportLogger.Printf("DUPLICATE of %v available in use as [%v]", message, u.Key)
		return dupe, nil
	}

	if err != nil {
		logHandler.ImportLogger.Panicf("Error=[%s]", err.Error())
		return Session_Store{}, err
	}

	// Save the entry instance to the database
	if u.Key == "" {
		logHandler.ImportLogger.Printf("[%v] ID is required, skipping", strings.ToUpper(domain))
		return Session_Store{}, nil
	}

	err = activeDB.Create(&u)
	if err != nil {
		// Log and panic if there is an error creating the entry instance
		logHandler.ImportLogger.Panicf("[%v] Create %s", strings.ToUpper(domain), err.Error())
		panic(err)
	}

	msg := fmt.Sprintf("Imported %v available Id=[%v] Message=[%v]", domain, original, message)
	logHandler.ImportLogger.Println(msg)
	// Return the created entry and nil error
	return u, nil
}
