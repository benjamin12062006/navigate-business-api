package authhttp

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"navigate-business-api/internal/auth"
)

type uidKey struct{}

type Authenticator struct {
	verifier auth.TokenVerifier
}

func NewAuthenticator(verifier auth.TokenVerifier) *Authenticator {
	return &Authenticator{verifier: verifier}
}

func (a *Authenticator) RequireGoogle(next http.Handler) http.Handler {
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
		uid, err := a.verifier.Verify(r.Context(), value[len(prefix):])
		if err != nil {
			switch {
			case errors.Is(err, auth.ErrInvalidIdentityToken):
				writeAuthError(w, http.StatusUnauthorized, "invalid_token")
			case errors.Is(err, auth.ErrIdentityProviderUnavailable):
				writeAuthError(w, http.StatusServiceUnavailable, "auth_unavailable")
			default:
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

func VerifiedUID(ctx context.Context) string {
	uid, _ := ctx.Value(uidKey{}).(string)
	return uid
}
