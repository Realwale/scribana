package handlers

import (
	"net/http"
	"strconv"

	"github.com/Realwale/scribana/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

// @Summary Create new comment
// @Description Create a new comment on a blog post
// @Tags comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param comment body CreateCommentRequest true "Comment details"
// @Success 201 {object} models.Comment
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := strconv.ParseUint(c.GetString("userID"), 10, 64)
	comment := models.Comment{
		Content: req.Content,
		PostID:  req.PostID,
		UserID:  uint(userID),
	}

	if err := h.db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary Update comment
// @Description Update an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Comment ID"
// @Param comment body CreateCommentRequest true "Comment details"
// @Success 200 {object} models.Comment
// @Failure 400,401,403,404 {object} ErrorResponse
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment
	if err := h.db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	userID, _ := strconv.ParseUint(c.GetString("userID"), 10, 64)
	if comment.UserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this comment"})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.Content = req.Content
	if err := h.db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Delete comment
// @Description Delete a comment
// @Tags comments
// @Security Bearer
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 401,403,404 {object} ErrorResponse
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment
	if err := h.db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	userID, _ := strconv.ParseUint(c.GetString("userID"), 10, 64)
	if comment.UserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this comment"})
		return
	}

	if err := h.db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
