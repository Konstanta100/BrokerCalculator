package rest

import (
	"encoding/json"
	"github.com/Konstanta100/BrokerCalculator/internal/operation"
	"log"
	"net/http"
	"time"
)

type OperationsRequest struct {
	DateFrom   string `json:"dateFrom"`
	DateTimeTo string `json:"dateTo,omitempty"`
	Figi       string `json:"figi,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type OperationsResponse struct {
	Operations *operation.Operations `json:"operations"`
}

type OperationHandler struct {
	OperationService *operation.Service
}

func (s *OperationHandler) CalculateCommission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var operationRequest OperationsRequest
	err := json.NewDecoder(r.Body).Decode(&operationRequest)
	if err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if operationRequest.DateFrom == "" {
		sendErrorResponse(w, "dateFrom are required", http.StatusBadRequest)
		return
	}

	dateFrom, err := time.Parse(time.DateOnly, operationRequest.DateFrom)
	if err != nil {
		sendErrorResponse(w, "Invalid dateFrom format", http.StatusBadRequest)
		return
	}

	var dateTo time.Time
	if operationRequest.DateTimeTo == "" {
		dateTo = time.Now().Truncate(24 * time.Hour)
	} else {
		dateTo, err = time.Parse(time.DateOnly, operationRequest.DateTimeTo)
		if err != nil {
			sendErrorResponse(w, "Invalid dateTo format", http.StatusBadRequest)
			return
		}
	}

	calculateCommission, err := s.OperationService.CalculateCommission(dateFrom, dateTo)
	if err != nil {
		sendErrorResponse(w, "Error getting operations", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(calculateCommission)
	if err != nil {
		log.Println("[ERROR] Encoding response", err)
		json.NewEncoder(w).Encode(err)
	}
}
func (s *OperationHandler) GetOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var operationRequest OperationsRequest
	err := json.NewDecoder(r.Body).Decode(&operationRequest)
	if err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if operationRequest.DateFrom == "" {
		sendErrorResponse(w, "dateFrom are required", http.StatusBadRequest)
		return
	}

	dateFrom, err := time.Parse(time.DateOnly, operationRequest.DateFrom)
	if err != nil {
		sendErrorResponse(w, "Invalid dateFrom format", http.StatusBadRequest)
		return
	}

	var dateTo time.Time

	if operationRequest.DateTimeTo == "" {
		dateTo = time.Now()
	} else {
		dateTo, err = time.Parse(time.DateOnly, operationRequest.DateTimeTo)
		if err != nil {
			sendErrorResponse(w, "Invalid dateTo format", http.StatusBadRequest)
			return
		}
	}

	operations, err := s.OperationService.GetOperation(operationRequest.Figi, dateFrom, dateTo)

	if err != nil {
		sendErrorResponse(w, "Error getting operations", http.StatusInternalServerError)
	}

	operationResponse := OperationsResponse{}
	if operations != nil {
		operationResponse.Operations = operations
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(operationResponse)
	if err != nil {
		log.Println("[ERROR] Encoding response", err)
		json.NewEncoder(w).Encode(err)
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}
