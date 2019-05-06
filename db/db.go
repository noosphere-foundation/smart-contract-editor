package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"noosphere.foundation/smart-contract-editor/config"
	"noosphere.foundation/smart-contract-editor/mailing"
	"noosphere.foundation/smart-contract-editor/utils"
)

const (
	PATH_TO_EMAIL_CONFIRMED_TEMPLATE = "www/html/email-confirmed-template.html"
	CREATE_TABLE_USERS               = `CREATE TABLE users (firstname TEXT, lastname TEXT, email TEXT, password TEXT,
		hash TEXT, active INTEGER, is_logged_in INTEGER, session_id TEXT, PRIMARY KEY (email, session_id));`
	CREATE_TABLE_SMART_CONTRACTS = `CREATE TABLE smart_contracts (id TEXT PRIMARY KEY, status TEXT, type TEXT,
		creation_date TEXT, price TEXT, last_started TEXT, comment TEXT, code TEXT, data BLOB, was_started INTEGER);`
	CREATE_TABLE_SMART_CONTRACT_BINDS = `CREATE TABLE binds (email TEXT, smart_contract_id TEXT PRIMARY KEY);`
	MAIN_DB_FILENAME                  = "users.db"
	// MAIN_DB_NAME = "file:users.db?mode=memory"
	DNS = "users.db?mode=memory"

	SMART_CONTRACT_ID_MAX_VALUE = 100000

	SMART_CONTRACT_STATUS_DRAFT   = "DRAFT"
	SMART_CONTRACT_STATUS_READY   = "READY TO START"
	SMART_CONTRACT_STATUS_WORKING = "WORKING"
	SMART_CONTRACT_STATUS_PAUSED  = "PAUSED"
)

var db *sql.DB
var timer = &utils.Timer{}
var SmartContractIDs map[string]bool = make(map[string]bool, 0)

func init() {
	log.Println("initializing database module...")
	// make new seed to unique random generation
	rand.Seed(time.Now().UTC().UnixNano())

	// create users and smart_contracts tables if not exists
	CreateSQLiteTablesIfDoesNotExist()

	err := PopulateSmartContractIDs()
	if err != nil {
		log.Fatal("Error when populating smartContractIDs:", err)
	}
	log.Println("ok")
}

func CreateSQLiteTablesIfDoesNotExist() {
	var err error
	db, err = sql.Open("sqlite3", DNS)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(MAIN_DB_FILENAME); err == nil {
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	// create 'users' table
	stmt, err := db.Prepare(CREATE_TABLE_USERS)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// create 'smart-contracts' table
	stmt, err = db.Prepare(CREATE_TABLE_SMART_CONTRACTS)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// create 'binds' table
	stmt, err = db.Prepare(CREATE_TABLE_SMART_CONTRACT_BINDS)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	FirstName      string
	LastName       string
	EMail          string
	Password       string
	ActivatingHash string
	Active         int8
	IsLoggedIn     int8
	SessionID      string
	Balance        float64
}

func AddNewUserToDB(user User) error {
	stmt, err := db.Prepare("insert into users values(?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Println("SOME SHIT HAPPEND")
		return err
	}

	fmt.Printf("trying insert values('%s', '%s', '%s', '%s', '%s', %d, %d, '%s') into users;\n",
		user.FirstName, user.LastName, user.EMail, user.Password, user.ActivatingHash, user.Active, user.IsLoggedIn,
		user.SessionID)
	_, err = stmt.Exec(user.FirstName, user.LastName, user.EMail, user.Password, user.ActivatingHash,
		user.Active, user.IsLoggedIn, user.SessionID)
	if err != nil {
		log.Println("SOME SHIT HAPPEND PART 2")
		return err
	}
	return nil
}

func CheckIfEMailExists(email string) bool {
	rows, err := db.Query("select email from users where email=?", email)
	if err != nil {
		log.Fatal(err)
	}
	emailsNumber := 0
	var s string
	for rows.Next() {
		emailsNumber++
		rows.Scan(&s)
	}
	return emailsNumber > 0
}

func GetLastRowID() int {
	rows := db.QueryRow("select max(rowid) from users")
	var lastRowID int
	rows.Scan(&lastRowID)
	return lastRowID
}

func ActivateUser(userID, userActivationHash string) (string, bool, error) {
	rows := db.QueryRow(
		`	SELECT firstname, lastname, email, active
			FROM users
			WHERE rowid = ?
			AND hash = ?`, userID, userActivationHash)
	var firstName, lastName, email string
	var active int8
	rows.Scan(&firstName, &lastName, &email, &active)

	if active == 1 {
		return "", false, nil
	}

	_, err := db.Exec(`UPDATE users SET active=? WHERE email=?`, 1, email)
	if err != nil {
		return "", false, err
	}

	// NEED TO REFACTOR!!!
	templateData := struct{ FullName string }{firstName + " " + lastName}
	emailText, err := utils.ParseTemplate(PATH_TO_EMAIL_CONFIRMED_TEMPLATE, templateData)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		timer.Start("sending confirmed email...")
		_, err := mailing.SendEMail(config.Data.EMailBlock.EMailValidation,
			config.Data.EMailBlock.EMailValidationPassword, email, "text/html", "Confirmation complete", emailText)
		if err != nil {
			log.Fatal(err)
		}
		timer.End("time elapsed:")
	}()
	// ###################
	return firstName + " " + lastName, true, nil
}

