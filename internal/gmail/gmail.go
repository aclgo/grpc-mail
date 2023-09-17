package gmail

import (
	"fmt"
	"net/smtp"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/pkg/errors"
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

	gmail := &Gmail{
		Auth: auth,
		Addr: fmt.Sprintf("%s:%s", config.Gmail.Host, config.Gmail.Port),
	}

	return gmail
}

func (g *Gmail) Send(data *models.MailBody) error {
	if err := smtp.SendMail(
		g.Addr,
		g.Auth,
		data.From,
		[]string{data.To},
		[]byte(data.Body),
	); err != nil {
		return errors.Wrap(err, "Send.smtp.SendMail")
	}

	return nil
}
