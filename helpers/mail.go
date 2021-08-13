package helpers

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/edwinnduti/dora/models"
)

var (
	USERNAME = "nduti316@gmail.com"
	PASSWORD = "tkqhbdrjjxhkrkyr"
	SCHOOL_NAME = "DORA SCHOOL"
)

func SendMailTo(student *models.Student) error {
	var to_send_email = []string{student.Email}

	message := &models.Message{
		TO_EMAIL : to_send_email,
		To_Name: student.FullName,
		From_Email : USERNAME,
		From_Name : SCHOOL_NAME,
		Subject: "Hello!",
	}
	body := fmt.Sprintf("Dear %s,\nWe saw that you signed-in andwe would like to urge you to fill in the questions as soon as possible.This questions are important for the evaluation of all lecturer performance necessary for enhancing better study for all students in the campus.\nKindly fill the for by today!\nRegards,\n%s",message.To_Name,message.From_Name)


	mime := fmt.Sprintln("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n")

	msg := []byte("From: "+message.From_Email+"\r\n"+"Subject:"+message.Subject+"\r\n"+mime+"<html><head><style>#rcorners {border-radius: 25px; background: #8AC007; padding: 20px; width: 90%; height: 100%;}</style></head><body id=\"rcorners\"><br><pre>"+body+"</pre></body></html>")

	message.Body = msg

	smtpServer := &models.SmtpServer{
		Host : "smtp.gmail.com",
		Port : "587",
	}

	address := fmt.Sprintln("%s:%s", smtpServer.Host, smtpServer.Port)
	auth := smtp.PlainAuth("", USERNAME, PASSWORD, smtpServer.Host)
	err := smtp.SendMail(address, auth, message.From_Email, message.TO_EMAIL, message.Body)
	if err != nil{
		return err
	}

	log.Println("Email Sent Successfully!")
	return nil
}
