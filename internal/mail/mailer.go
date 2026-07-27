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

	addr := fmt.Sprintf(
		"%s:%d",
		m.host,
		m.port,
	)

	message := fmt.Appendf(
		nil,
		"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s",
		subject,
		body,
	)

	var auth smtp.Auth

	if m.username != "" {
		auth = smtp.PlainAuth(
			"",
			m.username,
			m.password,
			m.host,
		)
	}

	return smtp.SendMail(
		addr,
		auth,
		m.from,
		[]string{to},
		message,
	)
}
