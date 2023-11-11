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

func NewMailBody(from string, to string, subject, body, template, svcName string) *MailBody {
	return &MailBody{
		From:        from,
		To:          to,
		Subject:     subject,
		Body:        body,
		Template:    template,
		ServiceName: svcName,
	}
}

var (
	reMailValid     = regexp.MustCompile(`^[a-zA-Z0-9._%-+]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	ErrInvalidEmail = errors.New("invalid email")
)

func emptyErr(arg string) error {
	return fmt.Errorf("empty %s", arg)
}

func (m *MailBody) Validate() error {

	if m.From == "" {
		return emptyErr("from")
	}

	if m.To == "" {
		return emptyErr("to")
	}

	if !reMailValid.MatchString(m.To) {
		return ErrInvalidEmail
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
