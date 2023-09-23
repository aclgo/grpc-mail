package gmail

import (
	"fmt"
	"net/smtp"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/pkg/errors"
)

var (
	MessageFormat = "Subject:%s\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s"
)

type Gmail struct {
	Auth smtp.Auth
	Addr string
}

func NewGmail(config *config.Config) *Gmail {
	auth := smtp.PlainAuth(
		config.Gmail.Identity,
		config.Gmail.Username,
		config.Gmail.Password,
		config.Gmail.Host,
	)

	gmail := Gmail{
		Auth: auth,
		Addr: fmt.Sprintf("%s:%d", config.Gmail.Host, config.Gmail.Port),
	}

	return &gmail
}

func (g *Gmail) Send(data *models.MailBody) error {

	bodyMessage := fmt.Sprintf(MessageFormat, data.Subject, data.Body)

	if err := smtp.SendMail(
		g.Addr,
		g.Auth,
		data.From,
		[]string{data.To},
		[]byte(bodyMessage),
	); err != nil {
		return errors.Wrap(err, "Send.smtp.SendMail")
	}

	return nil
}
