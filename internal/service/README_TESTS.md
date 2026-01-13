# Service Layer Unit Tests

Complete unit test coverage for all service layers with mocked dependencies for both happy path and unhappy path scenarios.

## Test Files

- [user_service_test.go](user_service_test.go) - Tests for UserService
- [topic_service_test.go](topic_service_test.go) - Tests for TopicService
- [post_service_test.go](post_service_test.go) - Tests for PostService

## Testing Approach

### Fully Mocked Unit Tests
All tests use:
- **Mocked Repositories** - Using `testify/mock` for repository operations
- **In-Memory Redis** - Using `miniredis` for Redis caching (no external Redis required)
- **No Database Required** - Tests run completely isolated without MySQL

This provides true unit testing with no external dependencies.

## Prerequisites

Install the required testing libraries:

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/alicebob/miniredis/v2
```

Install mockery for generating mocks:

```bash
go install github.com/vektra/mockery/v2@latest
```

These dependencies have already been added to the project's `go.mod` file.

## Running the Tests

### Run all service tests
```bash
go test ./internal/service/... -v
```

### Run tests for a specific service
```bash
go test ./internal/service -run TestUserService -v
go test ./internal/service -run TestTopicService -v
go test ./internal/service -run TestPostService -v
```

### Run a specific test
```bash
go test ./internal/service -run TestUserService_Create -v
go test ./internal/service -run TestTopicService_GetByID_WithCache -v
```

### Run with coverage
```bash
go test ./internal/service/... -cover
go test ./internal/service/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Coverage

### UserService Tests (9 test functions, 15+ scenarios)
**Happy Path:**
- ✓ Create user successfully
- ✓ Get user by ID (cache miss)
- ✓ Get user by ID (cache hit - repository not called)
- ✓ Get all users
- ✓ Get all users (cache hit)
- ✓ Get all users with empty list
- ✓ Update user (full update)
- ✓ Update user (partial update - username only)
- ✓ Delete user
- ✓ Cache invalidation on create/update/delete

**Unhappy Path:**
- ✓ Create fails with repository error
- ✓ Get by ID fails when user not found (sql.ErrNoRows)
- ✓ Get by ID fails with database error
- ✓ Get all fails with connection error
- ✓ Update fails when user not found
- ✓ Update fails with database error
- ✓ Delete fails when user not found
- ✓ Delete fails with foreign key constraint

### TopicService Tests (9 test functions, 16+ scenarios)
**Happy Path:**
- ✓ Create topic successfully
- ✓ Get topic by ID (cache miss)
- ✓ Get topic by ID (cache hit)
- ✓ Get all topics
- ✓ Get all topics (cache hit)
- ✓ Get all topics with empty list
- ✓ Update topic
- ✓ Delete topic
- ✓ Cache invalidation on create/update/delete

**Unhappy Path:**
- ✓ Create fails with invalid user (foreign key)
- ✓ Create fails with database connection error
- ✓ Get by ID fails when topic not found
- ✓ Get by ID fails with database error
- ✓ Get all fails with connection error
- ✓ Update fails when topic not found
- ✓ Update fails with database error
- ✓ Delete fails when topic not found
- ✓ Delete fails with foreign key constraint (has posts)

### PostService Tests (7 test functions, 14+ scenarios)
**Happy Path:**
- ✓ Create post successfully
- ✓ Get post by ID
- ✓ Get posts by topic ID (cache miss)
- ✓ Get posts by topic ID (cache hit)
- ✓ Get posts by topic ID with empty list
- ✓ Update post
- ✓ Delete post
- ✓ Cache invalidation on create/update/delete

**Unhappy Path:**
- ✓ Create fails with invalid topic ID
- ✓ Create fails with invalid user ID
- ✓ Create fails with database connection error
- ✓ Get by ID fails when post not found
- ✓ Get by ID fails with database error
- ✓ Get by topic ID fails with connection error
- ✓ Update fails when post not found
- ✓ Update fails with database error
- ✓ Delete fails when post not found
- ✓ Delete fails with database error during deletion

## Architecture

### Mock Strategy

Tests use mockery-generated mocks that implement the repository interfaces. Mocks are automatically generated from the interface definitions and provide type-safe, maintainable test doubles.

#### Mock Generation

Mocks are generated using the mockery tool with configuration in [.mockery.yaml](../../.mockery.yaml):

```yaml
with-expecter: true
all: true
dir: "{{.InterfaceDir}}/mocks"
filename: "mock_{{.InterfaceName}}.go"
mockname: "Mock{{.InterfaceName}}"
outpkg: mocks
packages:
  discussion-forum/internal/service:
    interfaces:
      UserRepositoryInterface:
      TopicRepositoryInterface:
      PostRepositoryInterface:
```

To regenerate mocks when interfaces change:

```bash
mockery
```

Generated mocks are located in [internal/service/mocks/](mocks/):
- [mock_UserRepositoryInterface.go](mocks/mock_UserRepositoryInterface.go)
- [mock_TopicRepositoryInterface.go](mocks/mock_TopicRepositoryInterface.go)
- [mock_PostRepositoryInterface.go](mocks/mock_PostRepositoryInterface.go)

#### Benefits of Mockery-Generated Mocks

1. **Type Safety** - Compile-time verification that mocks match interface signatures
2. **Auto-Generation** - No manual mock maintenance when interfaces change
3. **Better Error Messages** - Panics with clear messages if return values not specified
4. **EXPECT() Helper** - Advanced expectations for complex test scenarios
5. **Consistent Structure** - All mocks follow the same pattern

### In-Memory Redis

Tests use `miniredis` which provides a complete in-memory Redis implementation:

