package firebaseauth

import (
	"errors"
	"testing"

	adminauth "firebase.google.com/go/v4/auth"
)

func TestGoogleUID(t *testing.T) {
	for _, tc := range []struct {
		name     string
		token    *adminauth.Token
		wantUID  string
		wantFail bool
	}{
		{"google", &adminauth.Token{UID: "uid-1", Firebase: adminauth.FirebaseInfo{SignInProvider: "google.com"}}, "uid-1", false},
		{"password", &adminauth.Token{UID: "uid-1", Firebase: adminauth.FirebaseInfo{SignInProvider: "password"}}, "", true},
		{"anonymous", &adminauth.Token{UID: "uid-1", Firebase: adminauth.FirebaseInfo{SignInProvider: "anonymous"}}, "", true},
		{"missing_uid", &adminauth.Token{Firebase: adminauth.FirebaseInfo{SignInProvider: "google.com"}}, "", true},
		{"nil", nil, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			uid, err := googleUID(tc.token)
			if uid != tc.wantUID || (tc.wantFail && !errors.Is(err, ErrInvalidToken)) || (!tc.wantFail && err != nil) {
				t.Fatalf("googleUID() = %q, %v", uid, err)
			}
		})
	}
}
