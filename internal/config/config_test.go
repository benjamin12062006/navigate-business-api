package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("FIREBASE_PROJECT_ID", "project-test")
	t.Setenv("INSTANCE_CONNECTION_NAME", "project-test:us-east4:instance-test")
	t.Setenv("DB_NAME", "database-test")
	t.Setenv("DB_USER", "user-test")
	t.Setenv("DB_PASSWORD", "password-test")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "")

	cfg, err := Load()
	if err != nil || cfg.HTTPAddr != "127.0.0.1:8080" || cfg.FirebaseProjectID != "project-test" {
		t.Fatalf("Load() = %+v, %v", cfg, err)
	}

	t.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "127.0.0.1:9099")
	if _, err := Load(); err == nil {
		t.Fatal("emulador permitido fuera de desarrollo")
	}
	t.Setenv("APP_ENV", "development")
	if _, err := Load(); err != nil {
		t.Fatalf("emulador en desarrollo: %v", err)
	}
	t.Setenv("DB_PASSWORD", "")
	if _, err := Load(); err == nil {
		t.Fatal("DB_PASSWORD vacío fue aceptado")
	}
}
