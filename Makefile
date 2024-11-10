DB_URL = postgresql://root:secret@localhost:5432/bank_app?sslmode=disable
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root bank_app

dropdb:
	docker exec -it postgres dropdb bank_app

new_migration:
	migrate create -ext sql -seq -dir internal/migrations $(name)

migrate_up1:
	migrate -path internal/migrations -database "$(DB_URL)" -verbose up

migrate_down1:
	migrate -path internal/migrations -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb new_migration migrate_up1 migrate_down1 sqlc
