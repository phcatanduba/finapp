package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"finapp/internal/service"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, service.ErrNotFound):
		status = http.StatusNotFound
		msg = "not found"
	case errors.Is(err, service.ErrForbidden):
		status = http.StatusForbidden
		msg = "forbidden"
	case errors.Is(err, service.ErrInvalidCredentials):
		status = http.StatusUnauthorized
		msg = err.Error()
	case errors.Is(err, service.ErrEmailTaken):
		status = http.StatusConflict
		msg = err.Error()
	case errors.Is(err, service.ErrInvalidToken):
		status = http.StatusUnauthorized
		msg = "invalid or expired token"
	default:
		msg = err.Error()
		if status == http.StatusInternalServerError {
			msg = "internal server error"
		}
	}

	WriteJSON(w, status, ErrorResponse{Error: msg})
}

func BadRequest(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: msg})
}
