DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/sigo?sslmode=disable
POSTGRES_URL ?= postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
MIGRATE_IMAGE ?= migrate/migrate:v4.18.3

.PHONY: db-create migrate-up migrate-down migrate-version

db-create:
	@psql "$(POSTGRES_URL)" -v ON_ERROR_STOP=1 -tAc \
		"SELECT 1 FROM pg_database WHERE datname = 'sigo'" | grep -q 1 || \
		psql "$(POSTGRES_URL)" -v ON_ERROR_STOP=1 -c "CREATE DATABASE sigo"

migrate-up:
	@docker run --rm --network host \
		-v "$(CURDIR)/migrations:/migrations:ro" \
		$(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" up

migrate-down:
	@docker run --rm --network host \
		-v "$(CURDIR)/migrations:/migrations:ro" \
		$(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	@docker run --rm --network host \
		-v "$(CURDIR)/migrations:/migrations:ro" \
		$(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" version
