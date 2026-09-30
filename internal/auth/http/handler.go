package authhttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"navigate-business-api/internal/auth"
)

type Handler struct {
	verifier TokenVerifier
	service  *auth.Service
}

func NewHandler(verifier TokenVerifier, service *auth.Service) *Handler {
	return &Handler{verifier: verifier, service: service}
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	account, err := h.service.Sync(r.Context(), verifiedUID(r.Context()))
	if writeAccountError(w, err) {
		return
	}
	status := http.StatusOK
	if account.Created {
		status = http.StatusCreated
	}
	writeJSON(w, status, account)
}

func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	account, err := h.service.Current(r.Context(), verifiedUID(r.Context()))
	if writeAccountError(w, err) {
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func writeAccountError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, auth.ErrNotFound):
		writeAuthError(w, http.StatusNotFound, "account_not_synced")
	case errors.Is(err, auth.ErrSuspended):
		writeAuthError(w, http.StatusForbidden, "account_suspended")
	default:
		writeAuthError(w, http.StatusServiceUnavailable, "storage_unavailable")
	}
	return true
}

func writeAuthError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
