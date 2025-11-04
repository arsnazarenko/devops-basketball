# Builder stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /app/bin/devops-basketball ./cmd/devops-basketball

# Runtime stage
FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/bin/devops-basketball /app/devops-basketball
COPY ./config ./config

ENV HTTP_HOST=0.0.0.0 \
    HTTP_PORT=8080 \
    METRICS_HOST=0.0.0.0 \
    METRIC_PORT=8081

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/app/devops-basketball"]

