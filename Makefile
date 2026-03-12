DB_URL=postgresql://root:root@localhost:5432/passenger_service?sslmode=disable&x-migrations-table=schema_migrations_order

server:
	go run cmd/api/main.go

dropdb:
	docker exec -it postgres16 dropdb passenger_service

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root passenger_service

migrateup:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)



.PHONY:  server dropdb createdb migrateup migratedown migrateup1 migratedown1 new_migration
