## Build
FROM golang:1.23@sha256:60deed95d3888cc5e4d9ff8a10c54e5edc008c6ae3fba6187be6fb592e19e8c0 AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
