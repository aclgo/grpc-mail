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
	mailUCS  map[string]mail.MailUseCase
	logger   logger.Logger
	observer *mail.Observer
	proto.UnimplementedMailServiceServer
}

func NewMailServices(logger logger.Logger, obs *mail.Observer, mailsl ...*MailServiceLoad) *MailService {

	svcs := MailService{
		mailUCS:  make(map[string]mail.MailUseCase),
		logger:   logger,
		observer: obs,
	}

	for _, v := range mailsl {
		_, ok := svcs.mailUCS[v.NameService]
		if !ok {
			svcs.mailUCS[v.NameService] = v.Service
			continue
		}

		logger.Warnf("service name %s exist", v.NameService)
	}

	return &svcs
}
