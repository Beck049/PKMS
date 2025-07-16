package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ReadMarkdownFile reads a markdown file from the given path
func ReadMarkdownFile(basePath, filePath string) (string, error) {
	// Clean and join paths
	fullPath := filepath.Join(basePath, filePath)
	fullPath = filepath.Clean(fullPath)

	// Security check: ensure the path is within basePath
	relPath, err := filepath.Rel(basePath, fullPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return "", os.ErrPermission
	}

	// Read file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ValidateMarkdownPath checks if the path is valid and points to a markdown file
func ValidateMarkdownPath(path string) bool {
	// Check if path ends with .md
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return false
	}

	// Check if path contains any potentially dangerous patterns
	dangerousPatterns := []string{"..", "//", "\\"}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(path, pattern) {
			return false
		}
	}

	return true
}
