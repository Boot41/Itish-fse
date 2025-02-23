# Build stage
FROM golang:1.24-alpine as builder

# Install make and other build essentials
RUN apk add --no-cache make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code and Makefile
COPY . .

# Run make commands for formatting, vetting, and building
RUN make build

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/doctor_ai_assistant .
COPY --from=builder /app/.env .

# Copy migrations directory
COPY --from=builder /app/db/migrations ./db/migrations

# Create necessary directories
RUN mkdir -p /app/pdfs

# Expose port
EXPOSE 8080

# Run the binary directly since we don't need make in production
CMD ["./doctor_ai_assistant"]