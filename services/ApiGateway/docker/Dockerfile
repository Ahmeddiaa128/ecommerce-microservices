# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./services/ApiGateway/cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/api-gateway .

# Copy config files
COPY --from=builder /app/services/ApiGateway/config/.env ./config/.env

EXPOSE 8080

CMD ["./api-gateway"]
