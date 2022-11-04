package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/asaskevich/govalidator"
)

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalln("missing environment variable: ", key)
	}
	return value
}

func main() {
	fmt.Println("Hello, world!")
	// get working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// get excel filename from command line arguments
	filename := os.Args[1]
	if filename == "" {
		log.Fatalln("invalid excel file name: e.g. go run main.go test.xlsx")
	}

	// mail service info -- Host , Port , Username , Password , FromName , FromMail
	host := getEnvOrFatal("MAIL_HOST")
	port := getEnvOrFatal("MAIL_PORT")
	username := getEnvOrFatal("MAIL_USERNAME")
	password := getEnvOrFatal("MAIL_PASSWORD")
	fromName := getEnvOrFatal("MAIL_FROM_NAME")
	fromMail := getEnvOrFatal("MAIL_FROM_MAIL")
	mailSubject := getEnvOrFatal("MAIL_SUBJECT")

	// create application service
	config := NewMailConfig(host, port, username, password, fromMail, fromName)

	// read from excel
	path := dir + "/" + filename
	xlsx, err := excelize.OpenFile(path)
	if err != nil {
		log.Fatalln("error while opening excel file: ", err)
	}

	// emails
	var emails []string

	// read rows
	rows, err := xlsx.Rows("Sheet1")
	if err != nil {
		log.Fatal("error while reading rows: ", err)
	}

	// iterate over rows
	i := 0
	for rows.Next() {
		row := rows.Columns()

		// every row should just 1 column (email) and that is the first row
		if len(row) != 1 {
			log.Fatalln("invalid row length: ", len(row), " index: ", i, " row: ", row)
		}

		// check email format
		if !govalidator.IsEmail(row[0]) {
			log.Fatalln("invalid email format on row: ", i+1, " email: ", row[0])
		}

		// add email to emails
		emails = append(emails, row[0])

		i++
	}

	// log all emails are valid
	log.Println("all emails are valid and found: ", len(emails))

	// start sending emails. If success add to success queue, if failed add to failed queue and remove from emails list
	var successQueue []string
	var failedQueue []string

	emailsLength := len(emails)

	for index, email := range emails {
		// create mail content
		mailContent, err := PrepareEmailContent(dir, email, mailSubject)
		if err != nil {
			// remove email from emails array
			if index < emailsLength-1 {
				// add email to failed queue
				failedQueue = append(failedQueue, email)
			}

			// log error
			log.Println((index + 1), " - ERROR: ", email, " ", err)
			continue
		}

		// send email
		err = SendEmailViaSmtp(config, mailContent)
		if err != nil {
			// remove email from emails list
			if index < emailsLength-1 {
				// add email to failed queue
				failedQueue = append(failedQueue, email)
			}
			// log error
			log.Println((index + 1), " - ERROR: ", email, " ", err)
			continue
		}

		// add email to success queue
		successQueue = append(successQueue, email)

		// log success
		log.Println((index + 1), " - SUCCESS: ", email)
	}

	// write success queue to json file
	successQueueJson, err := json.Marshal(successQueue)
	if err != nil {
		log.Fatalln("error while converting success queue to json: ", err)
	}

	err = ioutil.WriteFile(dir+"/output/success.json", successQueueJson, 0644)
	if err != nil {
		log.Fatalln("error while writing success queue to json file: ", err)
	}

	// write failed queue to json file
	failedQueueJson, err := json.Marshal(failedQueue)
	if err != nil {
		log.Fatalln("error while converting failed queue to json: ", err)
	}

	err = ioutil.WriteFile(dir+"/output/failed.json", failedQueueJson, 0644)
	if err != nil {
		log.Fatalln("error while writing failed queue to json file: ", err)
	}

	// write emails list to json file
	emailsJson, err := json.Marshal(emails)
	if err != nil {
		log.Fatalln("error while converting emails list to json: ", err)
	}

	err = ioutil.WriteFile(dir+"/output/emails.json", emailsJson, 0644)
	if err != nil {
		log.Fatalln("error while writing emails list to json file: ", err)
	}

	// log success
	log.Println("SUCCESS: ", len(successQueue))
	// log failed
	log.Println("FAILED: ", len(failedQueue))
	// log remaining
	log.Println("REMAINING: ", len(emails))
}
