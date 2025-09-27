FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o event-weaver-svc ./cmd

# Stage 2: Minimal runtime image
FROM gcr.io/distroless/base-debian11

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/event-weaver-svc .

# Copy .env (optional)
COPY .env .env

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./event-weaver-svc"]