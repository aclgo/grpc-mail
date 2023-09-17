package ses

import (
	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/internal/models"
)

type Ses struct {
}

func NewSes(config *config.Config) *Ses {
	return &Ses{}
}

func (s *Ses) Connect(data models.MailBody) error {
	return nil
}

func (s *Ses) Send(data *models.MailBody) error {
	return nil
}
