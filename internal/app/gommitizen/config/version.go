package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type ConfigVersion struct {
	dirPath string

	Version               string    `json:"version" yaml:"version" plain:"version"`
	Commit                string    `json:"commit" yaml:"commit" plain:"commit"`
	VersionFiles          []string  `json:"version_files" yaml:"version_files" plain:"version_files"`
	Alias                 string    `json:"alias" yaml:"alias" plain:"alias"`
	Hooks                 HookTypes `json:"hooks,omitempty" yaml:"hooks,omitempty" plain:"hooks,omitempty"`
	UpdateChangelogOnBump bool      `json:"update_changelog_on_bump,omitempty" yaml:"update_changelog_on_bump,omitempty" plain:"update_changelog_on_bump,omitempty"`
}

type HookTypes struct {
	PreBump       string `json:"pre_bump,omitempty" yaml:"pre_bump,omitempty" plain:"pre_bump,omitempty"`
	PostBump      string `json:"post_bump,omitempty" yaml:"post_bump,omitempty" plain:"post_bump,omitempty"`
	PreChangelog  string `json:"pre_changelog,omitempty" yaml:"pre_changelog,omitempty" plain:"pre_changelog,omitempty"`
	PostChangelog string `json:"post_changelog,omitempty" yaml:"post_changelog,omitempty" plain:"post_changelog,omitempty"`
}

func NewConfigVersion(dirPath string, version string, commit string, alias string) *ConfigVersion {
	nAlias := alias
	if len(alias) == 0 {
		nAlias = filepath.Base(dirPath)
	}

	return &ConfigVersion{
		dirPath: dirPath,

		Version:      version,
		Commit:       commit,
		VersionFiles: make([]string, 0),
		Alias:        nAlias,
	}
}

func ReadConfigVersion(configVersionPath string) (*ConfigVersion, error) {
	data, err := os.ReadFile(configVersionPath)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %v", configVersionPath, err)
	}

	slog.Debug(fmt.Sprintf("reading config version in %s with data:\n%s", configVersionPath, string(data)))

	var version ConfigVersion
	err = json.Unmarshal(data, &version)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %v", err)
	}

	version.dirPath = filepath.Dir(configVersionPath)
	return &version, nil
}

func (v *ConfigVersion) Save() error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("parse struct to json: %v", err)
	}

	slog.Debug(fmt.Sprintf("saving config version in %s with data:\n%s", v.GetFilePath(), string(data)))

	err = os.WriteFile(v.GetFilePath(), data, 0644)
	if err != nil {
		return fmt.Errorf("write file %s: %v", v.GetFilePath(), err)
	}

	return nil
}

func (v *ConfigVersion) GetFilePath() string {
	return filepath.Join(v.dirPath, defaultFileName)
}

func (v *ConfigVersion) GetDirPath() string {
	return v.dirPath
}

func (v *ConfigVersion) GetGitTag() string {
	if len(v.Alias) > 0 {
		return v.Version + "+" + v.Alias
	}
	return v.Version
}

func (v *ConfigVersion) RunPreBump() error {
	return v.runHook("PreBump")
}

func (v *ConfigVersion) RunPostBump() error {
	return v.runHook("PostBump")
}

func (v *ConfigVersion) RunPreChangelog() error {
	return v.runHook("PreChangelog")
}

func (v *ConfigVersion) RunPostChangelog() error {
	return v.runHook("PostChangelog")
}

func (v *ConfigVersion) runHook(hookName string) error {
	hookValue := reflect.ValueOf(v.Hooks).FieldByName(hookName)
	if !hookValue.IsValid() {
		slog.Debug(fmt.Sprintf("hook %s not found", hookName))
		return nil
	}
	hook := hookValue.String()
	if len(hook) == 0 {
		slog.Debug(fmt.Sprintf("hook %s is empty", hookName))
		return nil
	}

	slog.Debug(fmt.Sprintf("running hook %s: %s", hookName, hook))

	output, err := exec.Command("bash", "-c", hook).CombinedOutput()
	if err != nil {
		return fmt.Errorf("run hook %s: %v", hookName, err)
	}

	// TODO: Pretty log info with colors
	if len(output) > 0 {
		slog.Info(fmt.Sprintf("\n\033[32mHook %s output:\n%s\033[0m", hookName, string(output)))
	} else {
		slog.Info(fmt.Sprintf("\n\033[32mLaunch hook %s\n\033[0m", hookName))
	}

	return nil
}

func (v *ConfigVersion) UpdateVersion(newVersion string, lastCommit string) ([]string, error) {
	modifiedFiles := make([]string, 0)

	v.Version = newVersion
	v.Commit = lastCommit
	err := v.Save()
	if err != nil {
		return nil, err
	}
	modifiedFiles = append(modifiedFiles, v.GetFilePath())

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
