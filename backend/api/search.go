package api

import (
	"database/sql"
	"net/http"
	"strings"

	"pkms/backend/services"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	DB             *sql.DB
	ContentService *services.ContentService
}

func NewSearchHandler(db *sql.DB, contentService *services.ContentService) *SearchHandler {
	return &SearchHandler{DB: db, ContentService: contentService}
}

type ArticleResult struct {
	ID       uint     `json:"id"`
	Title    string   `json:"title"`
	Pin      bool     `json:"pin"`
	RefCount int      `json:"ref_count"`
	Sort     int      `json:"sort"`
	Tags     []string `json:"tags"`
}

type searchArticle struct {
	ID    uint
	Title string
	Path  string
}

func (h *SearchHandler) SearchArticles(c *gin.Context) {
	query := c.Query("query")
	path := c.Query("path")
	tagStr := c.Query("tag")

	// 1. 先用 path/tag filter 拿到 searchList
	/*
	 * SELECT a.id, a.title, a.path
	 * FROM articles a
	 * WHERE LOWER(a.path) LIKE LOWER('%food%')
	 * AND a.id IN (
	 *	  SELECT at.article_id
	 *	  FROM article_tags at
	 *	  JOIN tags t ON at.tag_id = t.id
	 *	  WHERE t.name IN ('sweet')
	 * )
	 */
	var args []interface{}
	var where []string

	if path != "" {
		where = append(where, "LOWER(a.path) LIKE LOWER(?)")
		args = append(args, "%"+path+"%")
	}
	if tagStr != "" {
		tags := strings.Split(tagStr, ",")
		placeholders := make([]string, len(tags))
		for i, t := range tags {
			placeholders[i] = "?"
			args = append(args, t)
		}
		where = append(where, "a.id IN (SELECT at.article_id FROM article_tags at JOIN tags t ON at.tag_id = t.id WHERE t.name IN ("+strings.Join(placeholders, ",")+"))")
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	querySQL := `SELECT a.id, a.title, a.path FROM articles a ` + whereClause
	rows, err := h.DB.Query(querySQL, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var searchList []searchArticle
	for rows.Next() {
		var a searchArticle
		if err := rows.Scan(&a.ID, &a.Title, &a.Path); err == nil {
			searchList = append(searchList, a)
		}
	}

	// 2. maintain resultList: [{id, sort}]
	type resultItem struct {
		ID   uint
		Sort int
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"searchList": searchList,
	// })
	var resultList []resultItem

	if query != "" {
		for i := 0; i < len(searchList); i++ {
			if strings.Contains(searchList[i].Title, query) {
				// 3. 先做 title string match Sort:1
				resultList = append(resultList, resultItem{searchList[i].ID, 1})
			} else {
				// 4. 剩下的做 content match Sort:2
				content, err := h.ContentService.GetContent(searchList[i].Path)
				if err == nil && strings.Contains(content, query) {
					resultList = append(resultList, resultItem{searchList[i].ID, 2})
				}
			}
		}
	} else {
		for i := 0; i < len(searchList); i++ {
			resultList = append(resultList, resultItem{ID: searchList[i].ID, Sort: 0})
		}
	}

	// 5. 根據 resultList id 查 DB
	if len(resultList) == 0 {
		c.JSON(http.StatusOK, []ArticleResult{})
		return
	}
	idSort := map[uint]int{}
	idList := make([]string, 0, len(resultList))
	for _, r := range resultList {
		idSort[r.ID] = r.Sort
		idList = append(idList, "?")
		args = append(args, r.ID)
	}

	// 查詢最終 Result
	finalSQL := `SELECT id, title, pin, ref_count FROM articles WHERE id IN (` + strings.Join(idList, ",") + `)`
	finalRows, err := h.DB.Query(finalSQL, args[len(args)-len(resultList):]...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer finalRows.Close()

	var results []ArticleResult
	allTagsSet := map[string]struct{}{}
	for finalRows.Next() {
		var a ArticleResult
		if err := finalRows.Scan(&a.ID, &a.Title, &a.Pin, &a.RefCount); err == nil {
			a.Sort = idSort[a.ID]
			// 查詢 tags
			tagRows, err := h.DB.Query(`SELECT t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = ?`, a.ID)
			if err == nil {
				var tagNames []string
				for tagRows.Next() {
					var tag string
					if err := tagRows.Scan(&tag); err == nil {
						tagNames = append(tagNames, tag)
						allTagsSet[tag] = struct{}{}
					}
				}
				tagRows.Close()
				a.Tags = tagNames
			}
			results = append(results, a)
		}
	}
	// collect allTags
	allTags := make([]string, 0, len(allTagsSet))
	for tag := range allTagsSet {
		allTags = append(allTags, tag)
	}
	c.JSON(http.StatusOK, gin.H{
		"articles": results,
		"allTags":  allTags,
	})
}
