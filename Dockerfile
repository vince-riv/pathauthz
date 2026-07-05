## Build
FROM golang:1.25@sha256:f188e8c16ea47a8b22d2bdcf6d9bcd07b63ea7876c199749c07bf31e0ab33bad AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
