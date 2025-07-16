package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"pkms/backend/config"
	"pkms/backend/services"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	Service *services.ArticleService
	Cfg     *config.Config
}

func NewArticleHandler(service *services.ArticleService, cfg *config.Config) *ArticleHandler {
	return &ArticleHandler{Service: service, Cfg: cfg}
}

// CreateArticle godoc
// @Summary Create a new article
// @Accept json
// @Produce json
// @Param article body services.CreateArticleInput true "Article info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req services.CreateArticleInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	result, err := h.Service.CreateArticle(req, h.Cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Article created successfully",
		"article_id": result.ArticleID,
		"path":       result.Path,
	})
}

// DeleteArticle godoc
// @Summary Delete an article by ID
// @Param id path int true "Article ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid article ID"})
		return
	}

	err = h.Service.DeleteArticle(id, h.Cfg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Article not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Article deleted successfully"})
}

// UpdateArticle godoc
// @Summary Update an article by ID
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param article body services.UpdateArticleInput true "Updated article info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid article ID"})
		return
	}

	var req services.UpdateArticleInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	err = h.Service.UpdateArticle(id, req, h.Cfg)
	if err != nil {
		if err == services.ErrArticleNotFound {
			c.JSON(404, gin.H{"error": "Article not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Article updated successfully"})
}
