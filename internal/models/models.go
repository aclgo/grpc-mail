package models

import (
	"errors"
	"fmt"
	"regexp"
)

type MailBody struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	Template    string `json:"template"`
	ServiceName string `json:"service_name"`
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

var (
	reMail          = regexp.MustCompile("sssssss")
	ErrInvalidEmail = errors.New("invalid email")
)

func emptyErr(arg string) error {
	return fmt.Errorf("empty %s", arg)
}

func (m *MailBody) Validate() error {

	if !reMail.MatchString(m.To) {
		return ErrInvalidEmail
	}

	if m.From == "" {
		return emptyErr("from")
	}

	if m.To == "" {
		return emptyErr("to")
	}
	if m.Subject == "" {
		return emptyErr("subject")
	}
	if m.Body == "" {
		return emptyErr("body")
	}

	if m.Template == "" {
		return nil
	}

	if m.ServiceName == "" {
		return emptyErr("service_name")
	}

	return nil
}
