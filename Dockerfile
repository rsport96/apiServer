FROM golang:1.21.3-alpine

WORKDIR ./src/testService

COPY . .

RUN go get -u github.com/gin-gonic/gin \
github.com/caarlos0/env/v9 \
github.com/gin-gonic/gin \
github.com/pressly/goose/v3

ENTRYPOINT go build -o ./cmd . && ./cmd