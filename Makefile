BINARY_NAME=pipedrive-challenge
LINTER=golangci-lint

all: test build

build: go build -o pipedrive-challenge cmd/server/main.go

lint: $(LINTER) run