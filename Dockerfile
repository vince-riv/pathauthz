## Build
FROM golang:1.25@sha256:c0a2bd0756d92462a0d449124b039100ce447ebf69dc6c80a6d877503b36935e AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
