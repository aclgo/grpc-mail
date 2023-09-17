package server

import (
	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	config      *config.Config
	redisClient *redis.Client
	logger      logger.Logger
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() error {
	return nil
}
