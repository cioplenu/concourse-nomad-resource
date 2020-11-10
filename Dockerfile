FROM golang:alpine AS build

WORKDIR /tmp/nomad-resource
COPY . .

RUN go build -o dist/out out/main.go

FROM alpine:edge

RUN apk add --no-cache --update nomad

COPY --from=build /tmp/nomad-resource/dist/* /opt/resource/
