package service

import (
	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/telemetry"
	"github.com/aclgo/grpc-mail/pkg/logger"
)

type MailService struct {
	svcsMail map[string]mail.MailUseCase
	logger   logger.Logger
	tel      *telemetry.Provider
}

type MailServiceLoad struct {
	serviceName string
	mailService mail.MailUseCase
}

func NewMailServiceLoad(svcName string, mailService mail.MailUseCase) *MailServiceLoad {
	return &MailServiceLoad{
		serviceName: svcName,
		mailService: mailService,
	}
}

func NewMailService(logger logger.Logger, telemetry *telemetry.Provider, svcs ...*MailServiceLoad) *MailService {

	mailServices := MailService{
		svcsMail: make(map[string]mail.MailUseCase),
		logger:   logger,
		tel:      telemetry,
	}

	for _, value := range svcs {
		_, ok := mailServices.svcsMail[value.serviceName]
		if ok {
			mailServices.logger.Warnf("service name %s exist", value.serviceName)
			continue
		}

		mailServices.svcsMail[value.serviceName] = value.mailService
	}

	return &mailServices
}
