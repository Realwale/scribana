package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/blog-api/pkg/storage"
)

type UploadHandler struct {
	storage *storage.Storage
}

func NewUploadHandler(storage *storage.Storage) *UploadHandler {
	return &UploadHandler{storage: storage}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Single file upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file size (e.g., 5MB limit)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	filename, err := h.storage.SaveImage(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	imageURL := filepath.Join("/uploads", filename)
	c.JSON(http.StatusOK, gin.H{
		"url": imageURL,
	})
}
