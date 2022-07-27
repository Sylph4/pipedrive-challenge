FROM golang:alpine as builder
ADD . /src
WORKDIR /src
RUN go build -o pipedrive-challenge cmd/server/main.go

FROM alpine
WORKDIR /app

COPY --from=builder /src/pipedrive-challenge /app/
ADD  /migrations /app/migrations

ENTRYPOINT /app/pipedrive-challenge