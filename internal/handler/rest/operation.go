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
	AccountId  string `json:"accountId,omitempty"`
}

type OperationsResponse struct {
	Operations *operation.Operations `json:"operations"`
}

type OperationHandler struct {
	OperationService *operation.Service
}

func (h *OperationHandler) CalculateCommission(w http.ResponseWriter, r *http.Request) {
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

	calculateCommission, err := h.OperationService.CalculateCommission(operationRequest.AccountId, dateFrom, dateTo)
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
func (h *OperationHandler) GetOperations(w http.ResponseWriter, r *http.Request) {
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

	operations, err := h.OperationService.GetOperation(operationRequest.AccountId, operationRequest.Figi, dateFrom, dateTo)

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
