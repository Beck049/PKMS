package api

import (
	"database/sql"
	"net/http"

	"pkms/backend/services"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	DB *sql.DB
}

func NewTagHandler(db *sql.DB) *TagHandler {
	return &TagHandler{DB: db}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	query := c.Query("query")
	tags, err := services.GetTags(h.DB, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}
