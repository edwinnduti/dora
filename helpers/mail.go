package helpers

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/edwinnduti/dora/models"
)

var (
	USERNAME    = "tukscit@gmail.com"
	PASSWORD    = "tukit2017"
	SCHOOL_NAME = "TECHNICAL UNIVERSITY OF KENYA"
)

func SendMailTo(student *models.Student) error {
	var to_send_email = []string{student.Email}

	message := &models.Message{
		TO_EMAIL:   to_send_email,
		To_Name:    student.FullName,
		From_Email: USERNAME,
		From_Name:  SCHOOL_NAME,
		Subject:    "Hello!",
	}
	body := fmt.Sprintf("Dear %s,<br/><br/>We saw that you signed-up and we would like to urge you to fill in the questions as soon as possible.These questions are important for the evaluation of all lecturers performance, necessary for enhancing better study for all students in the campus.<br/>Kindly fill the form within 24hours!<br/><br/>Regards,<br/>%s.", message.To_Name, message.From_Name)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte("From: " + message.From_Email + "\r\n" + "Subject:" + message.Subject + "\r\n" + mime + "<html><head><style>#rcorners {border-radius: 25px; background: #8AC007; padding: 20px; width: 90%; height: 100%;}</style></head><body id=\"rcorners\"><br><h2 align=\"center\">" + message.From_Name + "</h2><br/><h3>" + body + "</h3></body></html>")

	message.Body = msg

	smtpServer := &models.SmtpServer{
		Host: "smtp.gmail.com",
		Port: "587",
	}

	address := fmt.Sprintf("%s:%s", smtpServer.Host, smtpServer.Port)
	auth := smtp.PlainAuth("", USERNAME, PASSWORD, smtpServer.Host)
	err := smtp.SendMail(address, auth, message.From_Email, message.TO_EMAIL, message.Body)
	if err != nil {
		return err
	}

	log.Println("Email Sent Successfully!")
	return nil
}
