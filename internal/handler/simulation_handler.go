package handler

import (
	"encoding/json"
	"net/http"

	"finapp/internal/model"
	"finapp/internal/service"
)

type SimulationHandler struct {
	svc service.SimulationService
}

func NewSimulationHandler(svc service.SimulationService) *SimulationHandler {
	return &SimulationHandler{svc: svc}
}

func (h *SimulationHandler) CompoundInterest(w http.ResponseWriter, r *http.Request) {
	var req model.CompoundInterestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	if req.Principal <= 0 || req.AnnualRate <= 0 || req.Years <= 0 {
		BadRequest(w, "principal, annual_rate, and years must be positive")
		return
	}
	result := h.svc.CompoundInterest(req)
	WriteJSON(w, http.StatusOK, result)
}

func (h *SimulationHandler) Loan(w http.ResponseWriter, r *http.Request) {
	var req model.LoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	if req.Principal <= 0 || req.AnnualRate <= 0 || req.Terms <= 0 {
		BadRequest(w, "principal, annual_rate, and terms must be positive")
		return
	}
	if req.Table != model.LoanTablePrice && req.Table != model.LoanTableSAC {
		req.Table = model.LoanTablePrice
	}
	result := h.svc.Loan(req)
	WriteJSON(w, http.StatusOK, result)
}
