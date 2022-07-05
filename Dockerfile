FROM golang:1.18-alpine AS builder
WORKDIR /app
RUN apk add git
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY cmd ./cmd
COPY pkg ./pkg
RUN go build -v ./cmd/panorama

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/panorama ./
COPY config.example.toml /etc/panorama/config.toml
COPY static static

ENTRYPOINT ["./panorama", "--config", "/etc/panorama/config.toml"]
