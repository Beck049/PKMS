package services

import (
	"errors"
	"os"

	"pkms/backend/config"
	"pkms/backend/utils"
)

var (
	ErrInvalidPath  = errors.New("invalid markdown file path")
	ErrFileNotFound = errors.New("file not found")
)

type ContentService struct {
	basePath string
}

func NewContentService(cfg *config.Config) *ContentService {
	return &ContentService{
		basePath: cfg.SearchPath,
	}
}

// GetContent retrieves the content of a markdown file
func (s *ContentService) GetContent(path string) (string, error) {
	// Validate path
	if !utils.ValidateMarkdownPath(path) {
		return "", ErrInvalidPath
	}

	// Read file content
	content, err := utils.ReadMarkdownFile(s.basePath, path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrFileNotFound
		}
		return "", err
	}

	return content, nil
}
