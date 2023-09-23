package models

import (
	"errors"
	"fmt"
)

type MailBody struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Template string `json:"template"`
}

func NewMailBody(from string, to string, subject, body, template string) *MailBody {
	return &MailBody{
		From:     from,
		To:       to,
		Subject:  subject,
		Body:     body,
		Template: template,
	}
}

func (m *MailBody) Validate() error {

	f := func(arg string) string {
		return fmt.Sprintf("empty %s", arg)
	}

	if m.From == "" {
		return errors.New(f("from"))
	}

	if m.To == "" {
		return errors.New(f("to"))
	}
	if m.Subject == "" {
		return errors.New(f("subject"))
	}
	if m.Body == "" {
		return errors.New(f("body"))
	}

	if m.Template == "" {
		return nil
	}

	return nil
}
