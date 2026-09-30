package auth

import "context"

type Repository interface {
	Sync(context.Context, Identity) (Account, error)
	Get(context.Context, Identity) (Account, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Sync(ctx context.Context, uid string) (Account, error) {
	account, err := s.repository.Sync(ctx, Identity{Provider: ProviderFirebase, Subject: uid})
	if err != nil {
		return Account{}, err
	}
	if account.Status == "suspended" {
		return Account{}, ErrSuspended
	}
	return account, nil
}

func (s *Service) Current(ctx context.Context, uid string) (Account, error) {
	account, err := s.repository.Get(ctx, Identity{Provider: ProviderFirebase, Subject: uid})
	if err != nil {
		return Account{}, err
	}
	if account.Status == "suspended" {
		return Account{}, ErrSuspended
	}
	return account, nil
}
