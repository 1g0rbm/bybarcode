#!make
include .env
export $(shell sed 's/=.*//' .env)

db-up:
	docker-compose up -d

run-bot:
	go run cmd/bot/main.go

run-api:
	go run cmd/api/main.go

migrations-up:
	~/go/bin/migrate -database ${POSTGRESQL_URL} -path migrations up

migrations-down:
	~/go/bin/migrate -database ${POSTGRESQL_URL} -path migrations down 1