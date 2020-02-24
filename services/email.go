package services

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendMail handles sending based on parameters via sendgrid
func SendMail(toMail string, name string, subject string, body string) (*rest.Response, error) {
	from := mail.NewEmail("Emania", "ewachira254@gmail.com")
	to := mail.NewEmail(name, toMail)
	message := mail.NewSingleEmail(from, subject, to, body, body)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_KEY"))
	response, err := client.Send(message)

	return response, err
}
