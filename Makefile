BINARY_NAME=pipedrive-challenge

build: go build -o pipedrive-challenge cmd/server/main.go

lint: golangci-lint run