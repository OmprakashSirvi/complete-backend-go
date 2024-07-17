FROM golang:latest

# TODO: Research more on these fields
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

RUN mkdir "/build"

# First copy the two files for installing the dependencies
COPY go.mod .
COPY go.sum .

# Install compile daemon dependency
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

COPY . .

# polling is the one which checks for new file changes..
# graceful kill does not kills the process imediately
ENTRYPOINT CompileDaemon -build="go build -o /build/app ./cmd/api" -command="/build/app" -polling -graceful-kill