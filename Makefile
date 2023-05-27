build:
	@go build -o ./build/SimpleBank

run: build
	@./build/SimpleBank

watch:
	@which gin || go install github.com/codegangsta/gin@latest
	@gin run .


createdb:
	@docker exec -it postgres15 createdb --username=username --owner=username simple_bank

# -U is neccessary to inform what users records to access, if not provided will find user with root dbs
dropdb:
	@docker exec -it postgres15 dropdb -U username simple_bank

migrateup:
	@migrate -path db/migration -database "postgresql://username:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	@migrate -path db/migration -database "postgresql://username:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	@sqlc generate

test:
	@go test -v -cover ./...
	@go clean -testcache