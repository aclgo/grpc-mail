package usecase

import (
	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/pkg/errors"
)

type mailUseCase struct {
	mailUC mail.MailUseCase
}

func NewmailUseCase(mailUC mail.MailUseCase) *mailUseCase {
	return &mailUseCase{
		mailUC: mailUC,
	}
}

func (u *mailUseCase) Send(data *models.MailBody) (string, error) {
	if err := u.mailUC.Send(data); err != nil {
		return "", errors.Wrap(err, "Send.mailUC.Send")
	}

	return mail.EmailSentSuccess, nil
}
