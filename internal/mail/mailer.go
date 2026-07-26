package mail

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewMailer(
	host string,
	port int,
	username string,
	password string,
	from string,
) *Mailer {

	return &Mailer{

		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}

}

func (m *Mailer) Send(
	to string,
	subject string,
	body string,
) error {

	auth := smtp.PlainAuth(
		"",
		m.username,
		m.password,
		m.host,
	)

	message := []byte(
		fmt.Sprintf(
			"Subject: %s\r\n\r\n%s",
			subject,
			body,
		),
	)

	return smtp.SendMail(
		fmt.Sprintf(
			"%s:%d",
			m.host,
			m.port,
		),
		auth,
		m.from,
		[]string{to},
		message,
	)
}
