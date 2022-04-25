# docker sqlc
dir = $(abspath .)
docker_sqlc = docker run --rm -v ${dir}:/src -w /src kjconroy/sqlc
sqlc_init:
	$(docker_sqlc) init
sqlc:
	$(docker_sqlc) generate


#docker postgre
postgre = docker exec -it postgres14
createdb:
	$(postgre) createdb --username=root --owner=root simple_bank
dropdb:
	$(postgre) dropdb simple_bank


# migrate
db = postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable
migrate = migrate -path db/migration -database $(db)
migrate_up:
	$(migrate) up
migrate_down:
	$(migrate) down


# go
server:
	go run main.go
test:
	go test -v ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc Store