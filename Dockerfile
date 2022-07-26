# Inherit the container from a pre-built `golang` image.
FROM golang:latest

# Install the static analysis tool(s)…
RUN go get -u -v github.com/GoASTScanner/gas
RUN go get -u -v golang.org/x/tools/cmd/goimports
RUN go get -u -v github.com/golang/lint/golint

# ...and one tool to rule them all!
RUN go get -u -v https://github.com/golangci/golangci-lint

WORKDIR /go/src/sylph4/pipedrive-challenge
COPY . .