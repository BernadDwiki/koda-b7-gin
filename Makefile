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