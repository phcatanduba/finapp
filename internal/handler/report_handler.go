package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"finapp/internal/middleware"
	"finapp/internal/service"
)

type ReportHandler struct {
	svc service.ReportService
}

func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	from, to := parsePeriod(r)

	result, err := h.svc.GetSummary(r.Context(), userID, from, to)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, result)
}

func (h *ReportHandler) ByCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	from, to := parsePeriod(r)

	result, err := h.svc.GetByCategory(r.Context(), userID, from, to)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, result)
}

func (h *ReportHandler) CashFlow(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	from, to := parsePeriod(r)

	result, err := h.svc.GetCashFlow(r.Context(), userID, from, to)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, result)
}

// parsePeriod parses ?from=YYYY-MM-DD&to=YYYY-MM-DD, defaulting to current month.
func parsePeriod(r *http.Request) (time.Time, time.Time) {
	now := time.Now()
	defaultFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	defaultTo := defaultFrom.AddDate(0, 1, -1)

	from := defaultFrom
	to := defaultTo

	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			to = t
		}
	}
	return from, to
}
