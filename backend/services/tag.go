package services

import (
	"database/sql"
	"strings"
)

type Tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetTags(db *sql.DB, query string) ([]Tag, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if query != "" {
		like := "%" + strings.ToLower(query) + "%"
		rows, err = db.Query("SELECT id, name FROM tags WHERE LOWER(name) LIKE ?", like)
	} else {
		rows, err = db.Query("SELECT id, name FROM tags")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
