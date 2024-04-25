include .envrc

# ============================= #
# HELPERS
# ============================= #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ============================= #
# DEVELOPMENT
# ============================= #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -jwt-secret-key=${JWT_SECRET_KEY} -db-dsn=${PAYTM_DB_DSN} -cors-trusted-origins=http://localhost:9000

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${PAYTM_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${PAYTM_DB_DSN} up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${PAYTM_DB_DSN} down

