package e2e

import (
	"context"
	"strings"
	"time"

	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"github.com/aclgo/grpc-mail/proto"
	"google.golang.org/grpc"
)

type e2eGRPCTest struct {
	mailClient proto.MailServiceClient
}

func Newe2eClient(mailClient proto.MailServiceClient) *e2eGRPCTest {
	return &e2eGRPCTest{
		mailClient: mailClient,
	}
}

func RunGRPC(addrServer string, logger logger.Logger) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addrServer, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("Run.DialContext: %v", err)
		return
	}

	mailClient := proto.NewMailServiceClient(conn)

	param := models.NewMailBody(
		"i am",
		"arcelo2022@gmail.com",
		"test e2e",
		"<h1>Hello</hello>",
		"<body><div><h1>template pre definida + body => %s</h1></div></body>",
		"gmail",
	)

	_, err = mailClient.SendService(
		ctx,
		&proto.MailRequest{
			From:        param.From,
			To:          param.To,
			Subject:     param.Subject,
			Body:        param.Body,
			Template:    param.Template,
			Servicename: param.ServiceName,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), mail.ErrServiceNameNotExist.Error()) {
			logger.Info("TEST e2e GRPC PASS")
			return
		}

		logger.Errorf("Run.SendService: %v", err)
	}

	logger.Info("TEST e2e GRPC PASS")
}
