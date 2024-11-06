package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
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

	err := filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == defaultFileName {
			list = append(list, subpath)
		}
		return nil
	})

	return list, err
}

func FindConfigVersionFilePathByPrefix(path, tagPrefix string) ([]string, error) {
	var list []string

	err := filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking subpath: %v", err)
		}
		if !info.IsDir() && info.Name() == defaultFileName {
			file, err := os.Open(subpath)
			if err != nil {
				return fmt.Errorf("error opening file: %v", err)
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					slog.Error(fmt.Sprintf("error closing file: %v", err))
				}
			}(file)

			var data map[string]interface{}
			if err := json.NewDecoder(file).Decode(&data); err != nil {
				return fmt.Errorf("error decoding file: %v", err)
			}

			if dataPrefix, ok := data["tag_prefix"].(string); ok && dataPrefix == tagPrefix {
				list = append(list, subpath)
			}
		}
		return nil
	})

	return list, err
}
