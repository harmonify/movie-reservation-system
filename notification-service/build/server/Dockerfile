FROM golang:1.20.5-alpine3.18 as build
WORKDIR /go/src/github.com/harmonify/movie-reservation-system/notification-service

RUN go install github.com/cosmtrek/air@v1.44.0
RUN mkdir -p .output

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd/server ./cmd/server
COPY ./build/server ./build/server
COPY ./internal ./internal

ENTRYPOINT ["air", "-c", "./build/server/.air.toml"]
