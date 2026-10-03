.PHONY: migrate-up migrate-down

-include .env
export DATABASE_URL

migrate-up:
	@test -n "$$DATABASE_URL" || { printf '%s\n' 'DATABASE_URL es obligatorio'; exit 1; }
	@migrate -path ./migrations -database "$$DATABASE_URL" up

migrate-down:
	@test -n "$$DATABASE_URL" || { printf '%s\n' 'DATABASE_URL es obligatorio'; exit 1; }
	@migrate -path ./migrations -database "$$DATABASE_URL" down 1
