# to silent error on parseJSON ${1} evaluation
1=.
# jq ARE REQUIRED to use this function
parseJSON=$(shell jq --raw-output ${1} secret.json)

POSTGRESQL_USERNAME=$(call parseJSON,.postgreSQL.username)
POSTGRESQL_PASSWORD=$(call parseJSON,.postgreSQL.password)
POSTGRESQL_ADDRESS=$(call parseJSON,.postgreSQL.address)
POSTGRESQL_PORT=$(call parseJSON,.postgreSQL.port)
POSTGRESQL_DATABASE=$(call parseJSON,.postgreSQL.database)
MIGRATION_DIR=

.EXPORT_ALL_VARIABLES:
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=./migrations/
GOOSE_DBSTRING=postgres://${POSTGRESQL_USERNAME}:${POSTGRESQL_PASSWORD}@${POSTGRESQL_ADDRESS}:${POSTGRESQL_PORT}/${POSTGRESQL_DATABASE}

httpserver:
	@CONFIG_FILE=secret.json go run cmd/http/main.go

seeder:
	@CONFIG_FILE=secret.json go run cmd/seed/main.go

devdb:
	@docker run --name postgres --detach \
		--publish ${POSTGRESQL_ADDRESS}:${POSTGRESQL_PORT}:5432 \
		--env POSTGRES_USER=${POSTGRESQL_USERNAME} \
		--env POSTGRES_PASSWORD=${POSTGRESQL_PASSWORD} \
		postgres:17.1-alpine3.20

oci:
	@docker build --tag user:1.0.0 .

test:
	@go test -v -cover ./...

test/k6:
	@k6 run test/script.js

goose/up:
	@goose up

goose/status:
	@goose status

.PHONY: http/server devdb test goose/up oci

