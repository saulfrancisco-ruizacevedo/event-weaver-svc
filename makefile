.PHONY: all deps wire run clean up down

all: deps wire run

deps:
	go mod tidy

wire:
	wire ./internal/infrastructure/di

run:
	go run ./cmd

clean:
	rm -f ./internal/infrastructure/di/wire_gen.go

up:
	docker-compose up -d --remove-orphans


down:
	docker-compose down

down-all:
	@echo "Stopping containers"
	docker-compose down -v --rmi all
