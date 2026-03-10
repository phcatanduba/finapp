package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"finapp/internal/middleware"
	"finapp/internal/model"
	"finapp/internal/service"
)

type PluggyHandler struct {
	svc           service.PluggySyncService
	webhookSecret string
}

func NewPluggyHandler(svc service.PluggySyncService, webhookSecret string) *PluggyHandler {
	return &PluggyHandler{svc: svc, webhookSecret: webhookSecret}
}

func (h *PluggyHandler) GenerateConnectToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	resp, err := h.svc.GenerateConnectToken(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, resp)
}

func (h *PluggyHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	items, err := h.svc.ListItems(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}
	if items == nil {
		items = []model.PluggyItem{}
	}
	WriteJSON(w, http.StatusOK, items)
}

func (h *PluggyHandler) DisconnectItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "invalid item id")
		return
	}

	if err := h.svc.DisconnectItem(r.Context(), id, userID); err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PluggyHandler) Sync(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	var req model.SyncRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.ItemID != "" {
		if err := h.svc.SyncItem(r.Context(), req.ItemID, userID); err != nil {
			WriteError(w, err)
			return
		}
	} else {
		if err := h.svc.SyncAllItems(r.Context(), userID); err != nil {
			WriteError(w, err)
			return
		}
	}

	WriteJSON(w, http.StatusOK, map[string]string{"status": "sync started"})
}

func (h *PluggyHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify HMAC signature if secret is configured
	if h.webhookSecret != "" {
		sig := r.Header.Get("X-Pluggy-Signature")
		if !h.verifySignature(body, sig) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	var payload model.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.svc.HandleWebhook(r.Context(), payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PluggyHandler) verifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
