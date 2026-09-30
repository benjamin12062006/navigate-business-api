.PHONY: migrate-up migrate-down

migrate-up:
	@test -n "$$DATABASE_URL" || { printf '%s\n' 'DATABASE_URL es obligatorio'; exit 1; }
	@migrate -path ./migrations -database "$$DATABASE_URL" up

migrate-down:
	@test -n "$$DATABASE_URL" || { printf '%s\n' 'DATABASE_URL es obligatorio'; exit 1; }
	@migrate -path ./migrations -database "$$DATABASE_URL" down 1
