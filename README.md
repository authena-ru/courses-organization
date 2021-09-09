[![codecov](https://codecov.io/gh/authena-ru/courses-organization/branch/main/graph/badge.svg?token=FYI97SHT3X)](https://codecov.io/gh/authena-ru/courses-organization)

# Authena courses organization service

Courses organization domain for Authena course passing project

## Build & Run (Locally)

### Prerequisites

- go 1.16
- Docker
- golangci-lint

### Environment

You can create .env file with following environment variables or set them manually:

```dotenv
APP_ENVIRONMENT=local # Environment name and config name to parse

MONGO_URI=mongodb://mongodb:27017
MONGO_USERNAME=admin
MONGO_PASSWORD=qwerty
```

### Commands

- ``make openapi`` — generates boilerplate code, types and server interface that conforms to OpenAPI
- ``make build`` — builds project
- ``make lint`` — runs linters
- ``make test-unit`` — runs unit tests and save cover profile
- ``make test-integration`` — runs integration tests and save cover profile
- ``make test-cover`` — builds code cover report
- ``make run-test-db`` — runs Docker with test Mongo DB
- ``make stop-test-db`` — stops Docker with test Mongo DB