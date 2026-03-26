# ===== Build Stage =====
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates bash

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Update Go dependencies
RUN go get -u ./... && go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./services/ApiGateway/cmd/main.go

# ===== Final Stage =====
FROM alpine:3.23.3

# Upgrade all packages to latest stable patch
RUN apk --no-cache upgrade && apk --no-cache add ca-certificates bash zlib

WORKDIR /root/

# Copy binary
COPY --from=builder /app/api-gateway .

# Copy config folder
COPY --from=builder /app/services/ApiGateway/config ./config

# Expose port
EXPOSE 8080


# Run the gateway
CMD ["./api-gateway"]
