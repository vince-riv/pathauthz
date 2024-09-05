## Build
FROM golang:1.21 AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest

WORKDIR /plugins-local/src/pathauthz

COPY --from=build /src .
