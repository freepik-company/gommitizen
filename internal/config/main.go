package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultFileName = ".version.json"
)

func NormalizePath(path string) (string, error) {
	if len(path) > 0 {
		if isRelativePath(path) {
			return toAbsolutePath(path)
		} else {
			return path, nil
		}
	}
	return getCurrentPath()
}

func getCurrentPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current path: %v", err)
	}
	return path, nil
}

func toAbsolutePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", fmt.Errorf("error converting to absolute path: %v", err)
	}
	return absPath, nil
}

func isRelativePath(path string) bool {
	return !filepath.IsAbs(path)
}
