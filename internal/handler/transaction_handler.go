package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"finapp/internal/middleware"
	"finapp/internal/model"
	"finapp/internal/service"
)

type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	q := r.URL.Query()

	filter := model.TransactionFilter{
		UserID: userID,
		Page:   parseIntQuery(q.Get("page"), 1),
		Limit:  parseIntQuery(q.Get("limit"), 50),
	}

	if v := q.Get("account_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			filter.AccountID = &id
		}
	}
	if v := q.Get("category_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			filter.CategoryID = &id
		}
	}
	if v := q.Get("type"); v != "" {
		t := model.TransactionType(v)
		filter.Type = &t
	}
	if v := q.Get("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filter.From = &t
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filter.To = &t
		}
	}

	resp, err := h.svc.ListTransactions(r.Context(), filter)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, resp)
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid transaction id")
		return
	}

	tx, err := h.svc.GetTransaction(r.Context(), id, userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, tx)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid transaction id")
		return
	}

	var req model.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	tx, err := h.svc.UpdateTransaction(r.Context(), id, userID, req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, tx)
}

func parseIntQuery(v string, fallback int) int {
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
