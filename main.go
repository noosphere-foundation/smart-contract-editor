package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"noosphere.foundation/smart-contract-editor/config"
	"noosphere.foundation/smart-contract-editor/db"
	"noosphere.foundation/smart-contract-editor/mailing"
	"noosphere.foundation/smart-contract-editor/utils"
)

const (
	PATH_TO_INDEX_PAGE                  = "www/html/index.html"
	PATH_TO_REGISTERED_PAGE             = "www/html/registered.html"
	PATH_TO_EMAIL_VALIDATION_TEMPLATE   = "www/html/email-validation-template.html"
	PATH_TO_SUCCESS_EMAIL_VERIFIED_PAGE = "www/html/email-verified-success.html"
	// PATH_TO_LOGGED_IN_PAGE                 = "www/html/logged-in.html"
	PATH_TO_EMAIL_IS_ALREADY_VERIFIED_PAGE = "www/html/email-is-already-verified.html"

	BIG_INTEGER_NUMBER = 99999999999

	// PVM_ADD_URL = "http://192.168.192.42:25873/pvm/add"
	// PVM_ADD_URL = "http://192.168.192.42:25874/pvm/add"
)

var globalSessionId = "fcetgvs7wie60eyiznv23n25i3wjznkh6hl3iwp7azsb2wcgrwadhcd0d1e4sryr"
var timer *utils.Timer = &utils.Timer{}
var idKeyC string

// path is "/"
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("content-type", "text/html")
		w.Write(config.ReadFile(PATH_TO_INDEX_PAGE))
	} else {
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte("The page isn't found."))
	}
}

// path is "/registered"
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Write(config.ReadFile(PATH_TO_REGISTERED_PAGE))
}

// path is "/check-if-email-exists"
func EMailExistenceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("email received by ajax is '%s'\n", r.FormValue("email"))
	if db.CheckIfEMailExists(r.FormValue("email")) {
		w.Write([]byte("User with such email already exists!"))
	}
}

// path is "/sign-up"
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	timer.Start("parsing html form...")
	if err := r.ParseForm(); err != nil {
		w.Write([]byte("Error when parsing form."))
		log.Fatal(err)
	}
	timer.End("time elapsed:")

	timer.Start("generating activating hash...")
	
	activatingHash, err := utils.GenerateActivatingHash()
	timer.End("time elapsed:")
	if err != nil {
		log.Fatal(err)
	}

	timer.Start("creating user object...")
	user := db.User{FirstName: r.FormValue("firstName"), LastName: r.FormValue("lastName"), EMail: r.FormValue("email"),
		Password: r.FormValue("password"), ActivatingHash: activatingHash, Active: 0, IsLoggedIn: 0, SessionID: ""}
	timer.End("time elapsed:")

	timer.Start("adding new user...")
	err = db.AddNewUserToDB(user)
	timer.End("time elapsed:")

	if err != nil {
		w.Write([]byte("Error when adding new user to DB."))
		log.Fatal(err)
	}

	timer.Start("preparing email confirmation template...")
	mailing.PrepareEMailConfirmationTemplate(PATH_TO_EMAIL_VALIDATION_TEMPLATE, user.FirstName, user.LastName,
		user.EMail, "http://"+config.Data.ServerBlock.URLBase+"/verify?id="+strconv.Itoa(db.GetLastRowID())+
			"&hash="+user.ActivatingHash)
	timer.End("time elapsed:")

	fmt.Printf("total -> adding new user took %s\n", time.Since(start))
}

// path is "/verify"
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	userIDs, ok := r.URL.Query()["id"]
	if !ok || len(userIDs) < 1 {
		log.Fatal("ERROR: missing 'id' parameter in the GET request")
	}
	userID := userIDs[0]

	userActivationHashes, ok := r.URL.Query()["hash"]
	if !ok || len(userActivationHashes) < 1 {
		log.Fatal("ERROR: missing 'hash' parameter in the GET request")
	}
	userActivationHash := userActivationHashes[0]

	fullName, ok, err := db.ActivateUser(userID, userActivationHash)
	if err != nil {
		log.Fatal(err)
	} else if !ok {
		w.Header().Set("content-type", "text/html")
		w.Write(config.ReadFile(PATH_TO_EMAIL_IS_ALREADY_VERIFIED_PAGE))
		return
	}

	templateData := struct{ FullName string }{fullName}
	htmlResponse, err := utils.ParseTemplate(PATH_TO_SUCCESS_EMAIL_VERIFIED_PAGE, templateData)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(htmlResponse))
}

