FROM golang:latest

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

RUN mkdir "/build"

COPY go.mod .
COPY go.sum .

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

COPY . .

ENTRYPOINT CompileDaemon -build="go build -o /build/app ./cmd/api" -command="/build/app" -polling -graceful-kill