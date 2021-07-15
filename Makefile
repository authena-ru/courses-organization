.PHONY:

openapi:
	oapi-codegen -generate types -o internal/adapter/delivery/http/v1/openapi_type.gen.go -package v1 api/openapi/courses-organization.yaml
	oapi-codegen -generate chi-server -o internal/adapter/delivery/http/v1/openapi_server.gen.go -package v1 api/openapi/courses-organization.yaml

build:
	go build ./...

test:
	go test -v -race ./...

cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

lint:
	golangci-lint run