// path is "/login"
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	email, hashedPassword, captchaResult := r.FormValue("email"), r.FormValue("password"), r.FormValue("captcha")
	fmt.Printf("email=%s, password_hash=%s, captcha_result=%s\n", email, hashedPassword, captchaResult)

	// verify captcha
	verifyResult := base64Captcha.VerifyCaptcha(idKeyC, captchaResult)
	fmt.Printf("captchaResult='%s'\n", captchaResult)
	if !verifyResult {
		w.Write([]byte("2"))
		return
	}

	rowID, accountActivated, err := db.Login(email, hashedPassword)
	if err != nil {
		w.Write([]byte("0"))
		return
	}
	fmt.Printf("rowID='%s'\n", rowID)
	if !accountActivated {
		w.Write([]byte("1"))
		return
	}

	// set is_logged_in DB field as 1 (means "user is logged in")
	err = db.SetUserIsLoggedInProperty(email, 1)
	if err != nil {
		log.Fatal(err)
	}

	// generate and set up browser cookie (unique session ID)
	sessionID := utils.GenerateSessionID(64)
	// set cookie
	expiration := time.Now().Add(time.Hour)
	cookie := http.Cookie{Name: utils.SMCE_SESSION_ID, Value: sessionID, Expires: expiration}
	http.SetCookie(w, &cookie)

	// write generated cookie to DB
	err = db.SetUserSessionID(email, sessionID)

	w.Write([]byte("login is ok"))
}

// path is "/generate-captcha"
func GenerateCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	//config struct for Character
	var configC = base64Captcha.ConfigCharacter{
		Height:             60,
		Width:              240,
		Mode:               base64Captcha.CaptchaModeArithmetic,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexHigh,
		IsUseSimpleFont:    true,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}

	//create a characters captcha.
	//GenerateCaptcha first parameter is empty string,so the package will generate a random uuid for you.
	var capC base64Captcha.CaptchaInterface
	idKeyC, capC = base64Captcha.GenerateCaptcha("", configC)
	//write to base64 string.
	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(capC)
	w.Write([]byte(base64stringC))
}

// path is "/contracts"
func ContractsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ContractsHandler is here!!!")
	// sessionID, err := utils.ReadSessionIDFromBrowserCookies(r)
	// if err != nil {
	// 	log.Println(err)
	// }
	//
	// fmt.Printf("sessionID='%s'\n", sessionID)
	//
	// if sessionID == "" {
	// 	fmt.Println("contracts_handler: redirecting to the login page...")
	// http.Redirect(w, r, "/", http.StatusMovedPermanently)
	// }

	w.Header().Set("content-type", "text/html")
	w.Write(config.ReadFile("www/html/contracts.html"))
}

