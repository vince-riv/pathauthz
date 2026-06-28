## Build
FROM golang:1.25@sha256:995e25c0e1868fa30a57236d5d8c2252b94b8716e53eae5895cd70dcce532cf0 AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
