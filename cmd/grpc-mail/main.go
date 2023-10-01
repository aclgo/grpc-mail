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
	"github.com/aclgo/grpc-mail/internal/telemetry"
	"github.com/aclgo/grpc-mail/pkg/logger"
)

func main() {

	cfg := config.Load(".")

	logger := logger.NewapiLogger(cfg)

	logger.Info("logger init")

	ses := ses.NewSes(cfg)
	gmail := gmail.NewGmail(cfg)

	sesUC := usecase.NewmailUseCase(ses, logger)
	gmailUC := usecase.NewmailUseCase(gmail, logger)

	// HTTP services
	serviceSesHTTP := httpService.NewMailService(sesUC, logger)
	serviceGmailHTTP := httpService.NewMailService(gmailUC, logger)

	// handlers http
	handlerSvcSes := server.NewHttpHandlerService("/ses", serviceSesHTTP)
	handlerSvcGmail := server.NewHttpHandlerService("/gmail", serviceGmailHTTP)

	// GRPC services
	servicesGRPC := grpcService.NewMailService(logger, sesUC)

	// providerConfig := telemetry.ProviderConfig{
	// 	Logger: logger,
	// }

	// providerConfig.Start()

	provider, err := telemetry.NewProvider(cfg, logger)
	defer func() {
		if err != nil {
			logger.Errorf("cannot initialize telemetry: %v", err)
			os.Exit(1)
		}
	}()

	defer provider.Shutdown()

	server := server.NewServer(cfg,
		logger,
		[]*server.HttpHandlerService{
			handlerSvcSes,
			handlerSvcGmail,
		},
		servicesGRPC,
		provider,
	)

	signal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(signal); err != nil {
		log.Fatal(err)
	}
}
