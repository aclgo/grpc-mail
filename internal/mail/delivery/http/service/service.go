package service

import (
	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/pkg/logger"
)

type MailService struct {
	mailUC mail.MailUseCase
	logger logger.Logger
}

func NewMailService(mailUC mail.MailUseCase, logger logger.Logger) *MailService {
	return &MailService{
		mailUC: mailUC,
		logger: logger,
	}
}
