package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponsOK struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

func JSON(w http.ResponseWriter, value any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		fmt.Fprint(w, ResponseError{Error: err.Error(), StatusCode: http.StatusInternalServerError})
	}
}