// path is "/generate-contract"
func GenerateContractHandler(w http.ResponseWriter, r *http.Request) {
	//sessionID, err := utils.ReadSessionIDFromBrowserCookies(r)
	sessionID := globalSessionId
	
	var response string
	htmlTemplateData, smartContract, contractInfo, id := utils.BuildTransactionHTMLTemplateData(r)

	// TESTING!!! (start)
	// run python parser and write result to the file
	resultParser := utils.RunPythonParser(smartContract.Code)
	err := utils.WriteFile("python_parser_result.tmp", resultParser)
	if err != nil {
		log.Println("error when writing python parser result:", err)
	}

	if id == "not_exist" {
		smartContract.ID = db.GenerateUniqueID()
		smartContract.Status = db.SMART_CONTRACT_STATUS_READY

		// save smart-contract to db
		err = db.AddSmartContractToDB(sessionID, smartContract)
		if err != nil {
			log.Fatal("Error when adding new smart-contract to DB:", err)
		}

		// generating resulting html response
		response, err = utils.ParseTemplate(contractInfo.TransactionHTMLTemplatePath, htmlTemplateData)
		if err != nil {
			log.Fatalf("Error when parsing '%s' contract template: %s\n", contractInfo.TransactionHTMLTemplatePath,
				err.Error())
		}
	} else {
		// update the smart-contract's code and price
		err = db.UpdateSmartContractByID(id, smartContract.Code,
			utils.CalculatePriceOfSmartContractFake(smartContract.Code))
		if err != nil {
			log.Fatal(err)
		}

		// update the smart-contract's data field (sqlite BLOB type)
		err = db.UpdateSmartContractDataByID(id, smartContract.Data)
		if err != nil {
			log.Fatal(err)
		}

		// update the smart-contract's comment
		err = db.UpdateSmartContractCommentByID(id, smartContract.Comment)
		if err != nil {
			log.Fatal(err)
		}
		response = "ok"
	}

	w.Header().Set("content-type", "text/html")
	w.Write([]byte(response))
}

// path is "/get-all-smart-contracts"
func GetAllSmartContractsHandler(w http.ResponseWriter, r *http.Request) {
	// read session id
// 	sessionID, err := utils.ReadSessionIDFromBrowserCookies(r)
sessionID := globalSessionId
	
	fmt.Println("sessionID=", sessionID)

	smartContracts, err := db.GetAllSmartContracts(sessionID)
	if err != nil {
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte("Error"))
		log.Fatal("Error when getting all smart-contracts:", err)
	}

	var smartContractsBuffer bytes.Buffer
	smartContractsBuffer.WriteString("[")
	for index, smartContract := range smartContracts {
		smartContractsBuffer.WriteString(smartContract.ToJSON())
		if index < len(smartContracts)-1 {
			smartContractsBuffer.WriteString(", ")
		}
	}
	smartContractsBuffer.WriteString("]")

	w.Header().Set("content-type", "text/plain")
	w.Write(smartContractsBuffer.Bytes())
}

// path is "/readSmartContract"
func ReadSmartContractHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	id := r.FormValue("ID")

	var smartContract string
	smartContract, err = db.GetSmartContractByID(id)
	if err != nil {
		log.Fatal(err)
	}

	pythonParserResult := utils.RunPythonParser(smartContract)
	smartContractTransaction := utils.SmartContractTransaction{TT: "SC", TST: utils.GetCurrentTimestamp(),
		CODE: smartContract, ANALYSIS: pythonParserResult}
	smartContractTransactionInJSON, err := smartContractTransaction.ToJSON(true)
	if err != nil {
		log.Fatal(err)
	}

	// read smc status
	smcStatus, err := db.GetSmartContractStatusByID(id)
	if err != nil {
		log.Fatal(err)
	}

	// read smc 'was_started' field
	smcWasStarted, err := db.GetSmartContractWasStartedFieldByID(id)
	if err != nil {
		log.Fatal(err)
	}

	finalStringToSend, err := (&utils.ExplorerTransaction{Transaction: smartContractTransactionInJSON,
		SmartContractStatus: smcStatus, SmartContractWasStarted: smcWasStarted == 1}).ToJSON(true)
	if err != nil {
		log.Fatal(err)
	}

	smartContractTransactionJSONSpaceTrimmed := strings.TrimSpace(finalStringToSend)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(smartContractTransactionJSONSpaceTrimmed))
}

// path is "/deleteSmartContract"
func DeleteSmartContractHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	id := r.FormValue("ID")

	err := db.DeleteSmartContractByID(id)
	w.Header().Set("content-type", "text/plain")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("Smart-contract with ID=%s was deleted successfully.", id)))
}

