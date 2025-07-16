package services

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pkms/backend/config"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

type Article struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Path       string    `json:"path"`
	Type       string    `json:"type"`
	CreateDate time.Time `json:"create_date"`
	EditDate   time.Time `json:"edit_date"`
	RefCount   int       `json:"ref_count"`
	Pin        bool      `json:"pin"`
}

type ArticleService struct {
	db *sql.DB
}

func NewArticleService(db *sql.DB) *ArticleService {
	return &ArticleService{db: db}
}

// GetArticleByID retrieves an article by its ID
func (s *ArticleService) GetArticleByID(id uint) (*Article, error) {
	query := `
		SELECT id, title, path, type, create_date, edit_date, ref_count, pin
		FROM articles 
		WHERE id = ?
	`

	var article Article
	err := s.db.QueryRow(query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Path,
		&article.Type,
		&article.CreateDate,
		&article.EditDate,
		&article.RefCount,
		&article.Pin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	return &article, nil
}

type CreateArticleInput struct {
	Title string
	Path  string
	Type  string
	Desc  string
	Tags  []string
}

type UpdateArticleInput struct {
	Title   *string  `json:"title,omitempty"`
	Path    *string  `json:"path,omitempty"`
	Type    *string  `json:"type,omitempty"`
	Pin     *bool    `json:"pin,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Content *string  `json:"content,omitempty"`
}

type CreateArticleResult struct {
	ArticleID int64
	Path      string
}

func (s *ArticleService) CreateArticle(input CreateArticleInput, cfg *config.Config) (*CreateArticleResult, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert to article table
	now := time.Now()
	result, err := tx.Exec(`
		INSERT INTO articles (title, path, type, create_date, edit_date, ref_count, pin)
		VALUES (?, ?, ?, ?, ?, 0, false)
	`, input.Title, input.Path, input.Type, now, now)
	if err != nil {
		return nil, err
	}
	articleID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	// Insert to article_tags
	for _, tagName := range input.Tags {
		var tagID int64
		// add tag if not exist
		err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			result, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
			if err != nil {
				return nil, err
			}
			tagID, err = result.LastInsertId()
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		_, err = tx.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articleID, tagID)
		if err != nil {
			return nil, err
		}
	}

	// Create article file with YAML frontmatter
	targetPath := filepath.Join(cfg.SearchPath, input.Path)
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, err
	}

	// Create YAML frontmatter
	tagsStr := ""
	if len(input.Tags) > 0 {
		tagsStr = "\"" + input.Tags[0] + "\""
		for i := 1; i < len(input.Tags); i++ {
			tagsStr += ", \"" + input.Tags[i] + "\""
		}
	}

	content := fmt.Sprintf("---\ntitle: '%s'\ntags: [%s]\ntype: '%s'\n---\n\n# %s\n\n%s",
		input.Title, tagsStr, input.Type, input.Title, input.Desc)

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer targetFile.Close()

	_, err = targetFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CreateArticleResult{
		ArticleID: articleID,
		Path:      input.Path,
	}, nil
}

func (s *ArticleService) DeleteArticle(id int64, cfg *config.Config) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 取得 path
	var path string
	err = tx.QueryRow("SELECT path FROM articles WHERE id = ?", id).Scan(&path)
	if err != nil {
		return err
	}

	// 2. 刪除檔案
	filePath := filepath.Join(cfg.SearchPath, path)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 3. 刪除 article_tags
	_, err = tx.Exec("DELETE FROM article_tags WHERE article_id = ?", id)
	if err != nil {
		return err
	}

	// 4. 刪除 articles
	_, err = tx.Exec("DELETE FROM articles WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ArticleService) UpdateArticle(id int64, input UpdateArticleInput, cfg *config.Config) error {
	// 先取得現有資料
	currentArticle, err := s.GetArticleByID(uint(id))
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 準備最新值
	title := currentArticle.Title
	if input.Title != nil {
		title = *input.Title
	}
	path := currentArticle.Path
	if input.Path != nil {
		path = *input.Path
	}
	typeValue := currentArticle.Type
	if input.Type != nil {
		typeValue = *input.Type
	}
	pin := currentArticle.Pin
	if input.Pin != nil {
		pin = *input.Pin
	}

	// 1. 更新 articles table
	now := time.Now()
	_, err = tx.Exec(`
		UPDATE articles 
		SET title = ?, path = ?, type = ?, pin = ?, edit_date = ?
		WHERE id = ?
	`, title, path, typeValue, pin, now, id)
	if err != nil {
		return err
	}

	// 2. tags 有提供才更新
	var tags []string
	if input.Tags != nil {
		// 刪除所有 article_tags
		_, err = tx.Exec("DELETE FROM article_tags WHERE article_id = ?", id)
		if err != nil {
			return err
		}
		// 重新插入
		for _, tagName := range input.Tags {
			var tagID int64
			err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
			if err == sql.ErrNoRows {
				result, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
				if err != nil {
					return err
				}
				tagID, err = result.LastInsertId()
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			_, err = tx.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", id, tagID)
			if err != nil {
				return err
			}
		}
		tags = input.Tags
	} else {
		// 沒有提供 tags，查詢現有 tags
		rows, err := tx.Query(`
			SELECT t.name 
			FROM tags t 
			JOIN article_tags at ON t.id = at.tag_id 
			WHERE at.article_id = ?
		`, id)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tagName string
			if err := rows.Scan(&tagName); err != nil {
				return err
			}
			tags = append(tags, tagName)
		}
	}

	// 3. 只要有 title、type、tags、content、path 任一有提供就重寫檔案
	needUpdateFile := input.Title != nil || input.Type != nil || input.Tags != nil || input.Content != nil || input.Path != nil
	if needUpdateFile {
		targetPath := filepath.Join(cfg.SearchPath, path)
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}

		// tags YAML 字串
		tagsStr := ""
		if len(tags) > 0 {
			tagsStr = "\"" + tags[0] + "\""
			for i := 1; i < len(tags); i++ {
				tagsStr += ", \"" + tags[i] + "\""
			}
		}

		// content
		content := ""
		if input.Content != nil {
			content = *input.Content
		} else {
			// 讀現有檔案內容，去除 frontmatter
			currentFilePath := filepath.Join(cfg.SearchPath, currentArticle.Path)
			if fileContent, err := os.ReadFile(currentFilePath); err == nil {
				str := string(fileContent)
				if idx := strings.Index(str, "\n---\n"); idx != -1 {
					content = str[idx+5:]
				} else {
					content = str
				}
			}
		}

		fileContent := fmt.Sprintf("---\ntitle: '%s'\ntags: [%s]\ntype: '%s'\n---\n\n%s",
			title, tagsStr, typeValue, content)

		targetFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = targetFile.WriteString(fileContent)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
