FROM docker.io/node:18-alpine AS ui_builder
WORKDIR /app/ui
COPY ui/package*.json ./
RUN npm install
COPY ui .
RUN npm run build

FROM docker.io/golang:1.22-alpine AS backend_builder
WORKDIR /app
RUN apk add git
COPY go.mod go.sum main.go ./
RUN go mod download && go mod verify
COPY --from=ui_builder /app/ui/build /app/ui/build
COPY internal ./internal
RUN go build -v

FROM scratch
WORKDIR /app
COPY --from=backend_builder /app/panorama ./
COPY config.example.toml /etc/panorama/config.toml

ENTRYPOINT ["./panorama", "run", "--config", "/etc/panorama/config.toml"]
