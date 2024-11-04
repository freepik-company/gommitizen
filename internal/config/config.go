package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	defaultFileName = ".version.json"
)

type PrintPlainOption int

const (
	PrintAll PrintPlainOption = iota
	PrintPathOnly
	PrintVersionOnly
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

func FindConfigVersionFilePathByPrefix(path, tag_prefix string) ([]string, error) {
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

			if dataPrefix, ok := data["tag_prefix"].(string); ok && dataPrefix == tag_prefix {
				list = append(list, subpath)
			}
		}
		return nil
	})

	return list, err
}

func PrintConfigVersions(configVersions []*ConfigVersion, format string) error {
	var err error

	switch format {
	case "json":
		err = printConfigVersionsJSON(configVersions)
	case "yaml":
		err = printConfigVersionsYAML(configVersions)
	}
	if err != nil {
		return fmt.Errorf("error printing config versions: %v", err)
	}
	return nil
}

func printConfigVersionsJSON(configVersions []*ConfigVersion) error {
	obj := make(map[string]interface{})
	obj["config_versions"] = configVersions

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	fmt.Println(string(data))
	return nil
}

func printConfigVersionsYAML(configVersions []*ConfigVersion) error {
	obj := make(map[string]interface{})
	obj["config_versions"] = configVersions

	data, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	fmt.Println(string(data))
	return nil
}

func PrintConfigVersionsPlain(configVersions []*ConfigVersion, option PrintPlainOption) error {
	for _, configVersion := range configVersions {
		switch option {
		case PrintPathOnly:
			fmt.Printf("%s: %s\n", configVersion.TagPrefix, configVersion.GetDirPath())
		case PrintVersionOnly:
			fmt.Printf("%s: %s\n", configVersion.TagPrefix, configVersion.Version)
		case PrintAll:
			fmt.Println("Tag Prefix:", configVersion.TagPrefix)
			fmt.Println("Directory:", configVersion.GetDirPath())
			fmt.Println("Version:", configVersion.Version)
			fmt.Println("Commit:", configVersion.Commit)
			fmt.Println("Version Files:")
			for _, file := range configVersion.VersionFiles {
				fmt.Println(" -", file)
			}
			fmt.Println()
		}
	}
	return nil
}
