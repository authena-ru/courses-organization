.PHONY:
.SILENT:

openapi:
	oapi-codegen -generate types -o internal/port/http/v1/openapi_type.gen.go -package v1 api/openapi/courses-organization.yaml
	oapi-codegen -generate chi-server -o internal/port/http/v1/openapi_server.gen.go -package v1 api/openapi/courses-organization.yaml

build:
	go build ./...

lint:
	golangci-lint run

test-unit:
	go test --short -v -race -coverpkg=./... -coverprofile=unit-all.out ./...
	cat unit-all.out | grep -v .gen.go > unit.out
	rm unit-all.out

test-integration:
	make run-test-db
	go test -v -race -coverprofile=integration.out ./internal/adapter/... || make stop-test-db

test-cover:
	go install github.com/wadey/gocovmerge@latest
	gocovmerge unit.out integration.out > cover.out
	go tool cover -html=cover.out -o cover.html


export TEST_DB_URI=mongodb://localhost:27019
export TEST_DB_NAME=test
export TEST_CONTAINER_NAME=test_db

run-test-db:
	docker run --rm -d -p 27019:27017 --name $$TEST_CONTAINER_NAME -e MONGODB_DATABASE=$$TEST_DB_NAME mongo:4.4-bionic

stop-test-db:
	docker stop $$TEST_CONTAINER_NAME