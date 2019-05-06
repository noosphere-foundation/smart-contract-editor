package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	// "html/template"
	"io/ioutil"
	"log"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

const (
	BIG_INTEGER_NUMBER = 99999999999

	DEFERRED_CONTRACT_TYPE   = "deferred"
	CONDITION_CONTRACT_TYPE  = "condition"
	AUTO_CONTRACT_TYPE       = "auto"
	COLLECTIVE_CONTRACT_TYPE = "collective"

	PATH_TO_DEFERRED_TRANSACTION_SCRIPT_TEMPLATE = "www/contract-templates/deferred.py"
	PATH_TO_DEFERRED_TRANSACTION_HTML_TEMPLATE   = "www/contract-templates/deferred.html"

	PATH_TO_CONDITION_TRANSACTION_SCRIPT_TEMPLATE = "www/contract-templates/condition.py"
	PATH_TO_CONDITION_TRANSACTION_HTML_TEMPLATE   = "www/contract-templates/condition.html"

	PATH_TO_AUTO_TRANSACTION_SCRIPT_TEMPLATE = "www/contract-templates/auto.py"
	PATH_TO_AUTO_TRANSACTION_HTML_TEMPLATE   = "www/contract-templates/auto.html"

	PATH_TO_COLLECTIVE_TRANSACTION_SCRIPT_TEMPLATE = "www/contract-templates/collective.py"
	PATH_TO_COLLECTIVE_TRANSACTION_HTML_TEMPLATE   = "www/contract-templates/collective.html"

	DEFERRED_TRANSACTION_SCRIPT_TEMP_FILENAME   = "deferred_payment_script.tmp"
	CONDITION_TRANSACTION_SCRIPT_TEMP_FILENAME  = "condition_payment_script.tmp"
	AUTO_TRANSACTION_SCRIPT_TEMP_FILENAME       = "auto_payment_script.tmp"
	COLLECTIVE_TRANSACTION_SCRIPT_TEMP_FILENAME = "collective_payment_script.tmp"

	SMART_CONTRACT_PRICE_COEFFICIENT_NZT float64 = 0.01

	SMCE_SESSION_ID = "smce_session_id"

	START_ACTION = "START"
	PAUSE_ACTION = "PAUSE"
)

// data: {ContractType: "deferred", ID: id, ContractDate: formatted, Receiver: $("#inputReceiver").val(),
// 	Data: $("#inputText").val(), TransactionMessage: $("#inputTransaction").val(),}

var deferredTransactionTemplateFieldNames = []string{"ContractDate", "ID", "Receiver", "Data", "TransactionMessage"}
var conditionTransactionTemplateFieldNames = []string{"ContractDate", "ID", "SelectCondition", "SelectOperator",
	"SelectValue", "Receiver", "Data", "TransactionMessage"}
var autoTransactionTemplateFieldNames = []string{"ContractDate", "ID", "AutoPaymentMode", "Receiver", "Data",
	"TransactionMessage"}
var collectiveTransactionTemplateFieldNames = []string{"ContractDate", "ID", "Receivers", "Data", "TransactionMessage"}

var ContractInfoByName map[string]*ContractInfo = make(map[string]*ContractInfo, 0)

type ContractInfo struct {
	TransactionScriptTemplatePath string
	TemplateFieldNames            []string
	TransactionScriptTempFilename string
	TransactionHTMLTemplatePath   string
}

type Timer struct {
	StartTime    time.Time
	DurationTime time.Duration
}

