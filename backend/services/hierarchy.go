package services

import (
	"os"
	"path/filepath"
)

type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children,omitempty"`
}

// GetHierarchy 遞迴取得目錄樹
func GetHierarchy(root string) ([]FileNode, error) {
	var nodes []FileNode
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		node := FileNode{
			Name:  entry.Name(),
			Path:  filepath.Join(root, entry.Name()),
			IsDir: entry.IsDir(),
		}
		if entry.IsDir() {
			children, err := GetHierarchy(filepath.Join(root, entry.Name()))
			if err == nil {
				node.Children = children
			}
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
