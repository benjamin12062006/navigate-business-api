package firebaseauth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	adminauth "firebase.google.com/go/v4/auth"
	"navigate-business-api/internal/auth"
)

type Verifier struct {
	client *adminauth.Client
}

var _ auth.TokenVerifier = (*Verifier)(nil)

func New(ctx context.Context, projectID string) (*Verifier, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return &Verifier{client: client}, nil
}

func (v *Verifier) Verify(ctx context.Context, raw string) (string, error) {
	token, err := v.client.VerifyIDTokenAndCheckRevoked(ctx, raw)
	if err != nil {
		if adminauth.IsIDTokenInvalid(err) || adminauth.IsUserDisabled(err) || adminauth.IsTenantIDMismatch(err) {
			return "", auth.ErrInvalidIdentityToken
		}
		return "", auth.ErrIdentityProviderUnavailable
	}
	return googleUID(token)
}

func googleUID(token *adminauth.Token) (string, error) {
	if token == nil || token.UID == "" || token.Firebase.SignInProvider != "google.com" {
		return "", auth.ErrInvalidIdentityToken
	}
	return token.UID, nil
}
