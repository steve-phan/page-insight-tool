# Page Insight Tool

A web application that analyzes web pages and extracts HTML structure, headings, links, and login form detection.

## Build and Run Instructions

### Quick Start with Docker Compose

**Prerequisites:** Docker and Docker Compose installed

```bash
# Start all services (Redis, Backend, Frontend)
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build
```

**Services:**

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html (or http://localhost:8080/swagger/)
- Redis: localhost:6379

**Note:** Swagger docs are generated and included. If Swagger doesn't load, regenerate docs with `make swagger` in the backend directory.

**Stop services:**

```bash
docker-compose down
```

### Manual Setup

**Prerequisites:**

- Go 1.21+
- Node.js 20+ and Yarn/npm
- Redis (for rate limiting)

**Backend:**

```bash
cd backend
go mod download

# Start Redis (required)
docker run -d -p 6379:6379 redis:7-alpine

# Run backend
go run cmd/main.go
```

**Frontend:**

```bash
cd frontend
yarn install

# Create .env.local (optional)
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local

# Run frontend
yarn dev
```

## Assumptions and Design Decisions

### Architecture

**Clean Architecture Pattern:**

- Dependency flow: `Config → ServiceFactory → Services → HandlerFactory → Handlers → Routes → Server`
- Dependency injection via factories ensures testability
- **Fail-fast** validation on startup prevents runtime errors

**Modular Extractor Pattern:**

- HTML analysis uses separate extractors (title, headings, links, login forms)
- Each extractor is a testable component following Single Responsibility Principle
- Easy to **extend** with new analysis features without modifying core logic

**Error Handling:**

- Three-layer system: Domain errors → Error mapping → Error middleware
- Centralized error handling with proper HTTP status codes
- Structured error responses with context

**Rate Limiting:**

- **Algorithm:** Fixed Window Counter (simple, efficient)
- **Storage:** Redis-based using sorted sets for distributed rate limiting
- **Configuration:** Per-endpoint limits:
  - `/api/v1/health`: 100 requests/minute
  - `/api/v1/analyze`: 5 requests/10 seconds
- **Headers:** Exposes `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After`
- **Rationale:** Fixed window chosen for simplicity; Token Bucket considered for future if burst handling needed

**SSR Architecture:**

- Next.js App Router for server-side rendering
- Client-side API calls to backend for dynamic updates
- Note: SSR server-to-server requests don't expose headers to browser

**API Documentation (Swagger):**

- **Tool:** Swagger/OpenAPI 2.0 using `swaggo/swag` and `gin-swagger`
- **Location:** Interactive Swagger UI available at `/swagger/index.html`
- **Documentation:** Generated from Go annotations (`@title`, `@version`, `@description`, `@tag`, etc.) in handler files
- **Generation:** Run `swag init` in the backend directory to regenerate docs when API changes
- **Integration:** Swagger docs are auto-registered on server startup via `swag.Register()` in `main.go`
- **Features:** Interactive API testing, request/response schemas, and endpoint documentation
- **Rationale:** Provides developer-friendly API exploration and testing without external tools

### Limitations

1. **Static HTML Analysis Only:** Client-Side Rendered (CSR) sites only return initial HTML. Dynamic content loaded via JavaScript won't be detected.

2. **Protected Sites:** Some sites (e.g., X.com/Twitter) may block automated requests despite rate limiting and proper headers.

3. **Login Form Detection:** Pattern-based detection (input type="password", form elements) may have false positives/negatives depending on implementation.

## Suggestions for Future Improvements

### Immediate Enhancements

1. **Enhanced CSR Support:**

   - Headless browser rendering (Playwright backend option)
   - Anti-detection measures for protected platforms
   - Progressive content loading detection
   - Content quality scoring for dynamic content completeness

2. **Caching Strategy:**

   - Redis caching for frequently analyzed URLs with TTL
   - Content-based cache invalidation
   - Cache hit/miss metrics

3. **Rate Limiting Enhancements:**
   - Token Bucket algorithm option for smoother rate limiting
   - Rate limit configuration via API/admin interface
   - Per-user rate limiting (if authentication added)

### Performance & Scalability

4. **Concurrent Processing:**

   - Batch analysis for multiple URLs
   - Request queuing system
   - Worker pool pattern for analysis tasks

5. **Observability:**
   - BetterStack integration for centralized log aggregation
   - Distributed tracing (OpenTelemetry) for request flow analysis

### Security & Reliability

6. **Security Hardening:**

   - CORS configuration with allowlist origins
   - Security headers (X-Frame-Options, CSP, HSTS)
   - Enhanced input validation and sanitization

### Advanced Features

7. **Analysis Extensions:**

   - SEO analysis extractor (meta tags, structured data)
   - Accessibility analysis (ARIA labels, alt text)
   - Performance analysis (render-blocking resources)
   - Security analysis (mixed content, insecure forms)
