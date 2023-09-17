package main

import (
	"log"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/internal/gmail"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/aclgo/grpc-mail/internal/ses"
	"github.com/aclgo/grpc-mail/pkg/redis"
)

func main() {

	cfg := config.NewConfig()

	redisClient := redis.Connect(cfg)

	_ = redisClient

	ses := ses.NewSes(cfg)
	gmail := gmail.NewGmail(cfg)

	data := models.NewMailBody("from", "to", "subject", "body", "template")

	err := ses.Send(data)
	if err != nil {
		log.Println(err)
	}

	err = gmail.Send(data)
	if err != nil {
		log.Fatalln(err)
	}
}
