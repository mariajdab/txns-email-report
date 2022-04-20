package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mariajdab/txns-email-report/database"
	"github.com/mariajdab/txns-email-report/models"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 2 {
		log.Fatalf("Wrong number of arguments: Usage %s <path-to-file>", os.Args[0])
	}

	filePath := os.Args[1]

	db, err := database.OpenDB()

	if err != nil {
		log.Fatal(err)
	}

	if err := database.CreateTable(db); err != nil {
		log.Fatal(err)
	}

	data, reportTxns, err := ProcessFile(filePath)
	if err != nil {
		log.Fatalln("an error happened processing the file")
	}

	fmt.Println(reportTxns)

	htmlBody, err := SendEmail(reportTxns)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(htmlBody)

	database.InsertTxns(db, data)

}

func ProcessFile(file string) ([]models.AccountTxn, models.Report, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, models.Report{}, err
	}

	defer f.Close()

	csvFile := csv.NewReader(f)

	fl, err := ReadFirstLine(csvFile)
	if err != nil {
		return nil, models.Report{}, err
	}

	if !ValidateFirstLine(fl) {
		err = errors.New("mismatch fields expected:(Id, Date, Transaction)")
		return nil, models.Report{}, err
	}

	var data []models.AccountTxn
	var r models.Report

	for {
		row, err := csvFile.Read()
		if err != nil {
			if err == io.EOF {
				r = NewReport()
				break
			}
			return nil, models.Report{}, errors.New("an error happened reading the csv")
		}

		rowData, err := ParseData(row)
		if err != nil {
			return nil, models.Report{}, err
		}

		ProcessReport(rowData.Transaction, row[1])

		data = append(data, models.AccountTxn{
			Id:          rowData.Id,
			Date:        rowData.Date,
			Transaction: rowData.Transaction,
		})
	}
	return data, r, nil
}

func SendEmail(report models.Report) (string, error) {

	user := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")

	addr := host + port
	from := user + "@emailtrap.io"

	to := []string{
		"mariajdab@gmail.com",
	}

	t, err := template.ParseFiles("email.html")
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, report); err != nil {
		return "", err
	}
	body := buf.String()

	//mime := "MIME-version: 1.0;"

	msg := []byte("From: 050e7e2127bc03@emailtrap.io\r\n" +
		"To: mariajdab@gmail.com\r\n" +
		"Subject: Test mail\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", user, password, host)

	err = smtp.SendMail(addr, auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Fake Email sent successfully")

	return body, nil
}
