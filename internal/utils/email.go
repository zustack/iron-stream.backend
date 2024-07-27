package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

func SendEmail(token int, email string, subject string) error {
	secretPassword := os.Getenv("EMAIL_SECRET_KEY")
	auth := smtp.PlainAuth(
		"",
		"agustfricke@gmail.com",
		secretPassword,
		"smtp.gmail.com",
	)

	path := os.Getenv("ROOT_PATH") + "/web/templates/email.html"
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return fmt.Errorf("error al analizar la plantilla: %w", err)
	}

	data := struct {
		Token int
	}{
		Token: token,
	}

	var bodyContent bytes.Buffer
	if err := tmpl.Execute(&bodyContent, data); err != nil {
		return fmt.Errorf("error al ejecutar la plantilla: %w", err)
	}

	emailContent := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s", email, subject, bodyContent.String())

	if err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"agustfricke@gmail.com",
		[]string{email},
		[]byte(emailContent),
	); err != nil {
		return fmt.Errorf("error al enviar el correo electrónico: %w", err)
	}

	return nil
}
