package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Konstanta100/BrokerCalculator/internal/dto"
	"github.com/Konstanta100/BrokerCalculator/internal/service"
)

type OperationsRequest struct {
	DateFrom   string `json:"dateFrom,omitempty"`
	DateTimeTo string `json:"dateTo,omitempty"`
	Figi       string `json:"figi,omitempty"`
	AccountID  string `json:"accountId"`
}

type OperationsResponse struct {
	Operations *dto.Operation `json:"operations"`
}

type OperationHandler struct {
	OperationService *service.OperationService
}

func (h *OperationHandler) CommissionFromBroker(w http.ResponseWriter, r *http.Request) {
	var operationRequest OperationsRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&operationRequest); err != nil {
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

	if operationRequest.AccountID == "" {
		sendErrorResponse(w, "accountId is required", http.StatusBadRequest)
		return
	}

	calculateCommission, err := h.OperationService.CalculateCommission(ctx, operationRequest.AccountID, dateFrom, dateTo)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error getting operations: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(calculateCommission)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error getting operations: %f", err), http.StatusInternalServerError)
	}
}

func (h *OperationHandler) LoadOperations(w http.ResponseWriter, r *http.Request) {
	var operationRequest OperationsRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&operationRequest); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if operationRequest.AccountID == "" {
		sendErrorResponse(w, "accountId is required", http.StatusBadRequest)
		return
	}

	operations, err := h.OperationService.LoadOperationsFromBroker(ctx, operationRequest.AccountID)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error getting operations: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(operations)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error getting operations: %f", err), http.StatusInternalServerError)
	}
}
