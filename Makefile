include ./.env

MIGRATION_PATH=db/migrations
DATABASE_URL=postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down

migrate-reset:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down 0

migrate-status:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version

# Seeders
SEED_PATH=db/seeders

# Run all seed files in the correct order (stops on first error)
seed-all:
	@for f in \
		$(SEED_PATH)/000001_seed_payment_methods.sql \
		$(SEED_PATH)/000002_seed_users.sql \
		$(SEED_PATH)/000003_seed_ewallets.sql \
		$(SEED_PATH)/000004_seed_transactions.sql \
		$(SEED_PATH)/000005_seed_transfer_details.sql \
		$(SEED_PATH)/000006_seed_top_up_details.sql; do \
		echo "Seeding $$f..."; \
		psql "$(DATABASE_URL)" -f $$f || exit 1; \
	done

# Run a single seed file: make seed-file file=db/seeders/000001_seed_payment_methods.sql
seed-file:
	@if [ -z "$(file)" ]; then echo "Please provide file variable, e.g. make seed-file file=$(SEED_PATH)/000001_seed_payment_methods.sql"; exit 1; fi
	psql "$(DATABASE_URL)" -f "$(file)"