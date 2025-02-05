package handlers

import (
	"net/http"

	"github.com/Realwale/scribana/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// @Summary Create new category
// @Description Create a new blog category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param category body CreateCategoryRequest true "Category details"
// @Success 201 {object} models.Category
// @Failure 400 {object} ErrorResponse
// @Failure 401,403 {object} ErrorResponse
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name: req.Name,
		Slug: slug.Make(req.Name),
	}

	if err := h.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// @Summary Update category
// @Description Update an existing category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID"
// @Param category body CreateCategoryRequest true "Category details"
// @Success 200 {object} models.Category
// @Failure 400,401,403,404 {object} ErrorResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := h.db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = req.Name
	category.Slug = slug.Make(req.Name)

	if err := h.db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Delete category
// @Description Delete a category (Admin only)
// @Tags categories
// @Security Bearer
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400,401,403,404 {object} ErrorResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := h.db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if category has posts
	var count int64
	h.db.Model(&models.Post{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with existing posts"})
		return
	}

	if err := h.db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
