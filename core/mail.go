package core

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type EmailData struct {
	Title   string
	Message string
	Order   Order
}

func SendMail(
	to string,
	message string,
	order Order,
) {
	smtpHost := "mailpit"
	smtpPort := "1025"
	smtpAddress := smtpHost + ":" + smtpPort
	from := "eshop@softsec.com"

	data := EmailData{
		Title:   "Your order was successfully placed",
		Message: message,
		Order:   order,
	}

	tmpl, err := template.ParseFiles("core/mail.gohtml")
	if err != nil {
		logger.Errorf("Error parsing email template: %s", err)
		return
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		logger.Errorf("Error processing email template: %s", err)
		return
	}
	mime := "MIME-version: 1.0;\nContent-type: text/html; charset=iso-8859-1;\n\n"

	subject := "Subject: Order placed\n"
	msg := []byte(subject + mime + body.String())

	err = smtp.SendMail(smtpAddress, nil, from, []string{to}, msg)
	if err != nil {
		logger.Errorf("Error sending email to %s: %s", to, err)
		return
	}
	_createHtmlFile(data)
	logger.Info("Email was sent successfully")
}

func _createHtmlFile(data EmailData) {
	tmpl, err := template.ParseFiles("core/mail.gohtml")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	outputFile, err := os.Create("email_tmp.html")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Execute the template and write to the file
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	log.Println("HTML file created successfully!")
}