func init() {
	log.Println("preparing utils.go...")
	ContractInfoByName[DEFERRED_CONTRACT_TYPE] = &ContractInfo{
		TransactionScriptTemplatePath: PATH_TO_DEFERRED_TRANSACTION_SCRIPT_TEMPLATE,
		TemplateFieldNames:            deferredTransactionTemplateFieldNames,
		TransactionScriptTempFilename: DEFERRED_TRANSACTION_SCRIPT_TEMP_FILENAME,
		TransactionHTMLTemplatePath:   PATH_TO_DEFERRED_TRANSACTION_HTML_TEMPLATE,
	}

	ContractInfoByName[CONDITION_CONTRACT_TYPE] = &ContractInfo{
		TransactionScriptTemplatePath: PATH_TO_CONDITION_TRANSACTION_SCRIPT_TEMPLATE,
		TemplateFieldNames:            conditionTransactionTemplateFieldNames,
		TransactionScriptTempFilename: CONDITION_TRANSACTION_SCRIPT_TEMP_FILENAME,
		TransactionHTMLTemplatePath:   PATH_TO_CONDITION_TRANSACTION_HTML_TEMPLATE,
	}

	ContractInfoByName[AUTO_CONTRACT_TYPE] = &ContractInfo{
		TransactionScriptTemplatePath: PATH_TO_AUTO_TRANSACTION_SCRIPT_TEMPLATE,
		TemplateFieldNames:            autoTransactionTemplateFieldNames,
		TransactionScriptTempFilename: AUTO_TRANSACTION_SCRIPT_TEMP_FILENAME,
		TransactionHTMLTemplatePath:   PATH_TO_AUTO_TRANSACTION_HTML_TEMPLATE,
	}

	ContractInfoByName[COLLECTIVE_CONTRACT_TYPE] = &ContractInfo{
		TransactionScriptTemplatePath: PATH_TO_COLLECTIVE_TRANSACTION_SCRIPT_TEMPLATE,
		TemplateFieldNames:            collectiveTransactionTemplateFieldNames,
		TransactionScriptTempFilename: COLLECTIVE_TRANSACTION_SCRIPT_TEMP_FILENAME,
		TransactionHTMLTemplatePath:   PATH_TO_COLLECTIVE_TRANSACTION_HTML_TEMPLATE,
	}
	log.Println("ok")
}

func (t *Timer) Start(message string) {
	if len(message) > 0 {
		fmt.Println(message)
	}
	t.StartTime = time.Now()
}

func (t *Timer) End(message string) {
	t.DurationTime = time.Since(t.StartTime)
	if len(message) > 0 {
		fmt.Println(message + " " + t.DurationTime.String())
	} else {
		fmt.Println("time elapsed:", t.DurationTime)
	}
	fmt.Println()
}

func GetImageExtension(filepath string) string {
	lastDotIndex := strings.LastIndex(filepath, ".")
	ext := "unknown"
	if lastDotIndex+1 < len(filepath) {
		ext = filepath[lastDotIndex+1:]
	}
	if ext == "ico" {
		ext = "x-icon"
	}
	return ext
}

func GetStringOfCryptoRandomInteger() (string, error) {
	max := *big.NewInt(BIG_INTEGER_NUMBER)
	num, err := rand.Int(rand.Reader, &max)
	if err != nil {
		return "", err
	}
	return num.String(), nil
}