```go
func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
    mr, err := miniredis.Run()
    if err != nil {
        t.Fatalf("Failed to start miniredis: %v", err)
    }

    client := redis.NewClient(&redis.Options{
        Addr: mr.Addr(),
    })

    return client, mr
}
```

Each test creates its own isolated mini redis instance that is cleaned up after the test completes.

### Service Construction for Tests

Services are created using special test constructors that accept interface types:

```go
service := NewUserServiceWithInterface(mockRepo, redisClient)
service := NewTopicServiceWithInterface(mockRepo, redisClient)
service := NewPostServiceWithInterface(mockRepo, redisClient)
```

These constructors allow injecting mock repositories while production code continues using concrete types via the normal constructors (`NewUserService`, etc.).

## Key Testing Patterns

### 1. Table-Driven Tests
Each test uses a table-driven approach for multiple scenarios:

```go
tests := []struct {
    name        string
    request     *models.CreateUserRequest
    mockReturn  *models.User
    mockError   error
    expectError bool
}{
    {
        name: "successful user creation",
        // ... test data
    },
    {
        name: "repository error",
        // ... test data
    },
}
```

### 2. Cache Testing
Tests verify both cache hits and misses:

```go
// First call - populates cache
mockRepo.On("GetByID", int64(1)).Return(user, nil).Once()
result1, _ := service.GetByID(1)

// Second call - hits cache (repository NOT called again)
result2, _ := service.GetByID(1)

// Verify repository was only called once
mockRepo.AssertExpectations(t)
```

### 3. Cache Invalidation Testing
Tests ensure cache is properly cleared on mutating operations:

```go
// Populate cache
service.GetAll()

// Perform mutation (create/update/delete)
service.Create(createReq)

// Verify next call hits repository (cache was cleared)
mockRepo.On("GetAll").Return(users, nil).Once()
service.GetAll()
```

### 4. Error Handling
Tests cover common database errors:
- `sql.ErrNoRows` - Record not found
- Foreign key constraint violations
- Generic database/connection errors

## Production Code Changes

To enable proper unit testing with mocks, the following changes were made to production code:

### 1. Repository Interfaces ([interfaces.go](interfaces.go))
Created interface definitions for all repositories:
- `UserRepositoryInterface`
- `TopicRepositoryInterface`
- `PostRepositoryInterface`

### 2. Service Struct Updates
Services now use interfaces internally:

```go
type UserService struct {
    repo  UserRepositoryInterface  // Was: *repository.UserRepository
    redis *redis.Client
}
```

### 3. Dual Constructors
Each service has two constructors:
- `NewXService(repo *repository.XRepository, ...)` - Production use (unchanged API)
- `NewXServiceWithInterface(repo XRepositoryInterface, ...)` - Testing use

This maintains backward compatibility while enabling dependency injection for tests.

## Benefits of This Approach

1. **No External Dependencies** - Tests run without MySQL or Redis
2. **Fast Execution** - All tests complete in < 1 second
3. **Isolated** - Each test is completely independent
4. **Repeatable** - No flaky tests due to shared state
5. **True Unit Tests** - Only testing the service layer logic
6. **Maintainable** - Clear test structure with table-driven tests
7. **Comprehensive** - Both happy and unhappy paths covered

## Example Test Output

```bash
$ go test ./internal/service/... -v

=== RUN   TestUserService_Create
=== RUN   TestUserService_Create/successful_user_creation
=== RUN   TestUserService_Create/repository_error
--- PASS: TestUserService_Create (0.00s)
    --- PASS: TestUserService_Create/successful_user_creation (0.00s)
    --- PASS: TestUserService_Create/repository_error (0.00s)
=== RUN   TestUserService_GetByID
...
PASS
ok      discussion-forum/internal/service       0.196s
```

## Common Test Patterns

### Testing a Create Operation
```go
mockRepo := mocks.NewMockUserRepositoryInterface(t)
redisClient, mr := setupTestRedis(t)
defer mr.Close()

service := NewUserServiceWithInterface(mockRepo, redisClient)

mockRepo.On("Create", request).Return(mockUser, nil)

result, err := service.Create(request)

assert.NoError(t, err)
assert.Equal(t, mockUser.ID, result.ID)
mockRepo.AssertExpectations(t)
```

### Testing Cache Behavior
```go
// Simulate cache hit by calling twice with Once() expectation
mockRepo.On("GetByID", int64(1)).Return(user, nil).Once()

result1, _ := service.GetByID(1)  // Cache miss
result2, _ := service.GetByID(1)  // Cache hit

// Verify repository was only called once
mockRepo.AssertExpectations(t)
```

### Testing Error Scenarios
```go
mockRepo.On("GetByID", int64(999)).Return(nil, sql.ErrNoRows)

result, err := service.GetByID(999)

assert.Error(t, err)
assert.Nil(t, result)
```

## Future Enhancements

Potential improvements:
- Add benchmark tests for performance testing
- Add integration tests with real database
- Test concurrent operations
- Property-based testing with `gopter`
- Mutation testing to verify test quality

## Mock Maintenance

When repository interfaces change, regenerate the mocks:

```bash
mockery
```

The mockery tool will automatically:
1. Read the interface definitions from [interfaces.go](interfaces.go)
2. Generate updated mock files in [mocks/](mocks/)
3. Ensure mocks match the current interface signatures

**Important**: Always regenerate mocks after modifying repository interfaces to keep tests in sync with production code.

## Notes

- Mock repositories are auto-generated using mockery v2
- Mocks use `testify/mock` which provides method call verification
- Miniredis supports most Redis commands used by the services
- Each test gets a fresh miniredis instance for isolation
- Tests are deterministic and can run in any order
- Production code maintains backward compatibility - existing code using `NewXService` continues to work
- Regenerate mocks with `mockery` when interfaces change
