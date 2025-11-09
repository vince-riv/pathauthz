## Build
FROM golang:1.23 AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
