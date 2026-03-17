# Discussion Forum API

A REST API for a discussion forum built with Go, MySQL, and Redis.

## Architecture

- **Language**: Go 1.21+
- **Database**: MySQL 8.0
- **Cache**: Redis 7
- **Framework**: Gin Web Framework
- **Containerization**: Docker & Docker Compose

## Features

- User management (CRUD operations)
- Topic/thread management
- Post/comment management
- Redis caching for improved performance
- Clean architecture with repository and service layers
- Docker support for easy deployment

## Project Structure

```
.
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routes
│   │   ├── routes.go
│   │   ├── user_handler.go
│   │   ├── topic_handler.go
│   │   └── post_handler.go
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── database/          # Database connections
│   │   ├── mysql.go
│   │   └── redis.go
│   ├── models/            # Data models
│   │   ├── user.go
│   │   ├── topic.go
│   │   └── post.go
│   ├── repository/        # Data access layer
│   │   ├── user_repository.go
│   │   ├── topic_repository.go
│   │   └── post_repository.go
│   └── service/           # Business logic layer
│       ├── user_service.go
│       ├── topic_service.go
│       └── post_service.go
├── migrations/            # Database migrations
│   ├── 01_init.sql
│   └── 02_seed.sql
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Getting Started

### Prerequisites

- Docker and Docker Compose installed
- Go 1.21+ (for local development)

### Running with Docker Compose

1. Clone the repository and navigate to the project directory

2. Start the services:
```bash
docker compose up -d
```

This will start:
- MySQL on port 3306
- Redis on port 6379
- API server on port 8080

3. Wait for services to be healthy (the MySQL health check may take a few seconds)

4. The API will be available at `http://localhost:8080`

### API Endpoints

#### Health Check
- `GET /health` - Check API health

#### Users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### Topics
- `POST /api/v1/topics` - Create a new topic
- `GET /api/v1/topics` - Get all topics
- `GET /api/v1/topics/:id` - Get topic by ID
- `PUT /api/v1/topics/:id` - Update topic
- `DELETE /api/v1/topics/:id` - Delete topic

#### Posts
- `POST /api/v1/posts` - Create a new post
- `GET /api/v1/posts/:id` - Get post by ID
- `GET /api/v1/topics/:topic_id/posts` - Get all posts for a topic
- `PUT /api/v1/posts/:id` - Update post
- `DELETE /api/v1/posts/:id` - Delete post

### Example API Requests

#### Create a user:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "email": "john@example.com"}'
```

#### Get all users:
```bash
curl http://localhost:8080/api/v1/users
```

#### Create a topic:
```bash
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{"title": "My First Topic", "user_id": 1}'
```

#### Create a post:
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"topic_id": 1, "user_id": 1, "content": "This is my first post!"}'
```

#### Get posts for a topic:
```bash
curl http://localhost:8080/api/v1/topics/1/posts
```

## Local Development

### Setup

1. Install Go dependencies:
```bash
go mod download
```

2. Start MySQL and Redis using Docker Compose:
```bash
docker compose up -d mysql redis
```

3. Set environment variables (copy from [.env.example](.env.example)):
```bash
export DATABASE_URL="root:password@tcp(localhost:3306)/discussion_forum?parseTime=true"
export REDIS_URL="localhost:6379"
export PORT="8080"
```

4. Run the application:
```bash
go run main.go
```

### Database Migrations

The migrations are automatically applied when the MySQL container starts. The migration files are located in the `migrations/` directory and are executed in alphabetical order.

To reset the database:
```bash
docker compose down -v
docker compose up -d
```

## Environment Variables

- `DATABASE_URL` - MySQL connection string (default: `root:password@tcp(mysql:3306)/discussion_forum?parseTime=true`)
- `REDIS_URL` - Redis connection address (default: `redis:6379`)
- `PORT` - API server port (default: `8080`)

## Caching Strategy

The API uses Redis for caching:
- User data: 5-minute TTL
- Topic data: 5-minute TTL
- Post data by topic: 3-minute TTL
- User and topic lists: 2-minute TTL

Caches are invalidated on create, update, and delete operations.

## Stopping the Application

```bash
docker compose down
```

To remove all data volumes:
```bash
docker compose down -v
```

## Connect Mysql locally
mysql -h 127.0.0.1 -P 3306 -u root -ppassword discussion_forum

## Next Steps

- Add authentication and authorization
- Implement pagination for list endpoints
- Add search functionality
- Add voting/rating system for posts
- Add user profile pictures
- Implement real-time updates with WebSockets
- Add comprehensive test coverage
