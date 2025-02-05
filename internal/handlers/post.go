package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/yourusername/blog-api/internal/models"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{db: db}
}

type CreatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
	ImageURL   string `json:"image_url"`
}

// @Summary Create new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Security Bearer
// @Param post body CreatePostRequest true "Post details"
// @Success 201 {object} models.Post
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := strconv.ParseUint(c.GetString("userID"), 10, 64)
	post := models.Post{
		Title:      req.Title,
		Slug:       slug.Make(req.Title),
		Content:    req.Content,
		AuthorID:   uint(userID),
		CategoryID: req.CategoryID,
		ImageURL:   req.ImageURL,
	}

	if err := h.db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// @Summary Get all posts
// @Description Get all blog posts
// @Tags posts
// @Produce json
// @Param category query string false "Filter by category slug"
// @Success 200 {array} models.Post
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	var posts []models.Post
	query := h.db.Preload("Author").Preload("Category").Preload("Comments")

	if category := c.Query("category"); category != "" {
		query = query.Joins("JOIN categories ON categories.id = posts.category_id").
			Where("categories.slug = ?", category)
	}

	if err := query.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// @Summary Get post by slug
// @Description Get a blog post by its slug
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} models.Post
// @Failure 404 {object} ErrorResponse
// @Router /posts/{slug} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")
	var post models.Post

	if err := h.db.Preload("Author").Preload("Category").Preload("Comments").
		Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Update post
// @Description Update an existing blog post
// @Tags posts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Post ID"
// @Param post body CreatePostRequest true "Post details"
// @Success 200 {object} models.Post
// @Failure 400,401,403,404 {object} ErrorResponse
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user is author or admin
	userID, _ := strconv.ParseUint(c.GetString("userID"), 10, 64)
	if post.AuthorID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this post"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = req.Title
	post.Slug = slug.Make(req.Title)
	post.Content = req.Content
	post.CategoryID = req.CategoryID
	post.ImageURL = req.ImageURL

	if err := h.db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Delete post
// @Description Delete a blog post
// @Tags posts
// @Security Bearer
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 401,403,404 {object} ErrorResponse
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user is author or admin
	userID, _ := strconv.ParseUint(c.GetString("userID"), 10, 64)
	if post.AuthorID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this post"})
		return
	}

	if err := h.db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
