package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aclgo/grpc-mail/internal/mail"

	"github.com/aclgo/grpc-mail/internal/models"
)

func (m *MailService) SendService(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data models.MailBody

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			m.logger.Errorf("SendService.json.NewDecoder: %v", err)

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		if err := data.Validate(); err != nil {
			m.logger.Errorf("SendService.Validate: %v", err)

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		if err := m.mailUC.Send(&data); err != nil {
			m.logger.Errorf("SendService.Send: %v", err)

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusInternalServerError,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		response := ResponsOK{
			Message: mail.EmailSentSuccess,
		}

		JSON(w, response, http.StatusOK)
	}
}
