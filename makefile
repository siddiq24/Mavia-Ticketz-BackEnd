include ./.env
DB_URL=postgres://siddiqrm24:bismillah246@localhost:5432/tickitz
DBURLM=$(DB_URL)?sslmode=disable
MIGRATION_PATH=db/migrations
SEED_PATH=db/seeds

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DBURLM) -path $(MIGRATION_PATH) up

insert-seed:
	for file in $$(ls $(SEED_PATH)/*.sql | sort); do \
		psql "$(DBURLM)" -f $$file; \
	done

migrate-down:
	migrate -database $(DBURLM) -path $(MIGRATION_PATH) down
