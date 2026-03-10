package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"finapp/internal/middleware"
	"finapp/internal/model"
	"finapp/internal/service"
)

type BudgetHandler struct {
	svc service.BudgetService
}

func NewBudgetHandler(svc service.BudgetService) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	budgets, err := h.svc.List(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	if budgets == nil {
		budgets = []model.Budget{}
	}
	WriteJSON(w, http.StatusOK, budgets)
}

func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	var req model.BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	b, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, b)
}

func (h *BudgetHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid budget id")
		return
	}
	b, err := h.svc.GetByID(r.Context(), id, userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, b)
}

func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid budget id")
		return
	}
	var req model.BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	b, err := h.svc.Update(r.Context(), id, userID, req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, b)
}

func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid budget id")
		return
	}
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BudgetHandler) Progress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid budget id")
		return
	}
	progress, err := h.svc.GetProgress(r.Context(), id, userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, progress)
}
