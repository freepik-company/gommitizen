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

func FindConfigVersionFilePathByAlias(path, alias string) ([]string, error) {
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

			if dataAlias, ok := data["alias"].(string); ok && dataAlias == alias {
				list = append(list, subpath)
			}
		}
		return nil
	})

	return list, err
}
