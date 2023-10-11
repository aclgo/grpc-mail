package service

import (
	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"github.com/aclgo/grpc-mail/proto"
)

type MailServiceLoad struct {
	NameService string
	Service     mail.MailUseCase
}

func NewMailServiceLoad(name string, service mail.MailUseCase) *MailServiceLoad {
	return &MailServiceLoad{
		NameService: name,
		Service:     service,
	}
}

type MailService struct {
	mailUCS []*MailServiceLoad
	logger  logger.Logger
	proto.UnimplementedMailServiceServer
}

func NewMailService(logger logger.Logger, mailsl ...*MailServiceLoad) *MailService {
	return &MailService{
		logger:  logger,
		mailUCS: mailsl,
	}
}
