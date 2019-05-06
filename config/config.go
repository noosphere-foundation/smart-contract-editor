package config

import (
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

const (
	CONFIG_FILE_PATH = "settings.cfg"

	EMAIL_VALIDATION_PARAMETER_NAME          = "emailValidation"
	EMAIL_VALIDATION_PASSWORD_PARAMETER_NAME = "emailValidationPassword"
	SERVER_ADDRESS_PARAMETER_NAME            = "server"
	URL_BASE_PARAMETER_NAME                  = "urlbase"
	PVM_URL                                  = "pvmurl"
	PVM_ACTION                               = "pvmaction"
)

type Config struct {
	EMailBlock  *EMailValidationBlockResult
	ServerBlock *ServerAddressBlockResult
}

type EMailValidationBlock struct {
	EMailValidationParameterName         string
	EMailValidationPasswordParameterName string
}

type ServerAddressBlock struct {
	ServerAddressParameterName string
	UrlBaseParameterName       string
	PVMURLParameterName        string
	PVMActionParameterName     string
}

type EMailValidationBlockResult struct {
	EMailValidation         string
	EMailValidationPassword string
}

type ServerAddressBlockResult struct {
	ServerAddress string
	Port          string
	URLBase       string
	PVMURL        string
	PVMAction     string
}

var Data *Config = &Config{}
var emailBlock *EMailValidationBlock = &EMailValidationBlock{}
var serverBlock *ServerAddressBlock = &ServerAddressBlock{}

func init() {
	log.Println("reading settings.cfg...")
	// setting config parameters names
	emailBlock.EMailValidationParameterName = EMAIL_VALIDATION_PARAMETER_NAME
	emailBlock.EMailValidationPasswordParameterName = EMAIL_VALIDATION_PASSWORD_PARAMETER_NAME

	serverBlock.ServerAddressParameterName = SERVER_ADDRESS_PARAMETER_NAME
	serverBlock.UrlBaseParameterName = URL_BASE_PARAMETER_NAME
	serverBlock.PVMURLParameterName = PVM_URL
	serverBlock.PVMActionParameterName = PVM_ACTION

	// reading the config file
	configFileContent := Read(CONFIG_FILE_PATH)

	// getting parsed config data
	emailBlock, err := GetEMailBlockConfigData(configFileContent)
	if err != nil {
		log.Fatal("Error when parsing email block in the config:", err)
	}

	serverBlock, err := GetServerBlockConfigData(configFileContent)
	if err != nil {
		log.Fatal("Error when parsing server block in the config:", err)
	}

	Data.EMailBlock = emailBlock
	Data.ServerBlock = serverBlock
	log.Println("ok")
}

func Read(filepath string) string {
	configFileContent := string(ReadFile(filepath))

	// remove lines starting with the '#' symbol (comments)
	configWithoutComments := []string{}
	for _, s := range strings.Split(configFileContent, "\n") {
		var cutted string = strings.TrimSpace(s)
		if !strings.HasPrefix(cutted, "#") && len(cutted) > 0 {
			configWithoutComments = append(configWithoutComments, cutted)
		}
	}

	return strings.Join(configWithoutComments, "\n")
}

func ReadFile(filepath string) []byte {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func parseConfigParameterValue(configParameterName, errorMessage, configContent string) (string, error) {
	r := regexp.MustCompile(configParameterName + "=([^\n]+)\n?")
	submatches := r.FindStringSubmatch(configContent)

	if len(submatches) > 1 {
		return submatches[1], nil
	}
	return "", errors.New(errorMessage)
}

func GetEMailBlockConfigData(configFileContent string) (*EMailValidationBlockResult, error) {
	email, err := parseConfigParameterValue(emailBlock.EMailValidationParameterName,
		"Config parsing error: getting email validation.", configFileContent)
	if err != nil {
		return &EMailValidationBlockResult{}, err
	}
	password, err := parseConfigParameterValue(emailBlock.EMailValidationPasswordParameterName,
		"Config parsing error: getting email validation password.", configFileContent)
	if err != nil {
		return &EMailValidationBlockResult{}, err
	}
	return &EMailValidationBlockResult{EMailValidation: email, EMailValidationPassword: password}, nil
}

func GetServerBlockConfigData(configFileContent string) (*ServerAddressBlockResult, error) {
	serverAddress, err := parseConfigParameterValue(serverBlock.ServerAddressParameterName,
		"Config parsing error: getting server address.", configFileContent)
	if err != nil {
		return &ServerAddressBlockResult{}, err
	}

	colonIndex := strings.Index(serverAddress, ":")
	port := serverAddress[colonIndex+1 : len(serverAddress)]

	urlBase, err := parseConfigParameterValue(serverBlock.UrlBaseParameterName, "Config parsing error: getting url base",
		configFileContent)
	if err != nil {
		return &ServerAddressBlockResult{}, err
	}

	pvmurl, err := parseConfigParameterValue(serverBlock.PVMURLParameterName, "Config parsing error: getting pvm url",
		configFileContent)
	if err != nil {
		return &ServerAddressBlockResult{}, err
	}

	pvmaction, err := parseConfigParameterValue(serverBlock.PVMActionParameterName,
		"Config parsing error: getting pvm action", configFileContent)
	if err != nil {
		return &ServerAddressBlockResult{}, err
	}

	return &ServerAddressBlockResult{ServerAddress: serverAddress, Port: port, URLBase: urlBase, PVMURL: pvmurl,
		PVMAction: pvmaction}, nil
}
