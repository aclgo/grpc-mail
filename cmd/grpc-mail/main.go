package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/e2e"
	"github.com/aclgo/grpc-mail/internal/adapters/gmail"
	"github.com/aclgo/grpc-mail/internal/adapters/ses"
	"github.com/aclgo/grpc-mail/internal/mail"
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

	// cfg.Meter.Name = "grpc-mail"
	// cfg.Tracer.Name = "grpc-mail"

	cfg.MeterExporterURL = "http://localhost:4317"
	cfg.Tracer.TracerExporterURL = "http://localhost:9411/api/v2/spans"

	logger := logger.NewapiLogger(cfg)

	logger.Info("logger init")

	tel := telemetry.NewProvider(cfg, logger)

	defer tel.Shutdown()

	tracer := tel.TracerProvider.Tracer("grpc-mail")
	meter := tel.MeterProvider.Meter("grpc-mail")

	observer := mail.NewObserver(logger, tracer, meter)

	logger.Info("observer init")

	ses := ses.NewSes(cfg)
	gmail := gmail.NewGmail(cfg)

	sesUC := usecase.NewmailUseCase(ses, logger)
	gmailUC := usecase.NewmailUseCase(gmail, logger)

	servicesHttpLoad := []*httpService.MailServiceLoad{
		httpService.NewMailServiceLoad("ses", sesUC),
		httpService.NewMailServiceLoad("gmail", gmailUC),
	}

	// HTTP services
	servicesHTTP := httpService.NewMailService(logger, observer, servicesHttpLoad...)

	// HTTP handlers
	handlerHTTP := server.NewHttpHandlerService("/api/v1/send", servicesHTTP)

	// GRPC services
	servicesGRPC := grpcService.NewMailServices(
		logger,
		observer,
		grpcService.NewMailServiceLoad("ses", sesUC),
		grpcService.NewMailServiceLoad("gmail", gmailUC),
	)

	server := server.NewServer(cfg,
		logger,
		handlerHTTP,
		servicesGRPC,
	)

	signal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		time.Sleep(time.Second * 10)
		e2e.RunGRPC(fmt.Sprintf("localhost:%v", cfg.ServiceGRPCPort), logger)
		e2e.RunHTTP(fmt.Sprintf("http://localhost:%v/api/v1/send", cfg.ServiceHTTPPort), logger)
	}()

	if err := server.Run(signal); err != nil {
		log.Fatal(err)
	}
}
