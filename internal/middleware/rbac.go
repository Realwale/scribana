package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/blog-api/internal/models"
	"gorm.io/gorm"
)

func RoleMiddleware(requiredRole models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		userID := c.GetString("userID")

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if user.Role != requiredRole && user.Role != models.AdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
