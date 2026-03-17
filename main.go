package main

import (
	"log"
	"os"

	"discussion-forum/internal/api"
	"discussion-forum/internal/config"
	"discussion-forum/internal/database"
	"discussion-forum/internal/repository"
	"discussion-forum/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMySQLConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	redisClient := database.NewRedisClient(cfg.RedisURL)

	userRepo := repository.NewUserRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	postRepo := repository.NewPostRepository(db)
	voteRepo := repository.NewVoteRepository(db)
	reportRepo := repository.NewReportRepository(db)

	userService := service.NewUserService(userRepo, redisClient)
	topicService := service.NewTopicService(topicRepo, redisClient)
	postService := service.NewPostService(postRepo, redisClient)
	voteService := service.NewVoteService(voteRepo, postRepo, redisClient)
	reportService := service.NewReportService(reportRepo)

	router := gin.Default()

	api.SetupRoutes(router, userService, topicService, postService, voteService, reportService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