// path is "/editSmartContract"
func EditSmartContractHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	id := r.FormValue("ID")

	smartContractCode, err := db.GetSmartContractByID(id)
	if err != nil {
		log.Fatal(err)
	}

	smartContractType, err := db.GetSmartContractTypeByID(id)
	if err != nil {
		log.Fatal(err)
	}

	smartContractFieldsMap, err := db.GetSmartContractDataByID(id)
	if err != nil {
		log.Fatal(err)
	}

	smartContract := struct {
		ID   string
		Code string
		Type string
		Data map[string]interface{}
	}{
		ID:   id,
		Code: smartContractCode,
		Type: smartContractType,
		Data: smartContractFieldsMap,
	}
	b, err := json.Marshal(&smartContract)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", "text/plain")
	w.Write(b)
}

// path is "/updateSmartContract"
func UpdateSmartContractHandler(w http.ResponseWriter, r *http.Request) {
// 	sessionID, err := utils.ReadSessionIDFromBrowserCookies(r)
sessionID := globalSessionId
	

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	id, status, updatedSmartContract := r.FormValue("ID"), r.FormValue("Status"), r.FormValue("SmartContractCode")

	if _, ok := db.SmartContractIDs[id]; !ok {
		db.AddSmartContractToDB(sessionID, &utils.SmartContract{ID: id, Status: status, Type: "custom",
			CreationDate: utils.GetCurrentTimestamp(), Price: utils.CalculatePriceOfSmartContractFake(updatedSmartContract),
			LastStarted: "unused", Comment: "no comment", Code: updatedSmartContract})
		return
	}

	err := db.UpdateSmartContractByID(id, updatedSmartContract, utils.CalculatePriceOfSmartContractFake(updatedSmartContract))
	if err != nil {
		log.Fatal(err)
	}
}

// path is "/generate-unique-id"
func GenerateUniqueID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(db.GenerateUniqueID()))
}

// path is "/run-pylint"
func RunPylintHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	smartContractCode := r.FormValue("SmartContractCode")
	utils.WriteFile("smart_contract_code.tmp", smartContractCode)
	analysisCodeResult := utils.RunPylint("smart_contract_code.tmp")
	err = utils.DeleteFile("smart_contract_code.tmp")
	if err != nil {
		log.Fatal(err)
	}

	pylintErrors, err := utils.FindErrorsInThePylintAnalysis(analysisCodeResult)
	if err != nil {
		log.Fatal(err)
	}

	// build result array of jsons
	var pylintErrorsBuffer bytes.Buffer
	pylintErrorsBuffer.WriteString("[")
	for index, pylintError := range pylintErrors {
		pylintErrorsBuffer.WriteString(pylintError.ToJSON())
		if index < len(pylintErrors)-1 {
			pylintErrorsBuffer.WriteString(", ")
		}
	}
	pylintErrorsBuffer.WriteString("]")

	w.Header().Set("content-type", "text/plain")
	w.Write(pylintErrorsBuffer.Bytes())
}

// path is "/updateStatusOfSmartContract"
func UpdateStatusOfSmartContractHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	id, status := r.FormValue("ID"), r.FormValue("Status")

	err := db.UpdateStatusOfSmartContract(id, status)
	if err != nil {
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte("Error"))
		log.Fatal(err)
	}
}

// path is "sendTransactionToPVM"
func SendTransactionToPVM(w http.ResponseWriter, r *http.Request) {
	var resultPVM string
	var err error
	if err = r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	// send to PVM
	transaction := r.FormValue("Transaction")
	resultPVM, err = utils.SendTransaction("pvm", config.Data.ServerBlock.PVMURL, transaction)
	log.Printf("\n\n\ntransaction=[%s]\n\n\n", transaction)
	if err != nil {
		log.Fatal(err)
	}

	id := r.FormValue("ID")
	err = db.MarkSmartContractWasStartedByID(id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("content-type", "text/plain")
	if err != nil {
		w.Write([]byte("error"))
		log.Fatal(err)
	}
	w.Write([]byte("ok"))

	log.Println("resultPVM (transaction): ", resultPVM)
}

// path is "/signout"
func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := utils.ReadSessionIDFromBrowserCookies(r)
	if err != nil {
		log.Println(err)
	}

	err = db.DeleteSessionID(sessionID)
	if err != nil {
		log.Fatal(err)
	}

	// delete cookie from browser
	cookie := http.Cookie{Name: utils.SMCE_SESSION_ID, Value: "", Expires: time.Unix(0, 0)}
	http.SetCookie(w, &cookie)
}

