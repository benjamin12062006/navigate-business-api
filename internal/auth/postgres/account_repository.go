package authpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"navigate-business-api/internal/auth"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) Get(ctx context.Context, identity auth.Identity) (auth.Account, error) {
	var account auth.Account
	err := r.pool.QueryRow(ctx, `
		SELECT a.id, a.status
		FROM auth_identities AS i
		JOIN accounts AS a ON a.id = i.account_id
		WHERE i.provider = $1 AND i.provider_subject = $2
	`, identity.Provider, identity.Subject).Scan(&account.ID, &account.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.Account{}, auth.ErrNotFound
	}
	if err != nil {
		return auth.Account{}, fmt.Errorf("get account: %w", err)
	}
	return account, nil
}

func (r *AccountRepository) Sync(ctx context.Context, identity auth.Identity) (auth.Account, error) {
	account, err := r.Get(ctx, identity)
	if err == nil {
		return account, nil
	}
	if !errors.Is(err, auth.ErrNotFound) {
		return auth.Account{}, err
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return auth.Account{}, fmt.Errorf("begin account transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `INSERT INTO accounts DEFAULT VALUES RETURNING id, status`).Scan(&account.ID, &account.Status)
	if err != nil {
		return auth.Account{}, fmt.Errorf("create account: %w", err)
	}

	var identityID string
	err = tx.QueryRow(ctx, `
		INSERT INTO auth_identities (account_id, provider, provider_subject)
		VALUES ($1, $2, $3)
		ON CONFLICT (provider, provider_subject) DO NOTHING
		RETURNING id
	`, account.ID, identity.Provider, identity.Subject).Scan(&identityID)
	if errors.Is(err, pgx.ErrNoRows) {
		// The other request won the unique key. Roll back our new account, then read theirs.
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return auth.Account{}, fmt.Errorf("roll back competing account: %w", rollbackErr)
		}
		return r.Get(ctx, identity)
	}
	if err != nil {
		return auth.Account{}, fmt.Errorf("create identity: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return auth.Account{}, fmt.Errorf("commit account: %w", err)
	}
	account.Created = true
	return account, nil
}
