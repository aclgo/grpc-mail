package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

type ErrServiceNotExit struct {
}

func (e ErrServiceNotExit) Error() string {
	return "service not exist"
}

var (
	spanServiceNameFormat     = "send-service-http-%s"
	meterServiceFormatSuccess = "send-service-http-%s-success"
	meterServiceFormatFail    = "send-service-http-%s-fail"
)

func (m *MailService) SendService(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, spanProccessing := m.tel.Tracer().Start(
			context.Background(),
			"send-service-http",
		)

		spanProccessing.AddEvent("request-processing")

		defer spanProccessing.End()

		if r.Method != http.MethodPost {
			spanProccessing.SetStatus(otelCodes.Error, http.StatusText(http.StatusMethodNotAllowed))
			spanProccessing.End()
			respError := ResponseError{
				Error:      http.StatusText(http.StatusMethodNotAllowed),
				StatusCode: http.StatusMethodNotAllowed,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		var data models.MailBody

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			m.logger.Errorf("SendService.json.NewDecoder: %v", err)
			spanProccessing.SetStatus(otelCodes.Error, err.Error())
			spanProccessing.End()
			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		svc, ok := m.svcsMail[data.ServiceName]
		if !ok {
			spanProccessing.SetStatus(otelCodes.Error, ErrServiceNotExit{}.Error())
			spanProccessing.End()

			respError := ResponseError{
				Error:      ErrServiceNotExit{}.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		if err := data.Validate(); err != nil {
			m.logger.Errorf("SendService.Validate: %v", err)
			spanProccessing.SetStatus(otelCodes.Error, err.Error())
			spanProccessing.End()

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		sendSuccess, _ := m.tel.Meter().Int64Counter(
			fmt.Sprintf(meterServiceFormatSuccess, data.ServiceName),
			metric.WithUnit("0"),
		)

		sendFail, _ := m.tel.Meter().Int64Counter(
			fmt.Sprintf(meterServiceFormatFail, data.ServiceName),
			metric.WithUnit("0"),
		)

		_, span := m.tel.Tracer().Start(
			context.Background(),
			fmt.Sprintf(spanServiceNameFormat, data.ServiceName),
		)

		defer span.End()

		span.AddEvent("send-mail")

		if err := svc.Send(&data); err != nil {
			m.logger.Errorf("SendService.Send: %v", err)
			sendFail.Add(context.Background(), 1)
			span.SetStatus(otelCodes.Error, err.Error())

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusInternalServerError,
			}

			JSON(w, respError, respError.StatusCode)

			return
		}

		sendSuccess.Add(context.Background(), 1)
		span.SetStatus(otelCodes.Ok, mail.EmailSentSuccess)

		response := ResponsOK{
			Message: mail.EmailSentSuccess,
		}

		JSON(w, response, http.StatusOK)
	}
}
