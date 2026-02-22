# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src

RUN apk add --no-cache build-base vips-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/goconverter ./cmd/goconverter

FROM alpine:3.22
WORKDIR /

RUN apk add --no-cache vips

COPY --from=builder /out/goconverter /goconverter

ENV PORT=8081
EXPOSE 8081

USER nobody

ENTRYPOINT ["/goconverter"]
