package interceptors

import (
	"context"
	"time"

	"github.com/aclgo/grpc-mail/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type interceptorGRPC struct {
	logger logger.Logger
}

func NewinterceptorGRPC(logger logger.Logger) *interceptorGRPC {
	return &interceptorGRPC{
		logger: logger,
	}
}

func (i *interceptorGRPC) Logger(ctx context.Context,
	req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)

	i.logger.Infof("METADATA: %v, METHOD: %v, TIME: %v, ERR: %v", md, info.FullMethod, time.Since(start), err)

	return reply, err
}
