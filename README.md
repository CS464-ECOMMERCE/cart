# Cart Service

The Cart Service is a gRPC-based microservice that provides shopping cart functionality for the ecommerce application. It uses Redis as the primary storage mechanism for both caching and persistence.

## Features

- **Shopping Cart Operations**:
  - Add items to cart
  - Remove items from cart
  - Update item quantities
  - Empty cart (clear all items)
  - Merge carts (useful for guest user to logged-in user transition)

- **Cart Persistence**:
  - Carts persist across user sessions
  - Redis used as both cache and persistent store
  - Time-to-Live (TTL) mechanism for cart expiration
  - TTL extension on user activity

- **Performance Optimizations**:
  - Redis LRU (Least Recently Used) eviction policy
  - Write-through mechanism for consistency
  - Efficient storage with automatic expiration

## Technical Architecture

### Storage

- **Primary Storage**: Redis
  - Both cache and persistent store
  - TTL: 30 minutes by default (configurable)
  - Eviction Policy: LRU (Least Recently Used)

### Communication

- **gRPC API** for backend communication
- **REST API** for direct client communication
- Protocol Buffers for efficient serialization

### Invalidation Strategy

- Redis entries are deleted on successful checkout
- Abandoned carts expire after TTL
- TTL is extended on any cart activity

## API Reference

### gRPC API

The service exposes the following gRPC methods:

- **AddItem**: Add an item to the cart
- **GetCart**: Retrieve the current cart
- **EmptyCart**: Clear all items from the cart
- **RemoveItem**: Remove a specific item from the cart
- **UpdateItemQuantity**: Update the quantity of an item
- **MergeCart**: Merge a guest cart into a user's cart
- **GetCartTTL**: Get the current TTL for a cart
- **ExtendCartTTL**: Extend the TTL for a cart

### REST API

The service also exposes the following REST endpoints:

- **GET /cart/:user_id**: Get a user's cart
- **POST /cart/:user_id/items**: Add an item to the cart
- **DELETE /cart/:user_id/items/:product_id**: Remove an item from the cart
- **PUT /cart/:user_id/items/:product_id**: Update item quantity
- **DELETE /cart/:user_id**: Empty the cart
- **POST /cart/:user_id/merge/:guest_user_id**: Merge a guest cart into a user's cart
- **GET /cart/:user_id/ttl**: Get the TTL for a cart
- **PUT /cart/:user_id/ttl**: Extend the TTL for a cart
- **GET /health**: Health check endpoint

#### Example REST API Requests

```bash
# Get a user's cart
curl -X GET http://localhost:8082/cart/1

# Add an item to the cart
curl -X POST http://localhost:8082/cart/1/items \
  -H "Content-Type: application/json" \
  -d '{"product_id": 123, "quantity": 2}'

# Update item quantity
curl -X PUT http://localhost:8082/cart/1/items/123 \
  -H "Content-Type: application/json" \
  -d '{"quantity": 5}'

# Remove an item from the cart
curl -X DELETE http://localhost:8082/cart/1/items/123

# Empty the cart
curl -X DELETE http://localhost:8082/cart/1

# Merge a guest cart into a user's cart
curl -X POST http://localhost:8082/cart/1/merge/2

# Get TTL for a cart
curl -X GET http://localhost:8082/cart/1/ttl

# Extend TTL for a cart
curl -X PUT http://localhost:8082/cart/1/ttl \
  -H "Content-Type: application/json" \
  -d '{"ttl_seconds": 3600}'
```

## Environment Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| GRPC_PORT | gRPC server port | 50051 |
| HTTP_PORT | HTTP server port | 8080 |
| REDIS_ADDR | Redis server address | redis:6379 |
| REDIS_PASSWORD | Redis password | redis_password |
| REDIS_DB | Redis database number | 0 |
| REDIS_DEFAULT_TTL | Default TTL for carts | 30m |

## Development Setup

1. Make sure Redis is running:
   ```bash
   docker-compose up redis -d
   ```

2. Run the cart service:
   ```bash
   docker-compose up cart
   ```

3. To run only the gRPC server or only the HTTP server:
   ```bash
   # Run gRPC server only
   docker-compose run cart go run main.go --server=grpc
   
   # Run HTTP server only
   docker-compose run cart go run main.go --server=http
   ```

4. For testing the gRPC API, you can use tools like [grpcurl](https://github.com/fullstorydev/grpcurl) or [BloomRPC](https://github.com/uw-labs/bloomrpc).

## Testing the API

### gRPC API

Using gRPCurl:

```bash
# Add an item to cart
grpcurl -plaintext -d '{"user_id": 1, "item": {"product_id": 123, "quantity": 2}}' localhost:50051 ecommerce.CartService/AddItem

# Get cart
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 ecommerce.CartService/GetCart
```

### REST API

Using curl:

```bash
# Get a user's cart
curl -X GET http://localhost:8082/cart/1

# Add an item to the cart
curl -X POST http://localhost:8082/cart/1/items \
  -H "Content-Type: application/json" \
  -d '{"product_id": 123, "quantity": 2}'
``` 