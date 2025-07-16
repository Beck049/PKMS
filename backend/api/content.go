package api

import (
	"log"
	"net/http"
	"strconv"

	"pkms/backend/services"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	contentService *services.ContentService
	articleService *services.ArticleService
}

func NewContentHandler(contentService *services.ContentService, articleService *services.ArticleService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
		articleService: articleService,
	}
}

// GetArticleContent godoc
// @Summary Get article content by ID
// @Description Get article metadata and content by article ID
// @Tags content
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/content/{id} [get]
func (h *ContentHandler) GetArticleContent(c *gin.Context) {
	// Get contentID from URL parameter
	contentIDStr := c.Param("id")
	if contentIDStr == "" {
		log.Printf("content ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "content ID is required",
		})
		return
	}

	// Parse contentID to uint
	contentID, err := strconv.ParseUint(contentIDStr, 10, 32)
	if err != nil {
		log.Printf("invalid content ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid content ID format",
		})
		return
	}

	// Get article metadata from database
	article, err := h.articleService.GetArticleByID(uint(contentID))
	if err != nil {
		switch err {
		case services.ErrArticleNotFound:
			log.Printf("article not found: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "article not found",
			})
		default:
			log.Printf("GetArticleByID error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	// Get file content using the article path
	rawData, err := h.contentService.GetContent(article.Path)
	if err != nil {
		switch err {
		case services.ErrInvalidPath:
			log.Printf("invalid file path: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid file path",
			})
		case services.ErrFileNotFound:
			log.Printf("file not found: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file not found",
			})
		default:
			log.Printf("GetContent error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	// Return combined data
	c.JSON(http.StatusOK, gin.H{
		"id":        article.ID,
		"title":     article.Title,
		"path":      article.Path,
		"type":      article.Type,
		"ref_count": article.RefCount,
		"pin":       article.Pin,
		"rawdata":   rawData,
	})
}
