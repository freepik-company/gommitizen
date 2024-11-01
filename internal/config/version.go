package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ConfigVersion struct {
	dirPath string

	Version      string   `json:"version"`
	Commit       string   `json:"commit"`
	VersionFiles []string `json:"version_files"`
	TagPrefix    string   `json:"tag_prefix"`
}

func NewConfigVersion(dirPath string, version string, commit string, tagPrefix string) *ConfigVersion {
	nDirPath, err := NormalizePath(dirPath)
	if err != nil {
		panic(fmt.Errorf("NormalizePath %s: %v", dirPath, err))
	}

	nTagPrefix := tagPrefix
	if len(tagPrefix) == 0 {
		nTagPrefix = filepath.Base(dirPath)
	}

	return &ConfigVersion{
		dirPath: nDirPath,

		Version:      version,
		Commit:       commit,
		VersionFiles: make([]string, 0),
		TagPrefix:    nTagPrefix,
	}
}

func ReadConfigVersion(configVersionPath string) (*ConfigVersion, error) {
	data, err := os.ReadFile(configVersionPath)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %v", configVersionPath, err)
	}

	var version ConfigVersion
	err = json.Unmarshal(data, &version)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %v", err)
	}

	version.dirPath = filepath.Dir(configVersionPath)
	return &version, nil
}

func (v ConfigVersion) Save() error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("parse struct to json: %v", err)
	}

	err = os.WriteFile(v.GetConfigVersionFilePath(), data, 0644)
	if err != nil {
		return fmt.Errorf("write file %s: %v", v.GetConfigVersionFilePath(), err)
	}

	return nil
}

func (v ConfigVersion) GetConfigVersionFilePath() string {
	return filepath.Join(v.dirPath, defaultFileName)
}

func (v ConfigVersion) GetDirPath() string {
	return v.dirPath
}

func (v ConfigVersion) GetTagVersion() string {
	if len(v.TagPrefix) > 0 {
		return v.TagPrefix + "_" + v.Version
	}
	return v.Version
}

func (v *ConfigVersion) UpdateVersion(newVersion string, lastCommit string) ([]string, error) {
	modifiedFiles := make([]string, 0)

	v.Version = newVersion
	v.Commit = lastCommit
	err := v.Save()
	if err != nil {
		return nil, err
	}
	modifiedFiles = append(modifiedFiles, v.GetConfigVersionFilePath())

	for _, versionFile := range v.VersionFiles {
		index := strings.Index(versionFile, ":")
		if index == -1 {
			slog.Info(fmt.Sprintf("warning, `%s` is not a valid format", versionFile))
			continue
		}

		fileName := versionFile[:index]
		substring := versionFile[index+1:]
		filePath := filepath.Join(v.dirPath, fileName)

		err := updateVersionOfFiles(filePath, substring, newVersion)
		if err != nil {
			return nil, err
		}
		modifiedFiles = append(modifiedFiles, filePath)
	}

	return modifiedFiles, nil
}

func updateVersionOfFiles(filePath, substring, newVersion string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	// Regular expression to find the version in the file
	regularExpression := ""
	validRegexp, err := isARegExp(substring)
	// Check if the substring is a regular expression that compiles
	if err != nil {
		return err
	}
	if validRegexp { // If it is a regular expression, use it as is
		regularExpression = substring
	} else { // If it is a literal string, use it as a word boundary
		regularExpression = fmt.Sprintf(`(?i)\b%s\b\s*[:=]\s*([0-9]+\.[0-9]+\.[0-9]+)`, substring)
	}
	versionRegex, err := regexp.Compile(regularExpression)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		match := versionRegex.FindStringSubmatch(line)
		if len(match) > 1 {
			// The line contains the substring given by 'substring'
			oldVersion := strings.TrimSpace(match[1])
			// Replace the old version with the new version
			newLine := strings.Replace(line, oldVersion, newVersion, 1)
			lines = append(lines, newLine)
		} else {
			// The line does not contain the substring given by 'substring'
			lines = append(lines, line)
		}
	}

	// Write the updated lines back to the file
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range lines {
		file.WriteString(line + "\n")
	}

	return nil
}

func isARegExp(s string) (bool, error) {
	// Compile the regular expression
	_, err := regexp.Compile(s)
	if err != nil {
		return false, err
	}
	// Check if the string contains any special regex characters
	specialChars := `.*+?^${}()|[]\`
	for _, char := range specialChars {
		if strings.ContainsRune(s, char) {
			return true, nil
		}
	}
	return false, nil
}
