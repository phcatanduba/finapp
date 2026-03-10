package handler

import (
	"net/http"

	"github.com/google/uuid"
	"finapp/internal/middleware"
	"finapp/internal/service"
)

type ProjectionHandler struct {
	svc service.ProjectionService
}

func NewProjectionHandler(svc service.ProjectionService) *ProjectionHandler {
	return &ProjectionHandler{svc: svc}
}

func (h *ProjectionHandler) Balance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	months := parseIntQuery(r.URL.Query().Get("months"), 6)
	if months > 120 {
		months = 120
	}

	result, err := h.svc.ProjectBalance(r.Context(), userID, months)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, result)
}
