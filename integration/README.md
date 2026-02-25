# Integration Tests

This directory contains end-to-end integration tests for the Discussion Forum API. These tests use **real MySQL and Redis infrastructure** rather than mocks to ensure actual integration points work correctly.

## Prerequisites

Before running integration tests, ensure you have:

1. **MySQL 8.0+** running on localhost:3306 (or configured via environment variable)
2. **Redis 7+** running on localhost:6379 (or configured via environment variable)

The easiest way to start the infrastructure is with Docker Compose:

```bash
docker-compose up -d mysql redis
```

## Test Database Setup

The tests automatically create the required schema and FULLTEXT index. By default, tests use a separate database (`discussion_forum_test`) to avoid conflicts with development data.

### Creating the Test Database

Connect to MySQL and create the test database:

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS discussion_forum_test;"
```

Or use Docker:

```bash
docker exec -i discussion-forum-mysql mysql -uroot -ppassword -e "CREATE DATABASE IF NOT EXISTS discussion_forum_test;"
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TEST_DATABASE_URL` | `root:password@tcp(localhost:3306)/discussion_forum_test?parseTime=true` | MySQL connection string |
| `TEST_REDIS_URL` | `localhost:6379` | Redis address |

## Running Tests

### Run All Integration Tests

```bash
cd /Users/treasure/Code/discussion
go test ./integration/... -v
```

### Run Specific Test

```bash
go test ./integration/... -v -run TestTopicSearch_HappyPath
```

### Run with Coverage

```bash
go test ./integration/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run with Race Detection

```bash
go test ./integration/... -v -race
```

## Test Categories

### Happy Path Tests
- `TestTopicSearch_HappyPath_FindsRelevantTopics` - Verifies basic search functionality
- `TestTopicSearch_HappyPath_ResultsOrderedByRelevance` - Verifies relevance-based ordering
- `TestTopicSearch_HappyPath_ReturnsCompleteTopicData` - Verifies all fields are returned

### Empty Query Tests
- `TestTopicSearch_EmptyQuery_ReturnsEmptyArray` - Empty `q` parameter
- `TestTopicSearch_MissingQueryParameter_ReturnsEmptyArray` - No `q` parameter
- `TestTopicSearch_WhitespaceOnlyQuery_ReturnsEmptyArray` - Whitespace-only query

### No Results Tests
- `TestTopicSearch_NoMatchingTopics_ReturnsEmptyArray` - Query with no matches
- `TestTopicSearch_EmptyDatabase_ReturnsEmptyArray` - Search on empty database

### Cache Behavior Tests
- `TestTopicSearch_CacheMiss_QueriesDatabase` - Verifies cache population
- `TestTopicSearch_CacheHit_ReturnsCachedResults` - Verifies cache retrieval
- `TestTopicSearch_DifferentQueries_UseSeparateCacheKeys` - Cache key isolation
- `TestTopicSearch_CacheTTL_ExpiresAfterTimeout` - TTL verification (skipped in CI)

### FULLTEXT Search Tests
- `TestTopicSearch_FullText_CaseInsensitiveSearch` - Case insensitivity
- `TestTopicSearch_FullText_PartialWordNotMatched` - Natural language mode behavior
- `TestTopicSearch_FullText_MultipleWords` - Multi-word queries
- `TestTopicSearch_FullText_SpecialCharacters` - Special character handling

### Database Integration Tests
- `TestTopicSearch_Integration_NewTopicBecomesSearchable` - New topics appear in search
- `TestTopicSearch_Integration_DeletedTopicNotSearchable` - Deleted topics removed from search
- `TestTopicSearch_Integration_UpdatedTopicReflectsChanges` - Updated titles are searchable

### Concurrency Tests
- `TestTopicSearch_Concurrent_MultipleSearchesSucceed` - Parallel search requests

## Test Data Isolation

Each test:
1. Creates its own test fixtures
2. Seeds necessary data with unique identifiers
3. Cleans up all data after completion
4. Uses Redis DB 1 to avoid conflicts with development

This ensures tests are:
- Independent of each other
- Can run in any order
- Leave no residual data

## Troubleshooting

### "Failed to connect to MySQL"
- Ensure MySQL is running: `docker-compose ps`
- Check credentials in `TEST_DATABASE_URL`
- Verify the test database exists

### "Failed to connect to Redis"
- Ensure Redis is running: `docker-compose ps`
- Check address in `TEST_REDIS_URL`

### "Failed to create FULLTEXT index"
- Ensure using MySQL 8.0+ with InnoDB engine
- Check that the topics table uses InnoDB: `SHOW CREATE TABLE topics;`

### Tests are slow
- FULLTEXT index creation takes time on first run
- Consider running with `-count=1` to disable test caching
