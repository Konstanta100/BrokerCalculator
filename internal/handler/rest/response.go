package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Println("[ERROR] Encoding response:", message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}
