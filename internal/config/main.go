package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultFileName = ".version.json"
)

func NormalizePath(dirPath string) (string, error) {
	if len(dirPath) > 0 {
		if isRelativeDirPath(dirPath) {
			return toAbsoluteDirPath(dirPath)
		} else {
			return dirPath, nil
		}
	}
	return getCurrentDirPath()
}

func getCurrentDirPath() (string, error) {
	dirPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current path: %v", err)
	}
	return dirPath, nil
}

func toAbsoluteDirPath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", fmt.Errorf("error converting to absolute path: %v", err)
	}
	return absPath, nil
}

func isRelativeDirPath(path string) bool {
	return !filepath.IsAbs(path)
}

func FindConfigVersionFilePath(path string) ([]string, error) {
	var list []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == defaultFileName {
			list = append(list, path)
		}
		return nil
	})

	return list, err
}
