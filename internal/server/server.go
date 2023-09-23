package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/aclgo/grpc-mail/config"
	grpcService "github.com/aclgo/grpc-mail/internal/mail/delivery/grpc/service"
	httpService "github.com/aclgo/grpc-mail/internal/mail/delivery/http/service"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"google.golang.org/grpc"
)

type Server struct {
	config       *config.Config
	logger       logger.Logger
	servicesHTTP []*HttpHandlerService
	serviceGRPC  []*grpcService.MailService
	stopFn       sync.Once
}

type HttpHandlerService struct {
	pattern string
	service *httpService.MailService
}

func NewHttpHandlerService(pattern string, service *httpService.MailService) *HttpHandlerService {
	return &HttpHandlerService{
		pattern: pattern,
		service: service,
	}
}

func NewServer(cfg *config.Config,
	logger logger.Logger,
	svcsHTTP []*HttpHandlerService,
	svcGRPC []*grpcService.MailService) *Server {
	return &Server{
		config:       cfg,
		logger:       logger,
		servicesHTTP: svcsHTTP,
		serviceGRPC:  svcGRPC,
	}
}

func (s *Server) Run(ctxSignal context.Context) error {

	ctxHttp := context.Background()

	var (
		errHTTP = make(chan error)
		errGRPC = make(chan error)
	)

	go func() {
		err := s.httpRun(ctxHttp)
		if err != nil {
			s.logger.Errorf("Run:%v", err)
			errHTTP <- fmt.Errorf("Run:%v", err)
		}
	}()

	go func() {
		err := s.grpcRun()
		if err != nil {
			s.logger.Errorf("Run:%v", err)
			errGRPC <- fmt.Errorf("Run:%v", err)
		}
	}()

	select {
	case eHTTP := <-errHTTP:
		return eHTTP
	case eGRPC := <-errGRPC:
		return eGRPC
	case <-ctxSignal.Done():
		s.logger.Info("shutting down servers")
		return nil
	}
}

func (s *Server) Stop() error {
	return nil
}

func (s *Server) httpRun(ctx context.Context) error {
	mux := http.NewServeMux()
	for _, svcHTTP := range s.servicesHTTP {
		mux.HandleFunc("/api"+svcHTTP.pattern, svcHTTP.service.SendService(ctx))
	}

	s.logger.Infof("server HTTP run on port %d", s.config.ServiceHTTPPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.ServiceHTTPPort), mux)
	if err != nil {
		s.logger.Infof("httpRun.ListenAndServe: %v", err)
		return fmt.Errorf("httpRun.ListenAndServe: %v", err)
	}

	return nil
}

func (s *Server) grpcRun() error {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.ServiceGRPCPort))
	if err != nil {
		s.logger.Infof("grpcRun.Listen: %v", err)
		return fmt.Errorf("grpcRun.Listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	srv := grpc.NewServer(opts...)

	s.logger.Infof("server GRPC run on port %d", s.config.ServiceGRPCPort)
	err = srv.Serve(l)
	if err != nil {
		s.logger.Infof("grpcRun.Serve: %v", err)
		return fmt.Errorf("grpcRun.Serve: %v", err)
	}

	return nil
}
