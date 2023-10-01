package service

import (
	"context"

	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/aclgo/grpc-mail/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MailService) SendService(ctx context.Context, req *proto.MailRequest) (*proto.MailResponse, error) {

	data := models.NewMailBody(req.From, req.To, req.Subject, req.Body, req.Template)

	if err := data.Validate(); err != nil {
		s.logger.Errorf("SendService.Validate: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ssss: %v", err)
	}

	err := s.mailUC.Send(data)
	if err != nil {
		s.logger.Errorf("SendService.Send: %v", err)
		return nil, status.Errorf(codes.Internal, "SendService.Send: %v", err)
	}

	return &proto.MailResponse{Message: mail.EmailSentSuccess}, nil
}
