FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./panel-server cmds/panel/main.go

FROM alpine:latest

ARG config='./config.yml'
ENV CONFIG $config
ENV GIN_MODE release

COPY --from=builder /app/panel-server .
COPY "$config" .
# We do this to get a more advanced stacktrace in sentry
COPY . ./app

EXPOSE 3000

CMD ./panel-server -config $CONFIG
