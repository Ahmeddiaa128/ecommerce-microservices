# ---------------- Build Stage ----------------
FROM golang:1.25.3-alpine AS builder
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o productservice ./services/ProductService/cmd/main.go

# ---------------- Development Stage ----------------
FROM golang:1.25.3-alpine AS dev
WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy code for hot reload
COPY --from=builder /app /app

# Expose gRPC and HTTP ports
EXPOSE 50053 8080

# Command to run Air
CMD ["/go/bin/air", "-c", "./services/ProductService/.air.toml"]

# ---------------- Production Stage ----------------
FROM gcr.io/distroless/static-debian12:nonroot AS prod
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/productservice .

# Numeric UID 'nonroot'
USER 65532:65532

# Expose gRPC 
EXPOSE 50053

# Command to run the binary
CMD ["./productservice"]

