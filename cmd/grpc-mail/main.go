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

	cfg.OtelExporter = "otlp"

	cfg.Meter.Name = "meter-name"
	cfg.Tracer.Name = "tracer name"

	cfg.MeterExporterURL = "otel-collector:4317"
	cfg.Tracer.TracerExporterURL = "http://zipkin:9411/api/v2/spans"

	logger := logger.NewapiLogger(cfg)

	logger.Info("logger init")

	tel, err := telemetry.NewProvider(cfg, logger)
	defer func() {
		if err != nil {
			logger.Errorf("cannot initialize telemetry: %v", err)
			os.Exit(1)
		}
	}()

	defer tel.Shutdown()

	logger.Info("provider init")

	ses := ses.NewSes(cfg)
	gmail := gmail.NewGmail(cfg)

	sesUC := usecase.NewmailUseCase(ses, logger)
	gmailUC := usecase.NewmailUseCase(gmail, logger)

	servicesHttpLoad := []*httpService.MailServiceLoad{
		httpService.NewMailServiceLoad("ses", sesUC),
		httpService.NewMailServiceLoad("gmail", gmailUC),
	}

	// HTTP services
	servicesHTTP := httpService.NewMailService(logger, tel, servicesHttpLoad...)

	// handlers http
	handlerHTTP := server.NewHttpHandlerService("/api/v1/send", servicesHTTP)

	// GRPC services
	servicesGRPC := grpcService.NewMailServices(
		logger,
		tel,
		grpcService.NewMailServiceLoad("ses", sesUC),
		grpcService.NewMailServiceLoad("gmail", gmailUC),
	)

	server := server.NewServer(cfg,
		logger,
		handlerHTTP,
		servicesGRPC,
		tel,
	)

	signal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(signal); err != nil {
		log.Fatal(err)
	}
}
