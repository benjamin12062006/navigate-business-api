package authmodule

import (
	"navigate-business-api/internal/auth"
	authhttp "navigate-business-api/internal/auth/http"
)

type Module struct {
	Service       *auth.Service
	Handler       *authhttp.Handler
	Authenticator *authhttp.Authenticator
}

func New(repository auth.Repository, verifier auth.TokenVerifier) *Module {
	service := auth.NewService(repository)
	return &Module{
		Service:       service,
		Handler:       authhttp.NewHandler(service),
		Authenticator: authhttp.NewAuthenticator(verifier),
	}
}
