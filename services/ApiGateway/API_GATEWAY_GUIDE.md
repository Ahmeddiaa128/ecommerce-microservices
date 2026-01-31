# API Gateway Integration Guide

## Overview

The API Gateway serves as the single entry point for all client requests in the e-commerce microservices architecture. It handles routing, authentication, authorization, rate limiting, and request/response transformation.

## Architecture Flow

```
Client Request
      ↓
API Gateway (Port 8080)
      ↓
[Middlewares Stack]
  - CORS
  - Recovery
  - Request ID
  - Logger
  - Timeout
  - Rate Limiter
  - Auth (if required)
  - Role Check (if required)
      ↓
[HTTP Handlers]
      ↓
[gRPC Clients]
      ↓
Microservices
  - User Service (50051)
  - Product Service (50052)
  - Cart Service (50053)
  - Order Service (50054)
```

## Key Components

### 1. Configuration (`config/`)

- **config.go**: Centralized configuration management
- **.env**: Environment variables for deployment

### 2. Clients (`internal/clients/`)

- **grpc_clients.go**: Manages gRPC connections to all microservices
- Connection pooling and health checking
- Automatic reconnection on failure

### 3. Middlewares (`internal/middleware/`)

#### auth.go

- JWT token validation
- User claims extraction
- Role-based access control
- Optional authentication for public endpoints

#### cors.go

- Cross-Origin Resource Sharing
- Configurable origins, methods, and headers
- Preflight request handling
- Panic recovery

#### logger.go

- Request/response logging
- Request ID generation and propagation
- Performance metrics (duration, status, size)

#### ratelimit.go

- Token bucket rate limiting
- Per-IP rate limiting
- Configurable limits and windows
- Rate limit headers in responses

#### timeout.go

- Request timeout protection
- Configurable timeout duration
- Graceful timeout handling

### 4. Handlers (`internal/handlers/`)

#### user_handler.go

**Public Endpoints:**

- Register: User registration
- Login: User authentication

**Authenticated Endpoints:**

- GetProfile: Current user profile
- UpdateUser: Update user information
- Address management (CRUD)

**Admin Only:**

- SearchUsers: List all users
- GetUserByID: Get any user by ID
- DeleteUser: Delete user account

#### product_handler.go

**Public Endpoints:**

- ListProducts: Browse products with pagination
- GetProductByID: Product details
- ListCategories: Browse categories
- GetCategoryByID: Category details

**Admin Only:**

- Product CRUD operations
- Category CRUD operations

#### cart_handler.go

**Authenticated Endpoints:**

- GetCart: View user's cart
- AddItem: Add product to cart
- UpdateItem: Update item quantity
- RemoveItem: Remove product from cart
- ClearCart: Empty cart

#### order_handler.go

**Authenticated Endpoints:**

- CreateOrder: Place new order
- ListOrders: View user's orders
- GetOrderByID: Order details
- AddOrderItem: Add item to order
- RemoveOrderItem: Remove item from order

**Admin Only:**

- UpdateOrderStatus: Update order status

### 5. Router (`internal/router/`)

- Route registration and organization
- Middleware chain configuration
- Method-based routing
- Handler composition

## Request Flow Example

### 1. User Registration

```
POST /api/v1/users/register
↓
CORS → Recovery → RequestID → Logger → Timeout → RateLimiter
↓
UserHandler.Register
↓
UserServiceClient.CreateUser (gRPC)
↓
User Service (gRPC Server)
↓
Response
```

### 2. Authenticated Request (Get Cart)

```
GET /api/v1/cart
Authorization: Bearer <token>
↓
CORS → Recovery → RequestID → Logger → Timeout → RateLimiter
↓
Auth Middleware (validates JWT)
↓
CartHandler.GetCart
↓
CartServiceClient.GetCart (gRPC)
↓
Cart Service (gRPC Server)
↓
Response
```

### 3. Admin Request (Create Product)

