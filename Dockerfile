## Build
FROM golang:1.26@sha256:595c7847cff97c9a9e76f015083c481d26078f961c9c8dca3923132f51fe12f1 AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN go mod vendor

## Final image
FROM alpine:latest@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

WORKDIR /plugins-local/src/github.com/vince-riv/pathauthz

COPY --from=build /src .
