package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ResponsOK struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

func JSON(span trace.Span, spanMessage string, w http.ResponseWriter, value any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	span.SetStatus(codes.Code(statusCode), spanMessage)
	span.End()

	if err := json.NewEncoder(w).Encode(value); err != nil {
		fmt.Fprint(w, ResponseError{Error: err.Error(), StatusCode: http.StatusInternalServerError})
	}
}
