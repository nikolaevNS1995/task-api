#########################
# Сборка бинарника
#########################
FROM golang:1.24.2-alpine AS build

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

RUN apk add --no-cache git gcc libc-dev

COPY go.mod go.sum ./
RUN go mod download && go mod tidy

COPY ./ ./

RUN go generate ./... && go fmt ./...

RUN go build -ldflags="-s -w" -o task-api ./cmd/app/main.go

#########################
# Финальный минимальный образ
#########################
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata ca-certificates

RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
 && echo "Europe/Moscow" >  /etc/timezone \
 && echo "hosts: files dns" > /etc/nsswitch.conf

COPY --from=build /app/task-api ./task-api
COPY --from=build /app/migrations/ ./migrations/

RUN chmod +x ./task-api

ENTRYPOINT ["./task-api"]