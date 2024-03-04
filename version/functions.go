package version

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

// Public functions

// FindFCVersionFiles Find .version.json files in a given directory and its subdirectories
func FindFCVersionFiles(rootDir string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".version.json") {
			fileList = append(fileList, path)
		}

		return nil
	})

	return fileList, err
}

// Private auxiliary functions

// Update the version in the files that contain it
func updateVersionOfFiles(filePath, substring, newVersion string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	// Regular expression to find the version in the file
	regularExpression := fmt.Sprintf(`(?i)\b%s\b\s*[:=]\s*([0-9]+\.[0-9]+\.[0-9]+)`, substring)
	versionRegex := regexp.MustCompile(regularExpression)

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

// Get the relative path to the current directory
func getRelativePath(filePath string) (string, error) {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error getting the current directory: %s", err)
	}

	// Get the relative path to the current directory
	relativePath, err := filepath.Rel(currentDir, filePath)
	if err != nil {
		return "", fmt.Errorf("Error getting the relative path: %s", err)
	}

	return relativePath, nil
}

// Get the base directory of a given file
func getBaseDirFromFilePath(filePath string) string {
	return filepath.Base(filePath)
}

// Determine the type of version increment based on the commit messages
func determineVersionBump(commitMessages []string) string {
	major := false
	minor := false
	patch := false

	for _, message := range commitMessages {
		if strings.HasPrefix(message, "Updated version") {
			continue
		}
		for _, prefix := range bcPrefix {
			if strings.HasPrefix(message, prefix) {
				major = true
				break
			}
		}
		for _, prefix := range featPrefix {
			if strings.HasPrefix(message, prefix) {
				minor = true
				break
			}
		}
		for _, prefix := range fixPrefix {
			if strings.HasPrefix(message, prefix) {
				patch = true
				break
			}
		}
		for _, prefix := range refactorPrefix {
			if strings.HasPrefix(message, prefix) {
				patch = true
				break
			}
		}
	}

	if major {
		return "major"
	} else if minor {
		return "minor"
	} else if patch {
		return "patch"
	}

	return "none"
}

// Increment the current version based on the given increment type and returns the new version
func incrementVersion(version string, incType string) (string, string, error) {
	currentVersion, err := semver.NewVersion(version)
	if err != nil {
		return "", "", err
	}

	var newVersion semver.Version
	if incType == "major" {
		newVersion = currentVersion.IncMajor() // Increment the major (for example, from 1.2.3 to 2.0.0)
	} else if incType == "minor" {
		newVersion = currentVersion.IncMinor() // Increment the minor (for example, from 1.2.3 to 1.3.0)
	} else if incType == "patch" {
		newVersion = currentVersion.IncPatch() // Increment the patch (for example, from 1.2.3 to 1.2.4)
	} else {
		newVersion = *currentVersion
	}

	return currentVersion.String(), newVersion.String(), nil
}
