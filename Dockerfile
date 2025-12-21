# Multi-stage Dockerfile for fizzbuzz-api

# Builder stage
FROM golang:1.25.1 AS builder
WORKDIR /src

# Download deps first for build cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the repo and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /app/fizzbuzz-api ./cmd/fizzbuzz-api

# Final image
FROM gcr.io/distroless/static-debian11

# Default to binding on all interfaces so the server is reachable from outside the container
ENV FBAPI_HOST=0.0.0.0
ENV FBAPI_PORT=4255
ENV FBAPI_MAX_FIZZBUZZ_LIMIT=100000
ENV FBAPI_MAX_STRING_LENGTH=30
ENV GIN_MODE=release

EXPOSE 4255/tcp

COPY --from=builder /app/fizzbuzz-api /usr/local/bin/fizzbuzz-api

ENTRYPOINT ["/usr/local/bin/fizzbuzz-api"]
