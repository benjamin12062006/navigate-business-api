package authhttp

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"navigate-business-api/internal/auth/firebase"
)

type TokenVerifier interface {
	Verify(context.Context, string) (string, error)
}

type uidKey struct{}

func (h *Handler) RequireGoogle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Values("Authorization")) != 1 {
			writeAuthError(w, http.StatusUnauthorized, "invalid_token")
			return
		}
		const prefix = "Bearer "
		value := r.Header.Get("Authorization")
		if !strings.HasPrefix(value, prefix) || len(value) == len(prefix) || strings.TrimSpace(value[len(prefix):]) != value[len(prefix):] {
			writeAuthError(w, http.StatusUnauthorized, "invalid_token")
			return
		}
		uid, err := h.verifier.Verify(r.Context(), value[len(prefix):])
		if err != nil {
			if errors.Is(err, firebaseauth.ErrInvalidToken) {
				writeAuthError(w, http.StatusUnauthorized, "invalid_token")
			} else {
				writeAuthError(w, http.StatusServiceUnavailable, "auth_unavailable")
			}
			return
		}
		if uid == "" {
			writeAuthError(w, http.StatusUnauthorized, "invalid_token")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), uidKey{}, uid)))
	})
}

func verifiedUID(ctx context.Context) string {
	uid, _ := ctx.Value(uidKey{}).(string)
	return uid
}
