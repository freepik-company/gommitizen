package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfigVersionWrapper struct {
	DirPath       string                 `json:"dir_path" yaml:"dir_path" plain:"dir_path"`
	FilePath      string                 `json:"file_path" yaml:"file_path" plain:"file_path"`
	LatestGitTag  string                 `json:"latest_git_tag" yaml:"latest_git_tag" plain:"latest_git_tag"`
	ConfigVersion map[string]interface{} `json:"config_version" yaml:"config_version" plain:"config_version"`
}

type Wrapper struct {
	ConfigVersionWrappers []ConfigVersionWrapper `json:"config_versions" yaml:"config_versions" plain:"config_versions"`
}

func PrintConfigVersions(configVersions []*ConfigVersion, fields []string, outputFormat string) (string, error) {
	tagFormat := outputFormat
	if tagFormat == "raw" {
		tagFormat = "plain"
	}
	wrapper := configVersionFilter(configVersions, fields, tagFormat)

	switch outputFormat {
	case "json":
		return printConfigVersionsJSON(wrapper)
	case "yaml":
		return printConfigVersionsYAML(wrapper)
	case "plain":
		return printConfigVersionsPlain(wrapper)
	case "raw":
		return printConfigVersionsRaw(wrapper)
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
	var sb strings.Builder
	sb.WriteString("config_versions:\n")
	for _, cvw := range wrapper.ConfigVersionWrappers {
		sb.WriteString(fmt.Sprintf("  dir_path: %s\n", cvw.DirPath))
		sb.WriteString(fmt.Sprintf("  file_path: %s\n", cvw.FilePath))
		sb.WriteString(fmt.Sprintf("  latest_git_tag: %s\n", cvw.LatestGitTag))
		sb.WriteString("  config_version:\n")
		for key, value := range cvw.ConfigVersion {
			typ := reflect.TypeOf(value)
			if typ.Kind() == reflect.Bool {
				sb.WriteString(fmt.Sprintf("    %s: %t\n", key, value))
			} else {
				sb.WriteString(fmt.Sprintf("    %s: %s\n", key, value))
			}
		}
	}

	return sb.String(), nil
}

func printConfigVersionsRaw(wrapper Wrapper) (string, error) {
	var lines []string
	single := len(wrapper.ConfigVersionWrappers) == 1

	for _, cvw := range wrapper.ConfigVersionWrappers {
		alias, _ := cvw.ConfigVersion["alias"]

		// Collect sorted non-alias keys for deterministic output
		var keys []string
		for key := range cvw.ConfigVersion {
			if key != "alias" {
				keys = append(keys, key)
			}
		}
		sort.Strings(keys)

		var line string
		if len(keys) == 0 {
			// Only alias field was requested (e.g., "get alias")
			line = fmt.Sprintf("%v", alias)
		} else if len(keys) == 1 {
			line = fmt.Sprintf("%v", cvw.ConfigVersion[keys[0]])
		} else {
			// Multiple fields: join as key:value pairs
			var pairs []string
			for _, key := range keys {
				pairs = append(pairs, fmt.Sprintf("%s:%v", key, cvw.ConfigVersion[key]))
			}
			line = strings.Join(pairs, ",")
		}

		if !single {
			line = fmt.Sprintf("%v=%s", alias, line)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

func configVersionFilter(configVersions []*ConfigVersion, fields []string, outputFormat string) Wrapper {
	wrapper := Wrapper{}

	for _, configVersion := range configVersions {
		cvw := ConfigVersionWrapper{
			DirPath:       configVersion.GetDirPath(),
			FilePath:      configVersion.GetFilePath(),
			LatestGitTag:  configVersion.GetGitTag(),
			ConfigVersion: make(map[string]interface{}),
		}

		val := reflect.ValueOf(configVersion).Elem()
		typ := val.Type()

		if len(fields) == 0 {
			fields = getAllFieldNames(typ)
		}

		for _, field := range fields {
			xField, ok := typ.FieldByName(field)
			if !ok {
				continue
			}
			xValue := val.FieldByName(field)
			if xValue.IsValid() && xValue.CanInterface() {
				xTag := strings.Split(xField.Tag.Get(outputFormat), ",")[0]
				if xTag == "" {
					xTag = field
				}
				cvw.ConfigVersion[xTag] = xValue.Interface()
			}
		}

		wrapper.ConfigVersionWrappers = append(wrapper.ConfigVersionWrappers, cvw)
	}

	return wrapper
}

func getAllFieldNames(typ reflect.Type) []string {
	var fields []string
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i).Name)
	}
	return fields
}
