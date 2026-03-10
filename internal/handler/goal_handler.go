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

type GoalHandler struct {
	svc service.GoalService
}

func NewGoalHandler(svc service.GoalService) *GoalHandler {
	return &GoalHandler{svc: svc}
}

func (h *GoalHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	goals, err := h.svc.List(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	if goals == nil {
		goals = []model.Goal{}
	}
	WriteJSON(w, http.StatusOK, goals)
}

func (h *GoalHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	var req model.GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	g, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, g)
}

func (h *GoalHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid goal id")
		return
	}
	g, err := h.svc.GetByID(r.Context(), id, userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, g)
}

func (h *GoalHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid goal id")
		return
	}
	var req model.GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	g, err := h.svc.Update(r.Context(), id, userID, req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, g)
}

func (h *GoalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid goal id")
		return
	}
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GoalHandler) Progress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid goal id")
		return
	}
	progress, err := h.svc.GetProgress(r.Context(), id, userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, progress)
}
