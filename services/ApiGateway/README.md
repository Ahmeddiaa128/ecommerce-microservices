# API Gateway

A production-ready API Gateway for the E-commerce microservices architecture, built with Go and following clean architecture principles.

## Features

### 🔐 Security & Authentication
- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control (RBAC)**: Fine-grained authorization
- **CORS**: Configurable cross-origin resource sharing
- **Rate Limiting**: Protection against abuse with configurable limits

### 🛡️ Middleware Stack
- **Request ID**: Unique identifier for request tracing
- **Logger**: Comprehensive request/response logging
- **Recovery**: Panic recovery for stability
- **Timeout**: Configurable request timeouts
- **CORS**: Cross-origin resource sharing
- **Authentication**: JWT token validation
- **Authorization**: Role-based access control

### 🔌 Service Integration
- **User Service**: User management and authentication
- **Product Service**: Product and category management
- **Cart Service**: Shopping cart operations
- **Order Service**: Order processing and management

### 📊 Observability
- Structured logging with request IDs
- Request/response logging with duration
- Error tracking and reporting

## Architecture

```
ApiGateway/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   ├── config.go            # Configuration management
│   └── .env                 # Environment variables
├── internal/
│   ├── clients/             # gRPC client connections
│   │   └── grpc_clients.go
│   ├── handlers/            # HTTP request handlers
│   │   ├── user_handler.go
│   │   ├── product_handler.go
│   │   ├── cart_handler.go
│   │   ├── order_handler.go
│   │   └── response.go
│   ├── middleware/          # HTTP middlewares
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── ratelimit.go
│   │   └── timeout.go
│   └── router/              # Route configuration
│       └── router.go
└── docker/
    ├── Dockerfile
    ├── Dockerfile.dev
    └── docker-compose.yml
```

## Configuration

### Environment Variables

```env
# Server Configuration
APP_PORT=8080
APP_ENV=development

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
ALLOWED_HEADERS=Accept,Authorization,Content-Type,X-Request-ID

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW_SECONDS=60

# Microservices URLs
USER_SERVICE_URL=localhost:50051
PRODUCT_SERVICE_URL=localhost:50052
CART_SERVICE_URL=localhost:50053
ORDER_SERVICE_URL=localhost:50054

# Timeouts
REQUEST_TIMEOUT_SECONDS=30
IDLE_TIMEOUT_SECONDS=120
READ_TIMEOUT_SECONDS=15
WRITE_TIMEOUT_SECONDS=15
```

## API Endpoints

### Health Check
- `GET /health` - Service health status
- `GET /api/v1/health` - API health status

### User Management

#### Public Endpoints
- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/login` - User login

#### Authenticated Endpoints
- `GET /api/v1/users/profile` - Get current user profile
- `PUT /api/v1/users` - Update user profile

#### Admin Only
- `GET /api/v1/users` - Search users (paginated)
- `GET /api/v1/users/{id}` - Get user by ID
- `DELETE /api/v1/users/{id}` - Delete user

### Address Management (Authenticated)
- `POST /api/v1/addresses` - Create new address
- `GET /api/v1/addresses` - List user addresses
- `PUT /api/v1/addresses/{id}` - Update address
- `DELETE /api/v1/addresses/{id}` - Delete address

### Product Management

#### Public Endpoints
- `GET /api/v1/products` - List products (paginated)
- `GET /api/v1/products/{id}` - Get product details

#### Admin Only
- `POST /api/v1/products` - Create product
- `PUT /api/v1/products/{id}` - Update product
- `DELETE /api/v1/products/{id}` - Delete product

### Category Management

#### Public Endpoints
- `GET /api/v1/categories` - List categories (paginated)
- `GET /api/v1/categories/{id}` - Get category details

#### Admin Only
- `POST /api/v1/categories` - Create category
- `PUT /api/v1/categories/{id}` - Update category
- `DELETE /api/v1/categories/{id}` - Delete category

### Cart Management (Authenticated)
- `GET /api/v1/cart` - Get user cart
- `POST /api/v1/cart/items` - Add item to cart
- `PUT /api/v1/cart/items` - Update cart item quantity
- `DELETE /api/v1/cart/items` - Remove item from cart
- `DELETE /api/v1/cart` - Clear cart

### Order Management (Authenticated)
- `POST /api/v1/orders` - Create new order
- `GET /api/v1/orders` - List user orders (paginated)
- `GET /api/v1/orders/{id}` - Get order details
- `POST /api/v1/orders/items` - Add item to order
- `DELETE /api/v1/orders/items` - Remove item from order

#### Admin Only
- `PATCH /api/v1/orders/status` - Update order status

## Getting Started

### Prerequisites
- Go 1.25.3 or higher
- Docker (optional)
- Running microservices (User, Product, Cart, Order)

### Local Development

1. **Clone the repository**
   ```bash
   cd services/ApiGateway
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp config/.env.example config/.env
   # Edit config/.env with your settings
   ```

4. **Run the service**
   ```bash
   go run cmd/main.go
   ```

### Using Docker

1. **Build and run with Docker Compose**
   ```bash
   cd services/ApiGateway/docker
   docker-compose up --build
   ```

2. **Development mode with hot reload**
   ```bash
   docker-compose -f docker-compose.dev.yml up
   ```

### Using Makefile

```bash
# Run the service
make run

# Build the service
make build

# Run tests
make test

# Run with Docker
make docker-up

# Clean up
make clean
```

## Authentication

### Register a User
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Use the Token
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Best Practices Implemented

### 🏗️ Architecture
- **Clean Architecture**: Separation of concerns with clear boundaries
- **Dependency Injection**: Loose coupling and testability
- **Interface Segregation**: Small, focused interfaces

### 🔒 Security
- **JWT with secure secret**: Token-based authentication
- **Role-based access control**: Fine-grained permissions
- **Rate limiting**: Protection against abuse
- **Input validation**: Request validation
- **Timeout protection**: Prevent long-running requests

### 📝 Code Quality
- **Error handling**: Comprehensive error handling
- **Logging**: Structured logging with context
- **Graceful shutdown**: Clean service termination
- **Recovery middleware**: Panic recovery

### 🚀 Performance
- **Connection pooling**: Efficient gRPC connections
- **Timeouts**: Prevent resource exhaustion
- **Rate limiting**: Control request rate

### 🧪 Maintainability
- **Clear structure**: Easy to navigate codebase
- **Configuration management**: Environment-based config
- **Documentation**: Comprehensive README and comments

## Monitoring & Debugging

### Health Check
```bash
curl http://localhost:8080/health
```

### Request ID Tracing
Every request includes an `X-Request-ID` header for tracing:
```bash
curl -i http://localhost:8080/api/v1/products
# Response includes: X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

### Rate Limit Headers
Rate limit information is included in responses:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests
- `Retry-After`: Seconds until reset (when limited)

## Troubleshooting

### Cannot connect to microservices
- Ensure all microservices are running
- Check service URLs in `.env` file
- Verify network connectivity

### Authentication fails
- Check JWT_SECRET matches across services
- Verify token hasn't expired
- Ensure Authorization header format: `Bearer <token>`

### Rate limit exceeded
- Wait for the rate limit window to reset
- Adjust `RATE_LIMIT_REQUESTS` and `RATE_LIMIT_WINDOW_SECONDS` if needed

## Contributing
1. Follow Go best practices
2. Write tests for new features
3. Update documentation
4. Follow the existing code structure

## License
MIT License
