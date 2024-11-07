package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v3"
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

func PrintConfigVersions(configVersions []*ConfigVersion, fields []string, format string) error {
	var err error

	switch format {
	case "json":
		err = printConfigVersionsJSON(configVersions, fields)
	case "yaml":
		err = printConfigVersionsYAML(configVersions, fields)
	case "plain":
		err = printConfigVersionsPlain(configVersions, fields)
	}
	if err != nil {
		return fmt.Errorf("error printing config versions: %v", err)
	}
	return nil
}

func printConfigVersionsJSON(configVersions []*ConfigVersion, fields []string) error {
	obj := make(map[string]interface{})
	obj["config_versions"] = configVersionFilter(configVersions, fields)

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	fmt.Println(string(data))
	return nil
}

func printConfigVersionsYAML(configVersions []*ConfigVersion, fields []string) error {
	obj := make(map[string]interface{})
	obj["config_versions"] = configVersionFilter(configVersions, fields)

	data, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	fmt.Println(string(data))
	return nil
}

func printConfigVersionsPlain(configVersions []*ConfigVersion, fields []string) error {
	obj := make(map[string]interface{})
	obj["config_versions"] = configVersionFilter(configVersions, fields)

	for _, configVersion := range obj["config_versions"].([]map[string]interface{}) {
		for fieldName, fieldValue := range configVersion {
			fmt.Printf("%s: %v\n", fieldName, fieldValue)
		}
		fmt.Println()
	}
	return nil
}

func configVersionFilter(configVersions []*ConfigVersion, fields []string) []map[string]interface{} {
	filteredConfigVersions := make([]map[string]interface{}, len(configVersions))

	for i, configVersion := range configVersions {
		val := reflect.ValueOf(configVersion).Elem()
		typ := val.Type()
		filteredConfigVersion := make(map[string]interface{})

		for j := 0; j < val.NumField(); j++ {
			var fieldName string
			if typ.Field(j).Name == "DirPath" {
				fieldName = "dir_path"
			} else {
				fieldName = typ.Field(j).Tag.Get("json")
			}
			if len(fields) == 0 {
				fieldValue := val.Field(j).Interface()
				filteredConfigVersion[fieldName] = fieldValue
				continue
			}
			for _, field := range fields {
				if fieldName == field {
					fieldValue := val.Field(j).Interface()
					filteredConfigVersion[fieldName] = fieldValue
					break
				}
			}
		}
		filteredConfigVersions[i] = filteredConfigVersion
	}
	return filteredConfigVersions
}