func Login(email, hashedPassword string) (string, bool, error) {
	rows := db.QueryRow(
		`	SELECT rowid, active
			FROM users
			WHERE email=?
			AND password=?`, email, hashedPassword)
	var rowID string
	var active int
	err := rows.Scan(&rowID, &active)
	if err != nil {
		return "", false, err
	}
	return rowID, active == 1, nil
}

func GetUserFullNameByID(userID string) (string, error) {
	rows := db.QueryRow(
		`	SELECT firstname, lastname
			FROM users
			WHERE rowid=?`, userID)
	var firstname, lastname string
	err := rows.Scan(&firstname, &lastname)
	if err != nil {
		return "", err
	}
	return firstname + " " + lastname, nil
}

func GetEMailBySessionID(sessionID string) (string, error) {
	row := db.QueryRow(
		`	SELECT 	email
			FROM 		users
			WHERE 	session_id = ?`, sessionID)

	var email string
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}

func AddSmartContractToDB(sessionID string, smartContract *utils.SmartContract) error {
	stmt, err := db.Prepare("insert into smart_contracts values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	fmt.Println("1")
	if err != nil {
		return err
	}

	fmt.Printf("trying insert values('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%v', '%d') into smart_contracts;\n",
		smartContract.ID, smartContract.Status, smartContract.Type, smartContract.CreationDate, smartContract.Price,
		smartContract.LastStarted, smartContract.Comment, smartContract.Code, smartContract.Data, 0)

	fmt.Println("2")
	mapBytes, err := utils.MapToBytes(smartContract.Data)
	if err != nil {
		return err
	}

	fmt.Println("3")
	_, err = stmt.Exec(smartContract.ID, smartContract.Status, smartContract.Type, smartContract.CreationDate,
		smartContract.Price, smartContract.LastStarted, smartContract.Comment, smartContract.Code, mapBytes, 0)
	if err != nil {
		return err
	}

	// add smart-contract-ID to smartContractIDs
	SmartContractIDs[smartContract.ID] = true

	fmt.Println("4")
	email, err := GetEMailBySessionID(sessionID)
	if err != nil {
		return err
	}

	// add smart-contract id to 'binds' table
	fmt.Println("5")
	stmt, err = db.Prepare("insert into binds values(?, ?);")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(email, smartContract.ID)
	fmt.Println("6")
	if err != nil {
		return err
	}

	return nil
}

func GetAllSmartContracts(sessionID string) ([]*utils.SmartContract, error) {
	fmt.Println("sessionID=", sessionID, "(inside GetAllSmartContracts func)")
	rows, err := db.Query(
		`	SELECT 	id, status, type, creation_date, price, last_started, comment
			FROM 		smart_contracts
			WHERE		id IN (
				SELECT 	smart_contract_id
				FROM		binds
				WHERE		email = (
					SELECT 	email
					FROM 		users
					WHERE		session_id = ?
				)
			)`, sessionID)
	if err != nil {
		return nil, err
	}

	smartContracts := []*utils.SmartContract{}
	for rows.Next() {
		smartContract := utils.SmartContract{}
		err := rows.Scan(&smartContract.ID, &smartContract.Status, &smartContract.Type, &smartContract.CreationDate,
			&smartContract.Price, &smartContract.LastStarted, &smartContract.Comment)
		if err != nil {
			return nil, err
		}
		smartContracts = append(smartContracts, &smartContract)
	}
	return smartContracts, nil
}

func GetSmartContractByID(id string) (string, error) {
	row := db.QueryRow(
		`	SELECT 	code
			FROM 		smart_contracts
			WHERE 	id = ?`, id)

	var smartContract string
	err := row.Scan(&smartContract)
	if err != nil {
		return "", err
	}

	return smartContract, nil
}

func DeleteSmartContractByID(id string) error {
	// delete from 'smart-contracts' table
	stmt, err := db.Prepare("delete from smart_contracts where id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	// delete from 'binds' table
	stmt, err = db.Prepare("delete from binds where smart_contract_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSmartContractByID(id, newSmartContract, newPrice string) error {
	stmt, err := db.Prepare("update smart_contracts set code = ?, price = ? where id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(newSmartContract, newPrice, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSmartContractDataByID(id string, dataMap map[string]interface{}) error {
	stmt, err := db.Prepare("update smart_contracts set data = ? where id = ?")
	if err != nil {
		return err
	}

	mapBytes, err := utils.MapToBytes(dataMap)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(mapBytes, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSmartContractCommentByID(id, comment string) error {
	stmt, err := db.Prepare("update smart_contracts set comment = ? where id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(comment, id)
	if err != nil {
		return err
	}
	return nil
}

func PopulateSmartContractIDs() error {
	rows, err := db.Query(
		`	SELECT 	id
			FROM 		smart_contracts`)
	if err != nil {
		return err
	}

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		SmartContractIDs[id] = true
	}
	return nil
}

// generate unique random number!
func GenerateUniqueID() string {
	newID := strconv.Itoa(rand.Intn(SMART_CONTRACT_ID_MAX_VALUE))
	_, ok := SmartContractIDs[newID]
	for ok {
		newID = strconv.Itoa(rand.Intn(SMART_CONTRACT_ID_MAX_VALUE))
		_, ok = SmartContractIDs[newID]
	}
	return newID
}

func UpdateStatusOfSmartContract(id, statusNew string) error {
	stmt, err := db.Prepare("update smart_contracts set status = ? where id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(statusNew, id)
	if err != nil {
		return err
	}
	return nil
}

func GetSmartContractPriceByID(id string) (float64, error) {
	row := db.QueryRow(
		`	SELECT 	price
			FROM 		smart_contracts
			WHERE 	id = ?`, id)

	var price string
	var err error
	err = row.Scan(&price)
	if err != nil {
		return -1, err
	}

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return -1, err
	}

	return priceFloat, nil
}

func GetSmartContractTypeByID(id string) (string, error) {
	row := db.QueryRow(
		`	SELECT 	type
			FROM 		smart_contracts
			WHERE 	id = ?`, id)

	var smartContractType string
	err := row.Scan(&smartContractType)
	if err != nil {
		return "", err
	}

	return smartContractType, nil
}

func GetSmartContractDataByID(id string) (map[string]interface{}, error) {
	row := db.QueryRow(
		`	SELECT 	data
			FROM 		smart_contracts
			WHERE 	id = ?`, id)

	var smartContractBytes []byte
	err := row.Scan(&smartContractBytes)
	if err != nil {
		return nil, err
	}

	var smartContractMap map[string]interface{}
	smartContractMap, err = utils.BytesToMap(smartContractBytes)
	if err != nil {
		return nil, err
	}

	return smartContractMap, nil
}

func SetUserIsLoggedInProperty(email string, isLoggedInValue int8) error {
	stmt, err := db.Prepare("update users set is_logged_in = ? where email = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(isLoggedInValue, email)
	if err != nil {
		return err
	}
	return nil
}

func SetUserSessionID(email, sessionID string) error {
	stmt, err := db.Prepare("update users set session_id = ? where email = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(sessionID, email)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSessionID(sessionID string) error {
	email, err := GetEMailBySessionID(sessionID)
	if err != nil {
		return err
	}

	err = SetUserIsLoggedInProperty(email, 0)
	if err != nil {
		return err
	}

	err = SetUserSessionID(email, "")
	if err != nil {
		return err
	}
	return nil
}

func IsUserLoggedIn(sessionID string) (bool, error) {
	row := db.QueryRow(
		`	SELECT 	is_logged_in
			FROM 		users
			WHERE 	session_id = ?`, sessionID)

	var isLoggedIn int8
	err := row.Scan(&isLoggedIn)
	if err != nil {
		return false, err
	}

	return isLoggedIn == 1, nil
}

func GetUserFullNameBySessionID(sessionID string) (string, error) {
	row := db.QueryRow(
		`	SELECT 	firstname, lastname
			FROM 		users
			WHERE 	session_id = ?`, sessionID)

	var firstName, lastName string
	err := row.Scan(&firstName, &lastName)
	if err != nil {
		return "", err
	}

	return firstName + " " + lastName, nil
}

func GetSmartContractStatusByID(id string) (string, error) {
	row := db.QueryRow(
		`	SELECT 	status
			FROM 		smart_contracts
			WHERE 	id = ?`, id)

	var status string
	err := row.Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func GetSmartContractWasStartedFieldByID(id string) (int8, error) {
	row := db.QueryRow(
		`	SELECT 	was_started
			FROM 		smart_contracts
			WHERE 	id = ?`, id)

	var wasStarted int8
	err := row.Scan(&wasStarted)
	if err != nil {
		return -1, err
	}

	return wasStarted, nil
}

func MarkSmartContractWasStartedByID(id string) error {
	stmt, err := db.Prepare("update smart_contracts set was_started = 1 where id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	db.Close()
}
