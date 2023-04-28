create_migration:
# make create_migration name=name_your_migration_without_spaces
	migrate create -ext sql -dir db/migrations -seq ${name}
migrate:
# make migrate password=postgres_password host=localhost port=5420 mode=up/down
	migrate -database 'postgres://postgres:${password}@${host}:${port}/advert_service?sslmode=disable' -path ./schema ${mode}
fmt:
	go fmt ./...
local:
	go build -o . cmd/main.go
	./main --use_db_config
build_image:
	docker build -t danponyavin/pl_advert_service:v1 .
run:
	docker run -d -p 6003:6003 -e POSTGRES_PASSWORD='' \
	-e POSTGRES_HOST='' -e POSTGRES_USER='' \
	-e POSTGRES_PORT='' -e POSTGRES_DB_NAME='' \
	-e GATEWAY_PORT='' -e GATEWAY_IP='' \
	-e ADVERTSERVICE_IP='' -e ADVERTSERVICE_PORT='' \
	--name pet_service_container danponyavin/pl_advert_service:v1
