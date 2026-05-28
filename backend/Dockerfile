# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0
RUN go build -o /app/server ./

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/server /app/server

EXPOSE 8080

ENTRYPOINT ["/app/server"]
