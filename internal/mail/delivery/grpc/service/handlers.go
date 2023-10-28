package service

import (
	"context"
	"fmt"

	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/aclgo/grpc-mail/proto"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	spanNameFormat         = "send-service-grpc-%s"
	meterNameFormatSuccess = "send-service-grpc-%s-success"
	meterNameFormatFail    = "send-service-grpc-%s-fail"
	ErrServiceNameNotExist = "service name not exist"
)

func (s *MailService) SendService(ctx context.Context, req *proto.MailRequest) (*proto.MailResponse, error) {

	svc, exist := s.mailUCS[req.Servicename]
	if !exist {
		return nil, status.Error(codes.NotFound, ErrServiceNameNotExist)
	}

	sendSuccess, _ := s.telemetry.Meter().Int64Counter(
		fmt.Sprintf(meterNameFormatSuccess, req.Servicename),
		metric.WithUnit("0"),
	)

	sendFail, _ := s.telemetry.Meter().Float64Counter(
		fmt.Sprintf(meterNameFormatFail, req.Servicename),
		metric.WithUnit("0"),
	)

	_, span := s.telemetry.Tracer().Start(ctx, fmt.Sprintf(spanNameFormat, s.mailUCS[req.Servicename]))
	defer span.End()

	data := models.NewMailBody(req.From, req.To, req.Subject, req.Body, req.Template)

	if err := data.Validate(); err != nil {
		sendFail.Add(context.Background(), 1)
		span.SetStatus(otelCodes.Error, err.Error())
		// s.logger.Errorf("SendService.Validate: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "SendService.Validate: %v", err)
	}

	err := svc.Send(data)
	if err != nil {
		sendFail.Add(context.Background(), 1)
		span.SetStatus(otelCodes.Error, err.Error())
		// s.logger.Errorf("SendService.Send: %v", err)
		return nil, status.Errorf(codes.Internal, "SendService.Send: %v", err)
	}

	sendSuccess.Add(context.Background(), 1)
	span.SetStatus(otelCodes.Error, mail.EmailSentSuccess)

	return &proto.MailResponse{Message: mail.EmailSentSuccess}, nil
}
