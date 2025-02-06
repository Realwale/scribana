package main

import (
	"github.com/Realwale/scribana/internal/handlers"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/Realwale/scribana/docs"
	"github.com/Realwale/scribana/internal/middleware"
	"github.com/Realwale/scribana/internal/models"
	"github.com/Realwale/scribana/internal/services"
	"github.com/Realwale/scribana/pkg/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Blog API
// @version 1.0
// @description A RESTful API for a blog platform
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Database connection
	dsn := "host=localhost user=postgres password=password dbname=blog_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Category{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize services
	authService := services.NewAuthService(db)

	// Setup upload directory
	uploadDir := filepath.Join("uploads")
	storageService := storage.NewStorage(uploadDir)
	uploadHandler := handlers.NewUploadHandler(storageService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(db)
	commentHandler := handlers.NewCommentHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)

	// Initialize Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Routes
	api := r.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public post routes
		api.GET("/posts", postHandler.GetPosts)
		api.GET("/posts/:slug", postHandler.GetPost)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Posts
			posts := protected.Group("/posts")
			{
				posts.POST("/", middleware.RoleMiddleware(models.AuthorRole), postHandler.CreatePost)
				posts.PUT("/:id", middleware.RoleMiddleware(models.AuthorRole), postHandler.UpdatePost)
				posts.DELETE("/:id", middleware.RoleMiddleware(models.AuthorRole), postHandler.DeletePost)
			}

			// Comments
			comments := protected.Group("/comments")
			{
				comments.POST("/", commentHandler.CreateComment)
				comments.PUT("/:id", commentHandler.UpdateComment)
				comments.DELETE("/:id", commentHandler.DeleteComment)
			}

			// Categories (Admin only)
			categories := protected.Group("/categories")
			categories.Use(middleware.RoleMiddleware(models.AdminRole))
			{
				categories.POST("/", categoryHandler.CreateCategory)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// Upload routes (restricted to authors and admins)
			uploads := protected.Group("/uploads")
			uploads.Use(middleware.RoleMiddleware(models.AuthorRole))
			{
				uploads.POST("/image", uploadHandler.UploadImage)
			}
		}
	}

	// Add Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files
	r.Static("/uploads", uploadDir)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
