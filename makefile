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
	docker build -t rodmul/pl_pet_service:v3 .
run:
	docker run -d -p 6003:6003 -e POSTGRES_PASSWORD='DNd72JDSufesosd9' \
	-e POSTGRES_HOST='79.137.198.139' -e POSTGRES_USER='postgres' \
	-e POSTGRES_PORT='58235' -e POSTGRES_DB_NAME='advert_service' \
	-e GATEWAY_PORT='6002' -e GATEWAY_IP='pl_api_gateway' \
	-e ADVERTSERVICE_IP='127.0.0.1' -e ADVERTSERVICE_PORT='6004' \
	--name pet_service_container rodmul/pl_pet_service:v3