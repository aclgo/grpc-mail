package interceptors

import (
	"context"
	"net/http"
	"time"

	"github.com/aclgo/grpc-mail/pkg/logger"
)

type interceptorHTTP struct {
	logger logger.Logger
}

func NewinterceptorHTTP(logger logger.Logger) *interceptorHTTP {
	return &interceptorHTTP{
		logger: logger,
	}
}

func (i *interceptorHTTP) Logger(ctx context.Context, fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		i.logger.Infof("METHOD: %v, BODY: %v, TIME: %v", r.Method, r.Body, time.Since(start))
	}
}
