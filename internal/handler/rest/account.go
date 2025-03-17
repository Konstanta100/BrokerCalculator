package rest

import (
	"encoding/json"
	"github.com/Konstanta100/BrokerCalculator/internal/account"
	"log"
	"net/http"
)

type AccountRequest struct {
	Status string `json:"status,omitempty"`
}

type AccountHandler struct {
	AccountService *account.Service
}

func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var accountRequest AccountRequest
	err := json.NewDecoder(r.Body).Decode(&accountRequest)
	if err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	accounts, err := h.AccountService.GetAccounts(accountRequest.Status)
	if err != nil {
		sendErrorResponse(w, "Error getting accounts", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accounts)
	if err != nil {
		log.Println("[ERROR] Encoding response", err)
		json.NewEncoder(w).Encode(err)
	}
}
