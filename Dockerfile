## Build
FROM golang:1.25@sha256:995e25c0e1868fa30a57236d5d8c2252b94b8716e53eae5895cd70dcce532cf0 AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
