package auth

import (
	"context"
	"errors"
)

var (
	ErrInvalidIdentityToken        = errors.New("invalid identity token")
	ErrIdentityProviderUnavailable = errors.New("identity provider unavailable")
)

type TokenVerifier interface {
	Verify(context.Context, string) (string, error)
}

// recibe:
//   1. context.Context → contexto de la petición
//   2. string          → token enviado por el cliente
//
// devuelve:
//   1. string → el UID / Subject de la identidad validada
//   2. error  → error si no se pudo verificar
