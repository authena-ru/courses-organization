.PHONY: openapi
openapi:
	oapi-codegen -generate types -o internal/coursesorg/adapter/delivery/http/v1/openapi_type.gen.go -package v1 api/openapi/courses-organization.yaml
	oapi-codegen -generate chi-server -o internal/coursesorg/adapter/delivery/http/v1/openapi_server.gen.go -package v1 api/openapi/courses-organization.yaml