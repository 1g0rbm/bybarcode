#!make
include .env
export $(shell sed 's/=.*//' .env)

migrations-up:
	~/go/bin/migrate -database ${POSTGRESQL_URL} -path migrations up

migrations-down:
	~/go/bin/migrate -database ${POSTGRESQL_URL} -path migrations down 1