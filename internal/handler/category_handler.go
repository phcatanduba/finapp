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

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	cats, err := h.svc.List(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	if cats == nil {
		cats = []model.Category{}
	}
	WriteJSON(w, http.StatusOK, cats)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	var req model.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	cat, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, cat)
}

func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid category id")
		return
	}
	cat, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, cat)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid category id")
		return
	}
	var req model.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	cat, err := h.svc.Update(r.Context(), id, userID, req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, cat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid category id")
		return
	}
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
