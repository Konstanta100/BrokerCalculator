package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/Konstanta100/BrokerCalculator/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRequest struct {
	Status string `json:"status,omitempty"`
}

type UserHandler struct {
	UserService *service.UserService
}

func (h *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	var (
		userID pgtype.UUID
		err    error
	)

	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, "missing user id", http.StatusBadRequest)
		return
	}

	err = userID.Scan(id)
	if err != nil {
		sendErrorResponse(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.FindByID(ctx, userID)
	if user == nil || err != nil {
		sendErrorResponse(w, "Error getting account", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		sendErrorResponse(w, "Encoding response", http.StatusInternalServerError)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, "Decoding request", http.StatusBadRequest)
		return
	}

	var params repository.UserCreateParams
	err = json.Unmarshal(body, &params)
	if err != nil {
		sendErrorResponse(w, "Decoding request", http.StatusBadRequest)
		return
	}

	if params.Name == "" || params.Email == "" {
		sendErrorResponse(w, "name and email is required", http.StatusBadRequest)
		return
	}

	uuid, err := h.UserService.CreateUser(ctx, params)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		fmt.Sprintf(`{"id": "%v"}`, uuid),
	))
}
