package rest

import (
	"encoding/json"
	"github.com/Konstanta100/BrokerCalculator/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type AccountHandler struct {
	AccountService *service.AccountService
}

func (h *AccountHandler) Accounts(w http.ResponseWriter, r *http.Request) {
	var (
		userID pgtype.UUID
		err    error
	)

	ctx := r.Context()
	id := r.URL.Query().Get("user_id")
	if id == "" {
		sendErrorResponse(w, "missing user id", http.StatusBadRequest)
		return
	}

	err = userID.Scan(id)
	if err != nil {
		sendErrorResponse(w, "invalid user id", http.StatusBadRequest)
		return
	}

	accounts, err := h.AccountService.FindAccounts(ctx, userID)
	if err != nil {
		sendErrorResponse(w, "Error getting accounts", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accounts)
	if err != nil {
		sendErrorResponse(w, "Error getting accounts", http.StatusInternalServerError)
		return
	}
}

func (h *AccountHandler) Account(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, "Id is required", http.StatusBadRequest)
		return
	}

	account, err := h.AccountService.FindById(ctx, id)
	if account == nil || err != nil {
		sendErrorResponse(w, "Error getting account", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		sendErrorResponse(w, "Error getting accounts", http.StatusInternalServerError)
	}
}

func (h *AccountHandler) LoadAccounts(w http.ResponseWriter, r *http.Request) {
	var (
		userID pgtype.UUID
		err    error
	)

	ctx := r.Context()
	id := r.URL.Query().Get("user_id")
	if id == "" {
		sendErrorResponse(w, "missing user id", http.StatusBadRequest)
		return
	}

	err = userID.Scan(id)
	if err != nil {
		sendErrorResponse(w, "invalid user id", http.StatusBadRequest)
		return
	}

	accounts, err := h.AccountService.LoadAccountsFromBroker(ctx, userID)
	if err != nil {
		sendErrorResponse(w, "Error load accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accounts)
	if err != nil {
		sendErrorResponse(w, "Error load accounts", http.StatusInternalServerError)
	}
}
