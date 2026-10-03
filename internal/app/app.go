package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"navigate-business-api/internal/auth/firebase"
	"navigate-business-api/internal/auth/module"
	authpostgres "navigate-business-api/internal/auth/postgres"
	"navigate-business-api/internal/config"
	"navigate-business-api/internal/httpapi"
)

type App struct {
	Handler        http.Handler
	pool           *pgxpool.Pool
	cloudSQLDialer *cloudsqlconn.Dialer
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	verifier, err := firebaseauth.New(ctx, cfg.FirebaseProjectID)
	if err != nil {
		return nil, fmt.Errorf("iniciar Firebase Auth: %w", err)
	}
	poolConfig, err := pgxpool.ParseConfig("host=localhost sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("crear configuración PostgreSQL: %w", err)
	}
	poolConfig.ConnConfig.User = cfg.DatabaseUser
	poolConfig.ConnConfig.Password = cfg.DatabasePassword
	poolConfig.ConnConfig.Database = cfg.DatabaseName
	poolConfig.MaxConns = 5

	dialer, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithLazyRefresh())
	if err != nil {
		return nil, fmt.Errorf("iniciar conector Cloud SQL: %w", err)
	}
	poolConfig.ConnConfig.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
		return dialer.Dial(ctx, cfg.InstanceConnectionName, cloudsqlconn.WithPublicIP())
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		_ = dialer.Close()
		return nil, fmt.Errorf("crear pool SQL: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		_ = dialer.Close()
		return nil, fmt.Errorf("conectar a PostgreSQL: %w", err)
	}

	auth := authmodule.New(authpostgres.NewAccountRepository(pool), verifier)
	return &App{
		Handler:        httpapi.NewRouter(auth.Handler, auth.Authenticator),
		pool:           pool,
		cloudSQLDialer: dialer,
	}, nil
}

func (a *App) Close() error {
	a.pool.Close()
	return a.cloudSQLDialer.Close()
}
