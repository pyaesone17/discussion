package api

import (
	"discussion-forum/internal/models"
	"discussion-forum/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	service *service.VoteService
}

func NewVoteHandler(service *service.VoteService) *VoteHandler {
	return &VoteHandler{service: service}
}

func (h *VoteHandler) CastVote(votableType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var req models.VoteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Value != 1 && req.Value != -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "value must be 1 or -1"})
			return
		}

		if err := h.service.CastVote(req.UserID, id, votableType, req.Value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "vote recorded"})
	}
}

func (h *VoteHandler) RemoveVote(votableType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var req models.VoteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.service.RemoveVote(req.UserID, id, votableType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "vote removed"})
	}
}
