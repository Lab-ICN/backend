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
GOOSE_DBSTRING="postgres://${POSTGRESQL_USERNAME}:${POSTGRESQL_PASSWORD}@${POSTGRESQL_ADDRESS}:${POSTGRESQL_PORT}/${POSTGRESQL_DATABASE}"

http:
	@go run cmd/api/main.go

db:
	@goose up
db/status:
	@goose status
db/migration:
	@goose create ${x} sql

.PHONY:
	run
	test
	db/migration

