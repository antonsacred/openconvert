# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/goconverter ./cmd/goconverter

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /

COPY --from=builder /out/goconverter /goconverter

ENV PORT=8081
EXPOSE 8081

ENTRYPOINT ["/goconverter"]
