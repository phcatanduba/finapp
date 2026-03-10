package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"finapp/internal/middleware"
	"finapp/internal/service"
)

type AccountHandler struct {
	svc service.AccountService
}

func NewAccountHandler(svc service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	accounts, err := h.svc.ListAccounts(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, accounts)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid account id")
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), id, userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, acc)
}
