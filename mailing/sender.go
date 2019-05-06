package mailing

import (
	"fmt"
	"log"
	"net/smtp"

	"noosphere.foundation/smart-contract-editor/config"
	"noosphere.foundation/smart-contract-editor/utils"
)

// var emailValidationBlockResult *config.EMailValidationBlockResult
var timer = &utils.Timer{}

// func init() {
// 	configFileContent := config.Read(config.FILE_PATH)
// 	var err error
// 	emailValidationBlockResult, err = config.GetEMailBlockConfigData(configFileContent)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func PrepareEMailConfirmationTemplate(validationEMailFileName, firstName, lastName, userEmail, validationLink string) {
	var greeting string = "Hi " + firstName + " " + lastName + ","
	templateData := struct {
		LineOfDashesAbove string
		FullName          string
		LineOfDashesUnder string
		ValidationLink    string
	}{
		utils.ProduceStringWithDashes(len(greeting)),
		firstName + " " + lastName,
		utils.ProduceStringWithDashes(len(greeting)),
		validationLink,
	}

	emailText, err := utils.ParseTemplate(validationEMailFileName, templateData)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		timer.Start("sending validation email...")
		_, err := SendEMail(config.Data.EMailBlock.EMailValidation, config.Data.EMailBlock.EMailValidationPassword,
			userEmail, "text/html", "Please confirm your email", emailText)
		if err != nil {
			log.Fatal(err)
		}
		timer.End("time elapsed:")
	}()
}

func SendEMail(emailFrom, emailFromPassword, userEmail, contentType, subject, emailBody string) (bool, error) {
	auth := smtp.PlainAuth("", emailFrom, emailFromPassword, "64.233.162.108")

	mime := fmt.Sprintf("MIME-version: 1.0;\nContent-Type: %s; charset=\"UTF-8\";\n\n", contentType)
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n%s", userEmail, subject, mime, emailBody)
	addr := "64.233.162.108:587"

	if err := smtp.SendMail(addr, auth, emailFrom, []string{userEmail}, []byte(msg)); err != nil {
		return false, err
	}
	return true, nil
}
