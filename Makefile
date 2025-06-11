.PHONY: pg-up pg-down up down

PG_UP_CMD = docker-compose up -d postgres
PG_DOWN_CMD = docker-compose down

up:
	go run cmd/api/main.go

pg-up:
	$(PG_UP_CMD)

pg-down:
	$(PG_DOWN_CMD)
