package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"pkms/backend/api"
	"pkms/backend/config"
	"pkms/backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	// Initialize services
	contentService := services.NewContentService(cfg)
	articleService := services.NewArticleService(db)

	// Initialize handlers
	contentHandler := api.NewContentHandler(contentService, articleService)
	hierarchyHandler := api.NewHierarchyHandler(cfg)
	tagHandler := api.NewTagHandler(db)
	searchHandler := api.NewSearchHandler(db, contentService)

	// 新增 ArticleHandler
	articleHandler := api.NewArticleHandler(articleService, cfg)

	// Setup router
	r := gin.Default()

	// set CROS settings
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))
	// API routes
	apiGroup := r.Group("/api")
	{
		// Content routes
		apiGroup.GET("/content/:id", contentHandler.GetArticleContent)
		// Hierarchy route
		apiGroup.GET("/hierarchy", hierarchyHandler.GetHierarchy)

		// Article routes
		apiGroup.POST("/articles", articleHandler.CreateArticle)
		apiGroup.DELETE("/articles/:id", articleHandler.DeleteArticle)
		apiGroup.PUT("/articles/:id", articleHandler.UpdateArticle)

		// Tag routes
		apiGroup.GET("/tags", tagHandler.GetTags)

		// Search routes
		apiGroup.GET("/search", searchHandler.SearchArticles)
	}

	// Start server
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
