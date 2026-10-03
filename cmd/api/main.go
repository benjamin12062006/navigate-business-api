package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"navigate-business-api/internal/app"
	"navigate-business-api/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	application, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("cerrar conector Cloud SQL: %v", err)
		}
	}()

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           application.Handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
	}

	log.Printf("Iniciando API en http://%s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("servidor HTTP: %w", err)
	}
	return nil
}
