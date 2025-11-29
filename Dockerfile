FROM golang:1.25 AS build

WORKDIR /usr/src/user-notes-api

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /usr/local/bin/user-notes-api ./cmd/main.go

FROM alpine:latest
COPY --from=build /usr/local/bin/user-notes-api /usr/local/bin/user-notes-api

ENTRYPOINT ["/usr/local/bin/user-notes-api"]


