# Websites HealthCheck Service

## Setup

Dependencies:

- GO 1.23

## Install dependencies

```bash
go mod tidy
```

## Setup environment

1. Create a `.env` file in the root directory of the project and copy the contents of the `.env.example` file into it. Fill in the necessary values.
2. Create a `sites.json` file in the root directory of the project and copy the contents of the `sites.example.json` file into it. Fill in the necessary values.

## Run app

### Without compiling
```bash
go run cmd/main.go
```

### Or with compiling

```bash
go build -o healthcheck cmd/main.go
./healthcheck
```