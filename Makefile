GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=pipedrive-challenge
LINTER=golangci-lint

all: test build

test:
      $(GOTEST) ./... -v

build:
      $(GOBUILD) -o $(BINARY_NAME) -v

lint:
      $(LINTER) run