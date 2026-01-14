# Discussion Forum API - Claude Context

## Project Overview
A REST API for a discussion forum built with Go, following clean architecture principles. The application provides CRUD operations for users, topics, and posts with Redis caching for improved performance.

## Tech Stack
- **Language**: Go 1.24.5
- **Web Framework**: Gin
- **Database**: MySQL 8.0
- **Cache**: Redis 7
- **Containerization**: Docker & Docker Compose

## Architecture Pattern
Clean architecture with clear separation of concerns:
```
Handlers (API Layer) → Services (Business Logic) → Repositories (Data Access) → Database
                                ↓
                            Redis Cache
```

## Project Structure
```
/Users/treasure/Code/discussion/
├── main.go                          # Entry point, dependency injection
├── internal/
│   ├── api/                         # HTTP handlers and routing
│   │   ├── routes.go               # Route definitions
│   │   ├── user_handler.go         # User endpoints
│   │   ├── topic_handler.go        # Topic endpoints
│   │   └── post_handler.go         # Post endpoints
│   ├── config/                      # Configuration
│   │   └── config.go               # Environment variable loading
│   ├── database/                    # Database connections
│   │   ├── mysql.go                # MySQL connection
│   │   └── redis.go                # Redis connection
│   ├── models/                      # Domain models and DTOs
│   │   ├── user.go
│   │   ├── topic.go
│   │   └── post.go
│   ├── repository/                  # Data access layer
│   │   ├── user_repository.go
│   │   ├── topic_repository.go
│   │   └── post_repository.go
│   └── service/                     # Business logic layer
│       ├── user_service.go
│       ├── topic_service.go
│       └── post_service.go
├── migrations/                      # SQL migrations
│   ├── 01_init.sql                 # Schema initialization
│   ├── 02_seed.sql                 # Sample data
│   └── 03_add_fulltext_search.sql  # FULLTEXT index for topics
├── docker-compose.yml               # Multi-container setup
├── Dockerfile                       # Application container
└── README.md                        # Project documentation
```

## Running the Application

### Option 1: Docker Compose (Recommended)
```bash
docker-compose up -d
```
This starts MySQL, Redis, and the API server together.

### Option 2: Local Development
1. Start dependencies:
```bash
docker-compose up -d mysql redis
```

2. Set environment variables:
```bash
export DATABASE_URL="root:password@tcp(localhost:3306)/discussion_forum?parseTime=true"
export REDIS_URL="localhost:6379"
export PORT="8080"
```

3. Run the server:
```bash
go run .
```

**Important**: When running locally with `go run .`, MySQL and Redis must be accessible on `localhost:3306` and `localhost:6379` respectively.

## API Routes

### Standard CRUD Pattern
All resources follow RESTful conventions:
- `POST /api/v1/{resource}` - Create
- `GET /api/v1/{resource}` - List all
- `GET /api/v1/{resource}/:id` - Get by ID
- `PUT /api/v1/{resource}/:id` - Update
- `DELETE /api/v1/{resource}/:id` - Delete

### All Endpoints
```
GET    /health                       - Health check
POST   /api/v1/users                 - Create user
GET    /api/v1/users                 - List users
GET    /api/v1/users/:id             - Get user
PUT    /api/v1/users/:id             - Update user
DELETE /api/v1/users/:id             - Delete user
POST   /api/v1/topics                - Create topic
GET    /api/v1/topics                - List topics
GET    /api/v1/topics/search?q=query - Search topics (full-text)
GET    /api/v1/topics/:id            - Get topic
PUT    /api/v1/topics/:id            - Update topic
DELETE /api/v1/topics/:id            - Delete topic
POST   /api/v1/posts                 - Create post
GET    /api/v1/posts/:id             - Get post
GET    /api/v1/topics/:id/posts      - Get posts for a topic
PUT    /api/v1/posts/:id             - Update post
DELETE /api/v1/posts/:id             - Delete post
```

## Coding Conventions

### URL Parameters
- **Always use `:id`** for path parameters, never `:user_id`, `:topic_id`, etc.
- Example: `/api/v1/topics/:id/posts` (correct)
- Counter-example: `/api/v1/topics/:topic_id/posts` (incorrect)
- **Rationale**: Gin router conflicts when mixing different param names at the same path level

### Route Ordering
- List routes (`GET /topics`) must come **before** parameterized routes (`GET /topics/:id`)
- More specific routes before wildcards
- **Rationale**: Gin router matches in order of registration

### Handler Pattern
```go
func (h *Handler) Method(c *gin.Context) {
    // 1. Parse and validate input
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    // 2. Call service layer
    result, err := h.service.Method(id)
    if err != nil {
        // Handle specific errors (sql.ErrNoRows, etc.)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 3. Return response
    c.JSON(http.StatusOK, result)
}
```

### Error Response Format
All errors return JSON with an `error` field:
```json
{"error": "error message here"}
```

### HTTP Status Codes
- `200 OK` - Successful GET, PUT
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Invalid input/parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server errors

## Environment Variables

