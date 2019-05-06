package config

import (
	"log"
	"strings"
	"testing"
)

const TEST_CONFIG_FILE_PATH = "settings_test.cfg"

func TestParseConfigParameterValue(t *testing.T) {
	content := Read(TEST_CONFIG_FILE_PATH)

	emailValidationValue, err := parseConfigParameterValue("emailValidation",
		"Error when parsing config parameter: email validation!", string(content))
	if err != nil {
		log.Fatal(err)
	}
	if emailValidationValue != "john@snow.com" {
		t.Errorf("Expected validation email is 'john@snow.com', but found '%s'\n", emailValidationValue)
	}

	emailValidationPasswordValue, err := parseConfigParameterValue("emailValidationPassword",
		"Error when parsing config parameter: email validation password!", string(content))
	if err != nil {
		log.Fatal(err)
	}
	if emailValidationPasswordValue != "qwerty777" {
		t.Errorf("Expected validation email password is 'qwerty777', but found '%s'\n", emailValidationPasswordValue)
	}
}

func TestGetEMailBlockConfigData(t *testing.T) {
	testConfigContent := Read(TEST_CONFIG_FILE_PATH)

	emailBlockConfigData, err := GetEMailBlockConfigData(string(testConfigContent))
	if err != nil {
		log.Fatal(err)
	}

	if emailBlockConfigData.EMailValidation != "john@snow.com" {
		t.Errorf("Expected validation email is 'john@snow.com', but found '%s'\n",
			emailBlockConfigData.EMailValidation)
	}

	if emailBlockConfigData.EMailValidationPassword != "qwerty777" {
		t.Errorf("Expected validation email password is 'qwerty777', but found '%s'\n",
			emailBlockConfigData.EMailValidationPassword)
	}
}

func TestRead(t *testing.T) {
	config := Read(TEST_CONFIG_FILE_PATH)
	splittedConfig := strings.Split(config, "\n")
	expected := []string{"emailValidation=john@snow.com", "emailValidationPassword=qwerty777", "server=localhost:1234"}

	if len(splittedConfig) != len(expected) {
		t.Errorf("expected len of config is '%d', but found '%d'\n", len(expected), len(splittedConfig))
	}

	for index := 0; index < len(splittedConfig); index++ {
		if splittedConfig[index] != expected[index] {
			t.Errorf("expected line is '%s', but found '%s'\n", expected[index], splittedConfig[index])
		}
	}
}
