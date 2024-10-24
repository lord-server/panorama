FROM docker.io/golang:1.23-alpine AS backend_builder
WORKDIR /app
RUN apk add git make
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile ./
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY static ./static
RUN make

FROM scratch
WORKDIR /app
COPY --from=backend_builder /app/bin/panorama ./
COPY config.example.toml /etc/panorama/config.toml

ENTRYPOINT ["./panorama", "run", "--config", "/etc/panorama/config.toml"]
