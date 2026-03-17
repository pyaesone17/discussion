package api

import (
	"discussion-forum/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	userService *service.UserService,
	topicService *service.TopicService,
	postService *service.PostService,
	voteService *service.VoteService,
	reportService *service.ReportService,
) {
	api := router.Group("/api/v1")

	userHandler := NewUserHandler(userService)
	api.POST("/users", userHandler.Create)
	api.GET("/users/:id", userHandler.GetByID)
	api.GET("/users", userHandler.GetAll)
	api.PUT("/users/:id", userHandler.Update)
	api.DELETE("/users/:id", userHandler.Delete)

	topicHandler := NewTopicHandler(topicService)
	api.POST("/topics", topicHandler.Create)
	api.GET("/topics", topicHandler.GetAll)
	api.GET("/topics/search", topicHandler.Search)
	api.GET("/topics/:id", topicHandler.GetByID)
	api.PUT("/topics/:id", topicHandler.Update)
	api.DELETE("/topics/:id", topicHandler.Delete)

	postHandler := NewPostHandler(postService)
	api.POST("/posts", postHandler.Create)
	api.GET("/posts/:id", postHandler.GetByID)
	api.GET("/topics/:id/posts", postHandler.GetByTopicID)
	api.PUT("/posts/:id", postHandler.Update)
	api.DELETE("/posts/:id", postHandler.Delete)

	voteHandler := NewVoteHandler(voteService)
	api.POST("/topics/:id/vote", voteHandler.CastVote("topic"))
	api.DELETE("/topics/:id/vote", voteHandler.RemoveVote("topic"))
	api.POST("/posts/:id/vote", voteHandler.CastVote("post"))
	api.DELETE("/posts/:id/vote", voteHandler.RemoveVote("post"))

	reportHandler := NewReportHandler(reportService)
	api.POST("/topics/:id/report", reportHandler.CreateReport("topic"))
	api.POST("/posts/:id/report", reportHandler.CreateReport("post"))
	api.GET("/reports", reportHandler.GetAll)
	api.PUT("/reports/:id", reportHandler.UpdateStatus)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
