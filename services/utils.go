package services

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse структура для ошибки в формате JSON
type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// RespondWithError отправляет JSON-ответ с ошибкой
func RespondWithError(w http.ResponseWriter, statusCode int, errorCode, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	})
}
