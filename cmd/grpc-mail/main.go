package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/internal/adapters/gmail"
	"github.com/aclgo/grpc-mail/internal/adapters/ses"
	grpcService "github.com/aclgo/grpc-mail/internal/mail/delivery/grpc/service"
	httpService "github.com/aclgo/grpc-mail/internal/mail/delivery/http/service"
	"github.com/aclgo/grpc-mail/internal/mail/usecase"
	"github.com/aclgo/grpc-mail/internal/server"
	"github.com/aclgo/grpc-mail/pkg/logger"
)

func main() {

	cfg := config.Load(".")
	cfg.ServiceHTTPPort = 3000

	logger := logger.NewapiLogger(cfg)

	ses := ses.NewSes(cfg)
	gmail := gmail.NewGmail(cfg)

	sesUC := usecase.NewmailUseCase(ses, logger)
	gmailUC := usecase.NewmailUseCase(gmail, logger)

	serviceSesHTTP := httpService.NewMailService(sesUC, logger)
	serviceSesGRPC := grpcService.NewMailService(sesUC, logger)

	serviceGmailHTTP := httpService.NewMailService(gmailUC, logger)
	serviceGmailGRPC := grpcService.NewMailService(gmailUC, logger)

	handlerSvcSes := server.NewHttpHandlerService("/ses", serviceSesHTTP)
	handlerSvcGmail := server.NewHttpHandlerService("/gmail", serviceGmailHTTP)

	server := server.NewServer(cfg,
		logger,
		[]*server.HttpHandlerService{
			handlerSvcSes,
			handlerSvcGmail,
		},
		[]*grpcService.MailService{
			serviceSesGRPC,
			serviceGmailGRPC,
		},
	)

	signal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	log.Fatal(server.Run(signal))
}
