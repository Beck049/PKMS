package api

import (
	"net/http"
	"os"

	"pkms/backend/config"
	"pkms/backend/services"

	"github.com/gin-gonic/gin"
)

type HierarchyHandler struct {
	cfg *config.Config
}

func NewHierarchyHandler(cfg *config.Config) *HierarchyHandler {
	return &HierarchyHandler{cfg: cfg}
}

// GetHierarchy 回傳 articles 目錄樹
func (h *HierarchyHandler) GetHierarchy(c *gin.Context) {
	root := h.cfg.SearchPath
	if _, err := os.Stat(root); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "articles directory not found"})
		return
	}
	nodes, err := services.GetHierarchy(root)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}
