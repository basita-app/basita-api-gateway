# API Gateway

A high-performance API Gateway built with Go Fiber that routes requests to microservices.

## Features

- **Fiber v2** - Fast HTTP framework
- **Service Routing** - Intelligent request proxying to backend services
- **CORS Support** - Cross-origin resource sharing enabled
- **Logging** - Request/response logging middleware
- **Auto Recovery** - Panic recovery middleware
- **Live Reload** - Air for development hot-reloading

## Architecture

The gateway acts as a single entry point for all client requests and routes them to the appropriate microservices:

```
Client → API Gateway (port 3001) → Microservices
                                 ├─ Auth Service (port 3000)
                                 └─ Car Listing Service (port 3002)
```

## Routes

### Health Check
- `GET /health` - Gateway health status

### Proxied Services
- `/api/auth/*` → Auth Service
- `/api/car-listing/*` → Car Listing Service

## Environment Variables

Create a `.env` file in the root directory:

```env
PORT=3001
AUTH_SERVICE_URL=http://localhost:3000
CAR_LISTING_SERVICE_URL=http://localhost:3002
```

## Installation

```bash
# Install dependencies
go mod download

# Install Air for live reload (optional)
go install github.com/air-verse/air@latest
```

## Development

Run with live reload:

```bash
~/go/bin/air
```

Or run directly:

```bash
go run main.go
```

## Production Build

```bash
# Build the binary
go build -o api-gateway

# Run the binary
./api-gateway
```

## Project Structure

```
api-gateway/
├── main.go           # Application entry point
├── go.mod            # Go module file
├── .air.toml         # Air configuration
├── .gitignore        # Git ignore patterns
└── README.md         # This file
```

## Request Flow

1. Client sends request to API Gateway
2. Gateway matches the request path to a service route
3. Request is proxied to the appropriate microservice with all headers and body
4. Response is returned to the client

All HTTP methods (GET, POST, PUT, DELETE, PATCH, etc.) are supported.