```
POST /api/v1/products/create
Authorization: Bearer <admin_token>
↓
CORS → Recovery → RequestID → Logger → Timeout → RateLimiter
↓
Auth Middleware (validates JWT)
↓
Role Middleware (checks admin role)
↓
ProductHandler.CreateProduct
↓
ProductServiceClient.CreateProduct (gRPC)
↓
Product Service (gRPC Server)
↓
Response
```

## Security Features

### JWT Authentication

- Token-based authentication
- Secure token validation
- Claims extraction (user ID, email, role)
- Token expiration handling

### Role-Based Access Control

- Customer role: Basic operations
- Admin role: Management operations
- Flexible role checking
- Multiple role support

### Rate Limiting

- Per-IP rate limiting
- Configurable limits
- Automatic cleanup
- Rate limit headers

### Input Validation

- Request body validation
- Query parameter validation
- Path parameter validation
- Error responses for invalid input

## Performance Optimizations

### Connection Pooling

- Persistent gRPC connections
- Connection reuse
- Reduced latency

### Timeout Management

- Request timeouts prevent resource exhaustion
- Configurable timeouts per operation
- Graceful timeout handling

### Rate Limiting

- Protect services from abuse
- Fair resource distribution
- Configurable per environment

## Error Handling

### Error Response Format

```json
{
  "error": "Unauthorized",
  "message": "invalid or expired token",
  "code": 401
}
```

### Error Types

- **400 Bad Request**: Invalid input
- **401 Unauthorized**: Missing or invalid token
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource not found
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Service error
- **504 Gateway Timeout**: Request timeout

## Monitoring & Observability

### Request ID

Every request gets a unique ID for tracing:

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

### Logging

Structured logs include:

- Request ID
- HTTP method and path
- Status code
- Duration
- Response size

### Rate Limit Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
Retry-After: 60 (when limited)
```

## Configuration Best Practices

### Development

- Lower rate limits for testing
- Verbose logging
- Permissive CORS
- Short JWT expiration

### Production

- Strict rate limits
- Essential logging only
- Specific CORS origins
- Longer JWT expiration
- Strong JWT secret
- HTTPS only

## Deployment

### Docker Compose

```bash
cd services/ApiGateway/docker
docker-compose up -d
```

### Kubernetes

- Horizontal pod autoscaling
- Health check probes
- Resource limits
- Secret management

### Environment Variables

All configuration via environment variables:

- Easy to change per environment
- No code changes needed
- Secret management support

## Testing

### Health Check

```bash
curl http://localhost:8080/health
```

### User Registration

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret123"}'
```

### Authenticated Request

```bash
curl -X GET http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer <token>"
```

## Troubleshooting

### High Latency

- Check gRPC connection health
- Verify service response times
- Review timeout configurations
- Check rate limiting

### Authentication Failures

- Verify JWT secret matches
- Check token expiration
- Validate token format
- Ensure Authorization header format

### Rate Limit Issues

- Adjust rate limit settings
- Check IP extraction logic
- Verify rate limit window

### CORS Errors

- Add client origin to ALLOWED_ORIGINS
- Check preflight request handling
- Verify allowed methods/headers

## Future Enhancements

### Planned Features

- [ ] Request/response caching
- [ ] Circuit breaker pattern
- [ ] Service mesh integration
- [ ] GraphQL support
- [ ] WebSocket support
- [ ] API versioning
- [ ] Request validation middleware
- [ ] Response compression
- [ ] Metrics export (Prometheus)
- [ ] Distributed tracing (Jaeger)
- [ ] API documentation (Swagger)

### Performance Improvements

- [ ] Response caching
- [ ] Request batching
- [ ] Connection pooling optimization
- [ ] Load balancing
- [ ] CDN integration

## Resources

- [Go HTTP Server Best Practices](https://golang.org/doc/)
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Rate Limiting Strategies](https://en.wikipedia.org/wiki/Rate_limiting)