// path is "/logged-as"
func LoggedAsHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := utils.ReadSessionIDFromBrowserCookies(r)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("content-type", "text/plain")
	if sessionID == "" {
		w.Write([]byte("user_is_not_logged"))
		return
	}

	fullname, err := db.GetUserFullNameBySessionID(sessionID)
	if err != nil {
		log.Println(err)
	}

	w.Write([]byte("(Logged as " + fullname + ".)"))
}

// path is "/sendActionTOPVM"
func SendActionTOPVM(w http.ResponseWriter, r *http.Request) {
	var resultPVM string
	var err error
	if err = r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	// check 'started' field of smart-contract
	actionJSON := r.FormValue("ActionJSON")
	log.Printf("\n\n\nActionJSON=[%s]\n\n\n", actionJSON)
	resultPVM, err = utils.SendAction(config.Data.ServerBlock.PVMAction, actionJSON)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("content-type", "text/plain")
	if err != nil {
		w.Write([]byte("error"))
		log.Fatal(err)
	}
	w.Write([]byte("ok"))

	log.Println("resultPVM (action): ", resultPVM)
}

func makeFakeUser() {
    //
    
    log.Println("Making fake user")
    
    if db.CheckIfEMailExists("test@mail.ru") {
		return
	}
	
	activatingHash, _ := utils.GenerateActivatingHash()
    
    user := db.User{FirstName: "User", LastName: "test", EMail: "test@mail.ru",
		Password: "passw0rd", ActivatingHash: activatingHash, Active: 1, IsLoggedIn: 1, SessionID: globalSessionId}
	

	db.AddNewUserToDB(user)
}

func main() {
	log.Println("main func is started")
	defer db.Close()
    makeFakeUser()
	log.Println("preparing to serving static files...")
	// serve static files
	fs := http.FileServer(http.Dir("www/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("ok")

	// the main handlers
	log.Println("initializing http handlers...")
	http.HandleFunc("/sendActionToPVM", SendActionTOPVM)
// 	http.HandleFunc("/logged-as", LoggedAsHandler)
// 	http.HandleFunc("/signout", SignoutHandler)
	http.HandleFunc("/sendTransactionToPVM", SendTransactionToPVM)
	http.HandleFunc("/updateStatusOfSmartContract", UpdateStatusOfSmartContractHandler)
	http.HandleFunc("/run-pylint", RunPylintHandler)
	http.HandleFunc("/generate-unique-id", GenerateUniqueID)
	http.HandleFunc("/updateSmartContract", UpdateSmartContractHandler)
	http.HandleFunc("/editSmartContract", EditSmartContractHandler)
	http.HandleFunc("/deleteSmartContract", DeleteSmartContractHandler)
	http.HandleFunc("/readSmartContract", ReadSmartContractHandler)
	http.HandleFunc("/get-all-smart-contracts", GetAllSmartContractsHandler)
	http.HandleFunc("/generate-contract", GenerateContractHandler)
	http.HandleFunc("/contracts", ContractsHandler)
// 	http.HandleFunc("/generate-captcha", GenerateCaptchaHandler)
// 	http.HandleFunc("/login", LoginHandler)
// 	http.HandleFunc("/verify", VerifyHandler)
// 	http.HandleFunc("/sign-up", SignUpHandler)
// 	http.HandleFunc("/check-if-email-exists", EMailExistenceHandler)
// 	http.HandleFunc("/registered", RegistrationHandler)
// 	http.HandleFunc("/", IndexHandler)
    //http.HandleFunc("/", ContractsHandler)
	log.Println("ok")

	log.Printf("listening on port %s...\n", config.Data.ServerBlock.Port)
	log.Fatal(http.ListenAndServe(config.Data.ServerBlock.ServerAddress, nil))
}
