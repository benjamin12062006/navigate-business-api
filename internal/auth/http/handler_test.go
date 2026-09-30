package authhttp_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"navigate-business-api/internal/auth"
	"navigate-business-api/internal/auth/firebase"
	authhttp "navigate-business-api/internal/auth/http"
	"navigate-business-api/internal/httpapi"
)

type fakeVerifier struct{ err error }

func (v fakeVerifier) Verify(context.Context, string) (string, error) {
	if v.err != nil {
		return "", v.err
	}
	return "firebase-uid", nil
}

type fakeRepository struct {
	account *auth.Account
	err     error
}

func (r *fakeRepository) Sync(_ context.Context, _ auth.Identity) (auth.Account, error) {
	if r.err != nil {
		return auth.Account{}, r.err
	}
	if r.account == nil {
		r.account = &auth.Account{ID: "account-1", Status: "active"}
		result := *r.account
		result.Created = true
		return result, nil
	}
	return *r.account, nil
}

func (r *fakeRepository) Get(_ context.Context, _ auth.Identity) (auth.Account, error) {
	if r.err != nil {
		return auth.Account{}, r.err
	}
	if r.account == nil {
		return auth.Account{}, auth.ErrNotFound
	}
	return *r.account, nil
}

func request(router http.Handler, method, path, bearer string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if bearer != "" {
		req.Header.Set("Authorization", bearer)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func TestAuthRoutes(t *testing.T) {
	repo := &fakeRepository{}
	router := httpapi.NewRouter(authhttp.NewHandler(fakeVerifier{}, auth.NewService(repo)))
	if got := request(router, http.MethodGet, "/api/v1/auth/me", "Bearer token"); got.Code != http.StatusNotFound {
		t.Fatalf("me before sync = %d", got.Code)
	}
	first := request(router, http.MethodPost, "/api/v1/auth/sync", "Bearer token")
	if first.Code != http.StatusCreated || first.Body.String() != "{\"account_id\":\"account-1\",\"status\":\"active\"}\n" {
		t.Fatalf("first sync = %d %s", first.Code, first.Body.String())
	}
	if got := request(router, http.MethodPost, "/api/v1/auth/sync", "Bearer token"); got.Code != http.StatusOK || got.Body.String() != first.Body.String() {
		t.Fatalf("repeat sync = %d %s", got.Code, got.Body.String())
	}
	if got := request(router, http.MethodGet, "/api/v1/auth/me", "Bearer token"); got.Code != http.StatusOK {
		t.Fatalf("me after sync = %d", got.Code)
	}
	repo.account.Status = "suspended"
	if got := request(router, http.MethodGet, "/api/v1/auth/me", "Bearer token"); got.Code != http.StatusForbidden {
		t.Fatalf("suspended account = %d", got.Code)
	}
}

func TestAuthRejectsMissingAndInvalidTokens(t *testing.T) {
	repo := &fakeRepository{}
	router := httpapi.NewRouter(authhttp.NewHandler(fakeVerifier{}, auth.NewService(repo)))
	for _, bearer := range []string{"", "Basic value", "Bearer "} {
		if got := request(router, http.MethodPost, "/api/v1/auth/sync", bearer); got.Code != http.StatusUnauthorized {
			t.Fatalf("bearer %q = %d", bearer, got.Code)
		}
	}
	if repo.account != nil {
		t.Fatal("una petición sin token creó una cuenta")
	}
	for _, tc := range []struct {
		err  error
		want int
	}{
		{firebaseauth.ErrInvalidToken, http.StatusUnauthorized},
		{errors.New("provider unavailable"), http.StatusServiceUnavailable},
	} {
		router := httpapi.NewRouter(authhttp.NewHandler(fakeVerifier{err: tc.err}, auth.NewService(&fakeRepository{})))
		if got := request(router, http.MethodPost, "/api/v1/auth/sync", "Bearer token"); got.Code != tc.want {
			t.Fatalf("verify error = %d, want %d", got.Code, tc.want)
		}
	}
}

func TestAuthStorageUnavailable(t *testing.T) {
	router := httpapi.NewRouter(authhttp.NewHandler(fakeVerifier{}, auth.NewService(&fakeRepository{err: errors.New("database unavailable")})))
	got := request(router, http.MethodPost, "/api/v1/auth/sync", "Bearer token")
	if got.Code != http.StatusServiceUnavailable || got.Body.String() != "{\"error\":\"storage_unavailable\"}\n" {
		t.Fatalf("storage failure = %d %s", got.Code, got.Body.String())
	}
}
