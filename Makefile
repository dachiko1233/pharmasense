.PHONY: dev build seed migrate test clean

dev:
	docker-compose up --build

build:
	docker-compose build

migrate:
	cd backend && go run ./cmd/api --migrate-only

seed:
	cd backend && go run ./cmd/api --seed-only

sqlc:
	cd backend && sqlc generate

test:
	cd backend && go test ./...

clean:
	docker-compose down -v
	rm -f backend/pharmasense-api

logs:
	docker-compose logs -f

backend-shell:
	docker exec -it pharmasense-backend sh

db-shell:
	docker exec -it pharmasense-db psql -U pharmasense -d pharmasense
