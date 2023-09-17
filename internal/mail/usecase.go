package mail

import "github.com/aclgo/grpc-mail/internal/models"

type MailUseCase interface {
	Send(*models.MailBody) error
}