The application reads from environment variables with fallback defaults in `config.go`:

```go
DATABASE_URL = getEnv("DATABASE_URL", "root:password@tcp(mysql:3306)/discussion_forum?parseTime=true")
REDIS_URL    = getEnv("REDIS_URL", "redis:6379")
PORT         = getEnv("PORT", "8080")
```

**Note**: The app doesn't auto-load `.env` files. Export variables manually when running locally.

## Caching Strategy

Redis is used for caching with TTLs:
- User data: 5 minutes
- Topic data: 5 minutes
- Post data by topic: 3 minutes
- User/topic lists: 2 minutes
- Topic search results: 3 minutes

Cache invalidation happens on:
- Create operations - clear list caches
- Update operations - clear specific item + list caches
- Delete operations - clear specific item + list caches

## Database

### Connection
MySQL connection with `parseTime=true` for automatic time parsing.

### Migrations
Located in `/migrations/` directory:
- `01_init.sql` - Creates tables (users, topics, posts)
- `02_seed.sql` - Inserts sample data
- `03_add_fulltext_search.sql` - Adds FULLTEXT index for topic search

Migrations run automatically when MySQL container starts via docker-entrypoint-initdb.d.

## Full-Text Search

### Overview
Topics support full-text search using MySQL's native FULLTEXT indexing. The search uses natural language mode with relevance ranking.

### Implementation Details
- **Endpoint**: `GET /api/v1/topics/search?q={query}`
- **Query Parameter**: `q` - The search query string
- **Search Mode**: MySQL NATURAL LANGUAGE MODE
- **Indexed Column**: `topics.title`
- **Results**: Ordered by relevance score (most relevant first)
- **Caching**: Search results cached in Redis for 3 minutes

### How It Works
1. Handler extracts query parameter from `?q=` in topic_handler.go:104
2. Service checks Redis cache using key `topics:search:{query}` in topic_service.go:116
3. Repository executes FULLTEXT query with `MATCH...AGAINST` in topic_repository.go:89
4. Results are sorted by relevance score automatically
5. Empty query returns empty array (no results)

### Database Schema
FULLTEXT index is created on the `topics.title` column:
```sql
ALTER TABLE topics ADD FULLTEXT INDEX idx_title_fulltext (title);
```

### Example Queries
```bash
# Search for topics about "Go"
curl "http://localhost:8080/api/v1/topics/search?q=Go"

# Search for topics about "performance optimization"
curl "http://localhost:8080/api/v1/topics/search?q=performance+optimization"

# Empty query returns empty array
curl "http://localhost:8080/api/v1/topics/search?q="
```

### Notes
- Search is case-insensitive
- MySQL's natural language search ignores common words (stopwords)
- Results ranked by relevance, not chronologically
- Minimum word length for indexing is typically 4 characters (MySQL default)
- Route must be registered before `GET /topics/:id` to avoid conflicts (routes.go:27)

## Common Issues & Solutions

### Issue: "Failed to connect to database: dial tcp: lookup mysql: no such host"
**Cause**: Running `go run .` without MySQL on localhost, or environment variables not set.
**Solution**:
- Start MySQL/Redis: `docker-compose up -d mysql redis`
- Export environment variables with `localhost` not `mysql`
- Or use docker-compose to run everything

### Issue: "panic: ':topic_id' in new path conflicts with existing wildcard ':id'"
**Cause**: Route parameter name mismatch in routes.go vs handlers.
**Solution**: Use consistent `:id` parameter names in both routes and handler code.

## Development Workflow

1. Make code changes
2. Stop running server if needed (Ctrl+C or `pkill -f "go run"`)
3. Restart: `go run .` (with env vars) or `docker-compose restart app`
4. Test with curl or API client
5. Check server logs for errors

## Testing Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Get all topics
curl http://localhost:8080/api/v1/topics

# Get posts for a topic
curl http://localhost:8080/api/v1/topics/1/posts

# Search topics
curl "http://localhost:8080/api/v1/topics/search?q=performance"

# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@example.com"}'
```

## Design Principles

1. **Clean Architecture** - Strict layer separation, no business logic in handlers
2. **Dependency Injection** - Dependencies passed via constructors in main.go
3. **Error Handling** - Always check errors, return appropriate HTTP status codes
4. **RESTful Design** - Standard HTTP methods and status codes
5. **Simplicity** - No over-engineering, YAGNI principle

## Future Enhancements (Not Yet Implemented)
- Authentication/Authorization
- Pagination for list endpoints
- Input validation middleware
- Structured logging
- Unit/integration tests
- .env file loading (currently requires manual export)

## Notes for Claude
- When fixing routing issues, always check both routes.go and the corresponding handler
- Parameter names must match between route definition and c.Param() calls
- The project doesn't use godotenv, so env vars must be exported manually for local runs
- Server runs on port 8080 by default
- MySQL is on 3306, Redis on 6379
- Sample data includes 3 users, 3 topics, and several posts
- Full-text search endpoint (`/topics/search`) must be registered before parameterized routes (`/topics/:id`) to avoid routing conflicts
