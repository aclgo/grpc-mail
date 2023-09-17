package models

type MailBody struct {
	From     string
	To       string
	Subject  string
	Body     string
	Template string
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