func GenerateActivatingHash() (string, error) {
	stringOfRandomInt, err := GetStringOfCryptoRandomInteger()
	fmt.Printf("generated random number: '%s'\n", stringOfRandomInt)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write([]byte(stringOfRandomInt))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ProduceStringWithDashes(numberOfDashes int) string {
	var dashes bytes.Buffer
	for i := 0; i < numberOfDashes; i++ {
		dashes.WriteString("-")
	}
	return dashes.String()
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RunPylint(filepath string) string {
	cmd := exec.Command("pylint", filepath)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		os.Stderr.WriteString("pylint error message: " + err.Error() + "\n")
	}
	return string(cmdOutput.Bytes())
}

func WriteFile(filepath, data string) error {
	err := ioutil.WriteFile(filepath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	return nil
}

// the right time format is "2014-11-12T11:45:26.371Z"
func StringToTime(timeString string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeString)

	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func BuildTransactionScriptTemplateData(r *http.Request) (map[string]interface{}, *ContractInfo, string, error) {
	var err error
	if err = r.ParseForm(); err != nil {
		return nil, nil, "", err
	}

	contractInfo := ContractInfoByName[r.FormValue("ContractType")]
	postValues := make([]string, len(contractInfo.TemplateFieldNames))
	templateData := make(map[string]interface{}, 0)
	var id string
	for index, templateFieldName := range contractInfo.TemplateFieldNames {
		if templateFieldName == "ID" {
			id = r.FormValue("ID")
		} else {
			// check contract date
			if templateFieldName == "ContractDate" {
				timeString := r.FormValue(templateFieldName)
				timeString = strings.Replace(timeString, " ", "T", 1)
				timeString += ".000Z"
				contractDateTime, err := StringToTime(timeString)
				if err != nil {
					log.Fatal("Error when converting string to time:", err)
				}
				if contractDateTime.Before(time.Now()) {
					log.Fatalf("Contract date ('%s') is incorrect! (Must be after then current date - '%s'.)\n",
						contractDateTime.String(), time.Now().String())
				}
			} else if templateFieldName == "Receiver" {
				receiver := r.FormValue(templateFieldName)
				match, _ := regexp.MatchString("[a-fA-F0-9]{66}", receiver)
				if match == false {
					log.Fatalf("Receiver '%s' is incorrect! (Must be the 66 digit HEX code.)", receiver)
				} else {
					log.Printf("Receiver '%s' is correct!", receiver)
				}
			} else if templateFieldName == "Receivers" {
				receivers := r.FormValue(templateFieldName)
				receiversRaw := strings.Split(receivers, ",")
				for index, receiverRaw := range receiversRaw {
					receiversRaw[index] = strings.Trim(receiverRaw, "[]'")
					match, _ := regexp.MatchString("[a-fA-F0-9]{66}", receiversRaw[index])
					if match == false {
						log.Fatalf("Receiver '%s' is incorrect! (Must be the 66 digit HEX code.)", receiversRaw[index])
					} else {
						log.Printf("Receiver '%s' is correct!", receiversRaw[index])
					}
				}
			}

			postValues[index] = r.FormValue(templateFieldName)
			templateData[templateFieldName] = postValues[index]
			fmt.Printf(templateFieldName+"='%s'\n", postValues[index])
		}
	}

	return templateData, contractInfo, id, nil
}

func BuildTransactionHTMLTemplateData(r *http.Request) (map[string]interface{}, *SmartContract, *ContractInfo, string) {
	// generating script template data
	scriptTemplateData, contractInfo, id, err := BuildTransactionScriptTemplateData(r)
	if err != nil {
		log.Fatal("Error when building deferred template data: ", err)
	}

	// generating smart-contract script (python)
	smartContract, err := ParseTemplate(contractInfo.TransactionScriptTemplatePath, scriptTemplateData)
	if err != nil {
		log.Fatal("Error when parsing script template file: ", err)
	}

	// writing python script to a file
	err = WriteFile(contractInfo.TransactionScriptTempFilename, smartContract)
	if err != nil {
		log.Fatal("Error when writing payment script to the file: ", err)
	}

	// run pylint
	staticAnalysisOfSmartContract := RunPylint(contractInfo.TransactionScriptTempFilename)

	// delete script file
	err = DeleteFile(contractInfo.TransactionScriptTempFilename)
	if err != nil {
		log.Fatal("Error when deleting payment script file: ", err)
	}

	htmlTemplateData := map[string]interface{}{
		"SmartContractCode":                 smartContract,
		"StaticAnalysisOfSmartContractCode": staticAnalysisOfSmartContract,
	}
	return htmlTemplateData, &SmartContract{Status: "Draft", Type: r.FormValue("ContractType"),
			CreationDate: GetCurrentTimestamp(), Price: CalculatePriceOfSmartContractFake(smartContract),
			LastStarted: "Unused", Comment: scriptTemplateData["TransactionMessage"].(string), Code: smartContract,
			Data: scriptTemplateData},
		contractInfo, id
}

type SmartContract struct {
	ID           string
	Status       string
	Type         string
	CreationDate string
	Price        string
	LastStarted  string
	Comment      string
	Code         string
	Data         map[string]interface{}
}

func (sm *SmartContract) ToJSON() string {
	b, err := json.Marshal(sm)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

type SmartContractTransaction struct {
	TT       string `json:"TT"`
	TST      string `json:"TST"`
	CODE     string `json:"CODE"`
	ANALYSIS string `json:"ANALYSIS"`
}

func (smt *SmartContractTransaction) ToJSON(escapeDoubleQuotes bool) (string, error) {
	if escapeDoubleQuotes {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(smt)
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	}
	bts, err := json.Marshal(smt)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

func GetCurrentTimestamp() string {
	// return time.Now().Format(time.RFC850)
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func CalculatePriceOfSmartContractFake(smartContractCode string) string {
	linesNumber := strings.Count(smartContractCode, "\n")
	return strconv.FormatFloat(float64(linesNumber)*SMART_CONTRACT_PRICE_COEFFICIENT_NZT, 'f', 2, 64)
}

type PylintError struct {
	FullText        string
	ErrorLineNumber int
}

func (pe *PylintError) ToJSON() string {
	b, err := json.Marshal(pe)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func FindErrorsInThePylintAnalysis(analysis string) ([]*PylintError, error) {
	// r := regexp.MustCompile(`(E:  (\d+), (\d+):[^\n]+)(\n)?`)
	r := regexp.MustCompile(`(E:[ ]+(\d+),[ ]+(\d+):[^\n]+)(\n)?`)
	matches := r.FindAllStringSubmatch(analysis, -1)

	pylintErrors := []*PylintError{}
	for i := range matches {
		errorLineNumber, err := strconv.Atoi(matches[i][2])
		if err != nil {
			return []*PylintError{}, err
		}
		pylintErrors = append(pylintErrors, &PylintError{FullText: matches[i][1], ErrorLineNumber: errorLineNumber})
	}
	return pylintErrors, nil
}

type ExplorerTransaction struct {
	Transaction             string `json:"transaction"`
	SmartContractStatus     string `json:"smcstatus"`
	SmartContractWasStarted bool   `json:"smcwasstarted"`
}

func (et *ExplorerTransaction) ToJSON(escapeDoubleQuotes bool) (string, error) {
	if escapeDoubleQuotes {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(et)
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	}
	bts, err := json.Marshal(et)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

func SendTransaction(mode string, url string, transaction interface{}) (string, error) {
	if strings.HasPrefix(url, "http") == false {
		url = "http://" + url
	}
	// make json
	
	fmt.Println("Url That we attack!!!::", url)
	
	var b []byte
	var err error
	if mode == "pvm" {
		var t string
		t, _ = transaction.(string)
		if err != nil {
			return "", err
		}
		b, err = json.Marshal(&ExplorerTransaction{Transaction: t})
		if err != nil {
			return "", err
		}
	} else {
		b, _ = transaction.([]byte)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response status:", resp.Status)
	fmt.Println("response headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

type PVMAction struct {
	Signature string `json:"SIGNATURE"`
	Sender    string `json:"SENDER"`
	Action    string `json:"ACTION"`
	SMS       string `json:"SMS"`
	TST       string `json:"TST"`
}

func SendAction(url string, pvmAction string) (string, error) {
	if strings.HasPrefix(url, "http") == false {
		url = "http://" + url
	}

    fmt.Println("Url That we attack 2!!!::", url)

	var err error

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(pvmAction)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response status:", resp.Status)
	fmt.Println("response headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func MapToBytes(m map[string]interface{}) ([]byte, error) {
	bytesFromMap, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytesFromMap, nil
}

func BytesToMap(bytes []byte) (map[string]interface{}, error) {
	mapRestored := make(map[string]interface{}, 0)
	err := json.Unmarshal(bytes, &mapRestored)
	if err != nil {
		return nil, err
	}
	return mapRestored, nil
}

// generate secure random session ID
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *mathRand.Rand = mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateSessionID(length int) string {
	return StringWithCharset(length, charset)
}

func ReadSessionIDFromBrowserCookies(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SMCE_SESSION_ID)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func RunPythonParser(smartContractCode string) string {
	cmd := exec.Command("python", "python_parser/parse_python_to_json.py", "--pp", smartContractCode)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		os.Stderr.WriteString("python-parser error message: " + err.Error() + "\n")
	}
	return string(cmdOutput.Bytes())
}

func EscapeCharacter(s, character string) string {
	return strings.Replace(s, character, `\`+character, -1)
}
