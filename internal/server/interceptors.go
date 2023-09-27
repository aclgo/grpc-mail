package server

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type InterceptorHTTP interface {
	Logger(ctx context.Context, fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc
}

type InterceptorGRPC interface {
	Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
}
