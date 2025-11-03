## Build
FROM golang:1.25@sha256:6bac879c5b77e0fc9c556a5ed8920e89dab1709bd510a854903509c828f67f96 AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
