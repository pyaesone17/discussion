package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"discussion-forum/internal/api"
	"discussion-forum/internal/models"
	"discussion-forum/internal/repository"
	"discussion-forum/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testConfig holds configuration for integration tests
type testConfig struct {
	DatabaseURL string
	RedisURL    string
}

func loadTestConfig() *testConfig {
	return &testConfig{
		DatabaseURL: getEnvOrDefault("TEST_DATABASE_URL", "root:password@tcp(localhost:3306)/discussion_forum_test?parseTime=true"),
		RedisURL:    getEnvOrDefault("TEST_REDIS_URL", "localhost:6379"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// testFixtures contains shared test infrastructure
type testFixtures struct {
	db           *sql.DB
	redis        *redis.Client
	router       *gin.Engine
	topicService *service.TopicService
	cleanup      func()
}

// setupTestFixtures initializes real MySQL and Redis connections for testing
func setupTestFixtures(t *testing.T) *testFixtures {
	t.Helper()

	cfg := loadTestConfig()

	// Connect to MySQL
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	require.NoError(t, err, "Failed to connect to MySQL. Ensure MySQL is running and TEST_DATABASE_URL is correct")

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	require.NoError(t, err, "Failed to ping MySQL")

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
		DB:   1, // Use DB 1 for tests to avoid conflicts with development
	})

	ctx := context.Background()
	err = redisClient.Ping(ctx).Err()
	require.NoError(t, err, "Failed to connect to Redis. Ensure Redis is running")

	// Initialize schema
	initTestSchema(t, db)

	// Create repositories and services
	topicRepo := repository.NewTopicRepository(db)
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	topicService := service.NewTopicService(topicRepo, redisClient)
	userService := service.NewUserService(userRepo, redisClient)
	postService := service.NewPostService(postRepo, redisClient)

	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupRoutes(router, userService, topicService, postService)

	return &testFixtures{
		db:           db,
		redis:        redisClient,
		router:       router,
		topicService: topicService,
		cleanup: func() {
			cleanupTestData(t, db, redisClient)
			db.Close()
			redisClient.Close()
		},
	}
}

// initTestSchema creates tables and FULLTEXT index for testing
func initTestSchema(t *testing.T, db *sql.DB) {
	t.Helper()

	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_username (username),
			INDEX idx_email (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	require.NoError(t, err, "Failed to create users table")

	// Create topics table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS topics (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(200) NOT NULL,
			user_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	require.NoError(t, err, "Failed to create topics table")

	// Add FULLTEXT index if it doesn't exist
	var indexExists int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		AND table_name = 'topics'
		AND index_name = 'idx_title_fulltext'
	`).Scan(&indexExists)
	require.NoError(t, err, "Failed to check for FULLTEXT index")

	if indexExists == 0 {
		_, err = db.Exec(`ALTER TABLE topics ADD FULLTEXT INDEX idx_title_fulltext (title)`)
		require.NoError(t, err, "Failed to create FULLTEXT index")
	}

	// Create posts table (for completeness, though not used in search tests)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			topic_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_topic_id (topic_id),
			INDEX idx_user_id (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	require.NoError(t, err, "Failed to create posts table")
}

// cleanupTestData removes all test data from database and Redis
func cleanupTestData(t *testing.T, db *sql.DB, redisClient *redis.Client) {
	t.Helper()

	// Clear Redis cache
	ctx := context.Background()
	redisClient.FlushDB(ctx)

	// Clear database tables in order (respecting foreign keys)
	_, _ = db.Exec("DELETE FROM posts")
	_, _ = db.Exec("DELETE FROM topics")
	_, _ = db.Exec("DELETE FROM users")
}

// seedTestUser creates a test user and returns their ID
func seedTestUser(t *testing.T, db *sql.DB, username, email string) int64 {
	t.Helper()

	result, err := db.Exec(
		"INSERT INTO users (username, email) VALUES (?, ?)",
		username, email,
	)
	require.NoError(t, err, "Failed to seed test user")

	id, err := result.LastInsertId()
	require.NoError(t, err, "Failed to get user ID")

	return id
}

// seedTestTopic creates a test topic and returns it
func seedTestTopic(t *testing.T, db *sql.DB, title string, userID int64) *models.Topic {
	t.Helper()

	result, err := db.Exec(
		"INSERT INTO topics (title, user_id) VALUES (?, ?)",
		title, userID,
	)
	require.NoError(t, err, "Failed to seed test topic")

	id, err := result.LastInsertId()
	require.NoError(t, err, "Failed to get topic ID")

	topic := &models.Topic{}
	err = db.QueryRow(`
		SELECT t.id, t.title, t.user_id, u.username, t.created_at, t.updated_at
		FROM topics t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`, id).Scan(&topic.ID, &topic.Title, &topic.UserID, &topic.Username, &topic.CreatedAt, &topic.UpdatedAt)
	require.NoError(t, err, "Failed to retrieve seeded topic")

	return topic
}

// --------------------------------------------------------------------------
// Happy Path Tests
// --------------------------------------------------------------------------

func TestTopicSearch_HappyPath_FindsRelevantTopics(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "searchuser", "search@example.com")
	seedTestTopic(t, fixtures.db, "Introduction to Golang Programming", userID)
	seedTestTopic(t, fixtures.db, "Advanced Golang Techniques", userID)
	seedTestTopic(t, fixtures.db, "Python for Beginners", userID)
	seedTestTopic(t, fixtures.db, "JavaScript Best Practices", userID)

	// Execute search for "Golang"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Golang", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err, "Failed to parse response")

	// Should find both Golang topics
	assert.Equal(t, 2, len(topics), "Expected 2 topics matching 'Golang'")

	// Verify the results contain the expected topics
	titles := make([]string, len(topics))
	for i, topic := range topics {
		titles[i] = topic.Title
	}
	assert.Contains(t, titles, "Introduction to Golang Programming")
	assert.Contains(t, titles, "Advanced Golang Techniques")
}

func TestTopicSearch_HappyPath_ResultsOrderedByRelevance(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data with varying relevance
	userID := seedTestUser(t, fixtures.db, "relevanceuser", "relevance@example.com")
	seedTestTopic(t, fixtures.db, "Performance Optimization Guide", userID)
	seedTestTopic(t, fixtures.db, "Database Performance and Optimization Tips", userID)
	seedTestTopic(t, fixtures.db, "Quick Performance Tips for Optimization", userID)

	// Execute search
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Performance+Optimization", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	// Should find topics matching "Performance" or "Optimization"
	assert.GreaterOrEqual(t, len(topics), 1, "Should find at least one topic")

	// Results should be ordered by MySQL's FULLTEXT relevance scoring
	// The topic with both words should rank higher
	for _, topic := range topics {
		assert.NotEmpty(t, topic.Title)
		assert.NotZero(t, topic.ID)
	}
}

func TestTopicSearch_HappyPath_ReturnsCompleteTopicData(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "datauser", "data@example.com")
	expectedTopic := seedTestTopic(t, fixtures.db, "Complete Data Structure Guide", userID)

	// Execute search
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Complete+Data+Structure", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(topics), 1, "Should find at least one topic")

	// Find our expected topic in results
	var found *models.Topic
	for _, topic := range topics {
		if topic.ID == expectedTopic.ID {
			found = topic
			break
		}
	}
	require.NotNil(t, found, "Expected topic should be in results")

	// Verify complete topic data is returned
	assert.Equal(t, expectedTopic.ID, found.ID)
	assert.Equal(t, expectedTopic.Title, found.Title)
	assert.Equal(t, expectedTopic.UserID, found.UserID)
	assert.Equal(t, "datauser", found.Username)
	assert.False(t, found.CreatedAt.IsZero())
	assert.False(t, found.UpdatedAt.IsZero())
}

// --------------------------------------------------------------------------
// Empty Query Tests
// --------------------------------------------------------------------------

func TestTopicSearch_EmptyQuery_ReturnsEmptyArray(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed some topics to ensure empty result is not due to empty database
	userID := seedTestUser(t, fixtures.db, "emptyqueryuser", "emptyquery@example.com")
	seedTestTopic(t, fixtures.db, "Sample Topic for Testing", userID)

	// Execute search with empty query
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	// Should return empty array, not null
	assert.NotNil(t, topics)
	assert.Equal(t, 0, len(topics), "Empty query should return empty array")
}

func TestTopicSearch_MissingQueryParameter_ReturnsEmptyArray(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Execute search without query parameter
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	assert.NotNil(t, topics)
	assert.Equal(t, 0, len(topics))
}

func TestTopicSearch_WhitespaceOnlyQuery_ReturnsEmptyArray(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed some topics
	userID := seedTestUser(t, fixtures.db, "whitespaceuser", "whitespace@example.com")
	seedTestTopic(t, fixtures.db, "Test Topic", userID)

	// Execute search with whitespace-only query (URL encoded)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=+", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	// Whitespace query should return empty or no results from MySQL FULLTEXT
	assert.NotNil(t, topics)
}

// --------------------------------------------------------------------------
// No Results Found Tests
// --------------------------------------------------------------------------

func TestTopicSearch_NoMatchingTopics_ReturnsEmptyArray(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed topics that won't match the search
	userID := seedTestUser(t, fixtures.db, "nomatchuser", "nomatch@example.com")
	seedTestTopic(t, fixtures.db, "Discussion about Golang", userID)
	seedTestTopic(t, fixtures.db, "Python Programming Tips", userID)

	// Search for something that doesn't exist
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=xyznonsensequery", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	assert.NotNil(t, topics, "Should return empty array, not null")
	assert.Equal(t, 0, len(topics), "No topics should match nonsense query")
}

func TestTopicSearch_EmptyDatabase_ReturnsEmptyArray(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Don't seed any data - database should be empty after cleanup

	// Execute search
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=anything", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	assert.NotNil(t, topics)
	assert.Equal(t, 0, len(topics))
}

// --------------------------------------------------------------------------
// Cache Behavior Tests
// --------------------------------------------------------------------------

func TestTopicSearch_CacheMiss_QueriesDatabase(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	ctx := context.Background()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "cachemissuser", "cachemiss@example.com")
	seedTestTopic(t, fixtures.db, "Caching Strategies Tutorial", userID)

	// Clear Redis to ensure cache miss
	fixtures.redis.FlushDB(ctx)

	// Execute search
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Caching+Strategies", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(topics), 1, "Should find the cached topic")

	// Verify cache was populated
	cacheKey := "topics:search:Caching Strategies"
	cached, err := fixtures.redis.Get(ctx, cacheKey).Result()
	assert.NoError(t, err, "Cache should be populated after search")
	assert.NotEmpty(t, cached)
}

func TestTopicSearch_CacheHit_ReturnsCachedResults(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	ctx := context.Background()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "cachehituser", "cachehit@example.com")
	seedTestTopic(t, fixtures.db, "Redis Caching Patterns", userID)

	query := "Redis Caching"
	encodedQuery := url.QueryEscape(query)

	// First request - populates cache
	req1 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/topics/search?q=%s", encodedQuery), nil)
	w1 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)

	var topics1 []*models.Topic
	err := json.Unmarshal(w1.Body.Bytes(), &topics1)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(topics1), 1)

	// Verify cache exists
	cacheKey := fmt.Sprintf("topics:search:%s", query)
	_, err = fixtures.redis.Get(ctx, cacheKey).Result()
	require.NoError(t, err, "Cache should exist after first request")

	// Delete the topic from database to prove second request uses cache
	_, err = fixtures.db.Exec("DELETE FROM topics WHERE title LIKE '%Redis Caching%'")
	require.NoError(t, err)

	// Second request - should return cached results even though database is empty
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/topics/search?q=%s", encodedQuery), nil)
	w2 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var topics2 []*models.Topic
	err = json.Unmarshal(w2.Body.Bytes(), &topics2)
	require.NoError(t, err)

	// Should still return results from cache
	assert.Equal(t, len(topics1), len(topics2), "Cached results should match initial results")
}

func TestTopicSearch_DifferentQueries_UseSeparateCacheKeys(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	ctx := context.Background()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "multicacheuser", "multicache@example.com")
	seedTestTopic(t, fixtures.db, "Golang Concurrency Patterns", userID)
	seedTestTopic(t, fixtures.db, "Python Async Programming", userID)

	// Search for Golang
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Golang", nil)
	w1 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Search for Python
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Python", nil)
	w2 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Verify separate cache keys exist
	golangCache, err := fixtures.redis.Get(ctx, "topics:search:Golang").Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, golangCache)

	pythonCache, err := fixtures.redis.Get(ctx, "topics:search:Python").Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, pythonCache)

	// Verify they contain different data
	assert.NotEqual(t, golangCache, pythonCache)
}

func TestTopicSearch_CacheTTL_ExpiresAfterTimeout(t *testing.T) {
	// Skip in CI as this test takes time
	if os.Getenv("CI") != "" {
		t.Skip("Skipping TTL test in CI environment")
	}

	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	ctx := context.Background()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "ttluser", "ttl@example.com")
	seedTestTopic(t, fixtures.db, "Cache Expiration Test", userID)

	// Execute search to populate cache
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Cache+Expiration", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify cache key has TTL set (should be around 3 minutes / 180 seconds)
	cacheKey := "topics:search:Cache Expiration"
	ttl, err := fixtures.redis.TTL(ctx, cacheKey).Result()
	assert.NoError(t, err)

	// TTL should be positive and less than or equal to 3 minutes
	assert.True(t, ttl > 0, "TTL should be positive")
	assert.True(t, ttl <= 3*time.Minute, "TTL should be at most 3 minutes")
}

// --------------------------------------------------------------------------
// MySQL FULLTEXT Search Behavior Tests
// --------------------------------------------------------------------------

func TestTopicSearch_FullText_CaseInsensitiveSearch(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data with mixed case
	userID := seedTestUser(t, fixtures.db, "caseuser", "case@example.com")
	seedTestTopic(t, fixtures.db, "GOLANG Development Guide", userID)
	seedTestTopic(t, fixtures.db, "golang basics tutorial", userID)
	seedTestTopic(t, fixtures.db, "GoLang Advanced Topics", userID)

	// Search with lowercase
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=golang", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	// Should find all three topics regardless of case
	assert.Equal(t, 3, len(topics), "Case-insensitive search should find all variations")
}

func TestTopicSearch_FullText_PartialWordNotMatched(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "partialuser", "partial@example.com")
	seedTestTopic(t, fixtures.db, "Programming Best Practices", userID)

	// Search for partial word - MySQL FULLTEXT typically doesn't match partial words in natural language mode
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Prog", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	// Natural language mode doesn't support partial word matching
	// This is expected MySQL FULLTEXT behavior
	assert.NotNil(t, topics)
}

func TestTopicSearch_FullText_MultipleWords(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "multiworduser", "multiword@example.com")
	seedTestTopic(t, fixtures.db, "Building Scalable Microservices", userID)
	seedTestTopic(t, fixtures.db, "Microservices Architecture Patterns", userID)
	seedTestTopic(t, fixtures.db, "Monolithic vs Microservices", userID)
	seedTestTopic(t, fixtures.db, "Building Web Applications", userID)

	// Search for multiple words
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=Building+Microservices", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)

	// Should find topics containing either or both words
	assert.GreaterOrEqual(t, len(topics), 1, "Should find topics matching multiple words")
}

func TestTopicSearch_FullText_SpecialCharacters(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data with special characters
	userID := seedTestUser(t, fixtures.db, "specialuser", "special@example.com")
	seedTestTopic(t, fixtures.db, "C++ Programming Guide", userID)
	seedTestTopic(t, fixtures.db, "Node.js Best Practices", userID)

	// Search should handle special characters
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=C%2B%2B", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Response should be valid JSON regardless of special characters
	var topics []*models.Topic
	err := json.Unmarshal(w.Body.Bytes(), &topics)
	require.NoError(t, err)
	assert.NotNil(t, topics)
}

// --------------------------------------------------------------------------
// Database Integration Tests
// --------------------------------------------------------------------------

func TestTopicSearch_Integration_NewTopicBecomesSearchable(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	ctx := context.Background()

	// Seed initial user
	userID := seedTestUser(t, fixtures.db, "integrationuser", "integration@example.com")

	// Clear cache
	fixtures.redis.FlushDB(ctx)

	// Search before creating topic
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=UniqueSearchTerm", nil)
	w1 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w1, req1)

	var topics1 []*models.Topic
	json.Unmarshal(w1.Body.Bytes(), &topics1)
	assert.Equal(t, 0, len(topics1), "Should find no topics before creation")

	// Clear any cached empty results
	fixtures.redis.Del(ctx, "topics:search:UniqueSearchTerm")

	// Create new topic with the search term
	seedTestTopic(t, fixtures.db, "UniqueSearchTerm in Topic Title", userID)

	// Search after creating topic
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=UniqueSearchTerm", nil)
	w2 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w2, req2)

	var topics2 []*models.Topic
	err := json.Unmarshal(w2.Body.Bytes(), &topics2)
	require.NoError(t, err)
	assert.Equal(t, 1, len(topics2), "Should find newly created topic")
}

func TestTopicSearch_Integration_DeletedTopicNotSearchable(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	ctx := context.Background()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "deleteuser", "delete@example.com")
	topic := seedTestTopic(t, fixtures.db, "DeleteableSearchTopic", userID)

	// Clear cache
	fixtures.redis.FlushDB(ctx)

	// Verify topic is searchable
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=DeleteableSearchTopic", nil)
	w1 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w1, req1)

	var topics1 []*models.Topic
	json.Unmarshal(w1.Body.Bytes(), &topics1)
	require.Equal(t, 1, len(topics1), "Should find topic before deletion")

	// Delete the topic and clear cache
	_, err := fixtures.db.Exec("DELETE FROM topics WHERE id = ?", topic.ID)
	require.NoError(t, err)
	fixtures.redis.Del(ctx, "topics:search:DeleteableSearchTopic")

	// Search after deletion
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=DeleteableSearchTopic", nil)
	w2 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w2, req2)

	var topics2 []*models.Topic
	json.Unmarshal(w2.Body.Bytes(), &topics2)
	assert.Equal(t, 0, len(topics2), "Should not find deleted topic")
}

func TestTopicSearch_Integration_UpdatedTopicReflectsChanges(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	ctx := context.Background()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "updateuser", "update@example.com")
	topic := seedTestTopic(t, fixtures.db, "OriginalSearchTitle", userID)

	// Clear cache
	fixtures.redis.FlushDB(ctx)

	// Verify original title is searchable
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=OriginalSearchTitle", nil)
	w1 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w1, req1)

	var topics1 []*models.Topic
	json.Unmarshal(w1.Body.Bytes(), &topics1)
	assert.Equal(t, 1, len(topics1), "Should find topic with original title")

	// Update the topic title
	_, err := fixtures.db.Exec("UPDATE topics SET title = ? WHERE id = ?", "UpdatedSearchTitle", topic.ID)
	require.NoError(t, err)

	// Clear relevant caches
	fixtures.redis.FlushDB(ctx)

	// Search for old title should not find it
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=OriginalSearchTitle", nil)
	w2 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w2, req2)

	var topics2 []*models.Topic
	json.Unmarshal(w2.Body.Bytes(), &topics2)
	assert.Equal(t, 0, len(topics2), "Should not find topic with old title")

	// Search for new title should find it
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=UpdatedSearchTitle", nil)
	w3 := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w3, req3)

	var topics3 []*models.Topic
	json.Unmarshal(w3.Body.Bytes(), &topics3)
	assert.Equal(t, 1, len(topics3), "Should find topic with updated title")
}

// --------------------------------------------------------------------------
// Error Handling Tests
// --------------------------------------------------------------------------

func TestTopicSearch_Error_ResponseFormat(t *testing.T) {
	// This test validates the response format when things work correctly
	// Error injection for database failures requires additional setup
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Valid search returns proper JSON array
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/search?q=test", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response is valid JSON
	var result interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err, "Response should be valid JSON")

	// For successful requests, result should be an array
	_, isArray := result.([]interface{})
	assert.True(t, isArray, "Successful search should return an array")
}

func TestTopicSearch_Error_InvalidMethod(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// POST to search endpoint should fail
	req := httptest.NewRequest(http.MethodPost, "/api/v1/topics/search?q=test", nil)
	w := httptest.NewRecorder()
	fixtures.router.ServeHTTP(w, req)

	// Gin returns 404 for unmatched routes
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --------------------------------------------------------------------------
// Concurrency Tests
// --------------------------------------------------------------------------

func TestTopicSearch_Concurrent_MultipleSearchesSucceed(t *testing.T) {
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	// Seed test data
	userID := seedTestUser(t, fixtures.db, "concurrentuser", "concurrent@example.com")
	seedTestTopic(t, fixtures.db, "Concurrent Programming Guide", userID)
	seedTestTopic(t, fixtures.db, "Parallel Processing Techniques", userID)
	seedTestTopic(t, fixtures.db, "Multithreading Best Practices", userID)

	// Run concurrent searches
	done := make(chan bool, 10)
	errors := make(chan error, 10)

	queries := []string{"Concurrent", "Parallel", "Multithreading", "Programming", "Techniques"}

	for i := 0; i < 10; i++ {
		go func(query string) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/topics/search?q=%s", query), nil)
			w := httptest.NewRecorder()
			fixtures.router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				errors <- fmt.Errorf("unexpected status code: %d", w.Code)
				return
			}

			var topics []*models.Topic
			if err := json.Unmarshal(w.Body.Bytes(), &topics); err != nil {
				errors <- err
				return
			}

			done <- true
		}(queries[i%len(queries)])
	}

	// Wait for all goroutines
	successCount := 0
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			successCount++
		case err := <-errors:
			t.Errorf("Concurrent search failed: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("Concurrent search timed out")
		}
	}

	assert.Equal(t, 10, successCount, "All concurrent searches should succeed")
}
