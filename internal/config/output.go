package config

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

type ConfigVersionWrapper struct {
	DirPath       string                 `json:"dir_path" yaml:"dir_path"`
	ConfigVersion map[string]interface{} `json:"config_version" yaml:"config_version"`
}

type Wrapper struct {
	ConfigVersionWrappers []ConfigVersionWrapper `json:"config_versions" yaml:"config_versions"`
}

func PrintConfigVersions(configVersions []*ConfigVersion, fields []string, outputFormat string) (string, error) {
	wrapper := configVersionFilter(configVersions, fields, outputFormat)

	switch outputFormat {
	case "json":
		return printConfigVersionsJSON(wrapper)
	case "yaml":
		return printConfigVersionsYAML(wrapper)
	case "plain":
		return printConfigVersionsPlain(wrapper)
	}

	return "", fmt.Errorf("unsupported format: %s", outputFormat)
}

func printConfigVersionsJSON(wrapper Wrapper) (string, error) {
	data, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling data: %v", err)
	}

	return string(data), nil
}

func printConfigVersionsYAML(wrapper Wrapper) (string, error) {
	data, err := yaml.Marshal(wrapper)
	if err != nil {
		return "", fmt.Errorf("error marshalling data: %v", err)
	}

	return string(data), nil
}

func printConfigVersionsPlain(wrapper Wrapper) (string, error) {
	str := "config_versions:\n"
	for _, cvw := range wrapper.ConfigVersionWrappers {
		str += fmt.Sprintf("  dir_path: %s\n", cvw.DirPath)
		str += "  config_version:\n"
		for key, value := range cvw.ConfigVersion {
			str += fmt.Sprintf("    %s: %s\n", key, value)
		}
	}

	return str, nil
}

func configVersionFilter(configVersions []*ConfigVersion, fields []string, outputFormat string) Wrapper {
	wrapper := Wrapper{}

	for _, configVersion := range configVersions {
		cvw := ConfigVersionWrapper{
			DirPath:       configVersion.GetDirPath(),
			ConfigVersion: make(map[string]interface{}, 0),
		}

		val := reflect.ValueOf(configVersion).Elem()
		typ := val.Type()

		if len(fields) > 0 {
			for _, field := range fields {
				xField, ok := typ.FieldByName(field)
				if !ok {
					continue // Si el campo no existe, pasamos al siguiente
				}
				xValue := val.FieldByName(field)
				if xValue.IsValid() && xValue.CanInterface() {
					xTag := xField.Tag.Get(outputFormat)
					if xTag == "" {
						xTag = field
					}
					cvw.ConfigVersion[xTag] = xValue.Interface()
				}
			}
		} else {
			for i := 0; i < val.NumField(); i++ {
				xField := typ.Field(i)
				xValue := val.Field(i)
				if xValue.IsValid() && xValue.CanInterface() {
					xTag := xField.Tag.Get(outputFormat)
					if xTag == "" {
						xTag = xField.Name
					}
					cvw.ConfigVersion[xTag] = xValue.Interface()
				}
			}
		}

		wrapper.ConfigVersionWrappers = append(wrapper.ConfigVersionWrappers, cvw)
	}

	return wrapper
}
