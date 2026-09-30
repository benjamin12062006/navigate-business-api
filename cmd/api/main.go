package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"navigate-business-api/internal/auth"
	"navigate-business-api/internal/auth/firebase"
	authhttp "navigate-business-api/internal/auth/http"
	authpostgres "navigate-business-api/internal/auth/postgres"
	"navigate-business-api/internal/httpapi"
)

func main() {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	databaseURL := os.Getenv("DATABASE_URL")
	if projectID == "" || databaseURL == "" {
		log.Fatal("FIREBASE_PROJECT_ID y DATABASE_URL son obligatorios")
	}
	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") != "" && os.Getenv("APP_ENV") != "development" {
		log.Fatal("FIREBASE_AUTH_EMULATOR_HOST solo se permite con APP_ENV=development")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	verifier, err := firebaseauth.New(ctx, projectID)
	if err != nil {
		log.Fatalf("no se pudo iniciar Firebase Auth: %v", err)
	}
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal("DATABASE_URL no es válido")
	}
	poolConfig.MaxConns = 5
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("no se pudo crear el pool SQL: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("no se pudo conectar a PostgreSQL: %v", err)
	}

	service := auth.NewService(authpostgres.NewAccountRepository(pool))
	handler := authhttp.NewHandler(verifier, service)
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           httpapi.NewRouter(handler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
	}

	log.Printf("Iniciando API en http://%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
