package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"

	"gommitizen/git"
)

// Personalized error type

const versionFile = ".version.json"

type VersionError struct {
	Message string
}

func (e *VersionError) Error() string {
	return e.Message
}

// Manage the version information for our project

type VersionData struct {
	Version  string `json:"version"`
	Commit   string `json:"commit"`
	filePath string
}

func NewVersionData(version string, commit string, filePath string) *VersionData {
	return &VersionData{Version: version, Commit: commit, filePath: filePath}
}

// Getters

func (version *VersionData) GetVersion() string {
	return version.Version
}

func (version *VersionData) GetCommit() string {
	return version.Commit
}

func (version *VersionData) GetFilePath() string {
	return version.filePath
}

// Public functions

func (version *VersionData) Initialize(path string) error {
	// check .version.json does not exist
	configFile := path + "/.version.json"
	if _, err := os.Stat(configFile); err == nil {
		fmt.Println("The repository is already initialized")
		os.Exit(1)
	}

	version.Commit = "HEAD^"
	version.Version = "0.0.0"
	version.filePath = configFile

	err := version.Save()
	if err != nil {
		fmt.Println("Error saving .version.json file:", err)
		return err
	}

	return nil
}

// Save the version and commit values in the .version.json file
func (version *VersionData) Save() error {
	jsonData, err := version.String()

	err = os.WriteFile(version.filePath, []byte(jsonData), 0644)

	if err != nil {
		return err
	}

	return nil
}

// String returns the JSON representation of the VersionData struct
func (version *VersionData) String() (string, error) {
	jsonData, err := json.MarshalIndent(version, "", "  ")

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// Find all .version.json files in a given directory and its subdirectories
func FindFCVersionFiles(rootDir string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), versionFile) {
			fileList = append(fileList, path)
		}

		return nil
	})

	return fileList, err
}

// Public methods

// Get the version and commit values from the .version.json file
func (version *VersionData) ReadData(filePath string) error {
	version.filePath = filePath

	ver, errVersion := version.getCurrentVersionFromJsonFile()
	if errVersion != nil {
		return errVersion
	}
	version.Version = ver

	commit, errCommit := version.getCommitValueFromJsonFile()
	if errCommit != nil {
		return errCommit
	}
	version.Commit = commit

	return nil
}

// Returns true if some file has been modified in Git from a given commit in a given directory
func (version *VersionData) IsSomeFileModified() (bool, error) {
	if version.Commit == "" || version.filePath == "" {
		return false, &VersionError{
			Message: "A commit or a .version.json file has not been specified",
		}
	}

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return false, &VersionError{
			Message: "Error obtaining current directory",
		}
	}

	// Get the relative path to the current directory
	relativePath, err := filepath.Rel(currentDir, version.filePath)
	if err != nil {
		return false, &VersionError{
			Message: "Error obtaining relative path",
		}
	}

	// Get the base path of the file
	dirPath := filepath.Dir(relativePath)

	// Get the list of modified files in Git from a given commit in a given directory
	git := git.Git{
		DirPath:    dirPath,
		FromCommit: version.Commit,
	}
	errUpdate := git.UpdateData()
	if errUpdate != nil {
		return false, &VersionError{
			Message: "Error updating Git data: " + errUpdate.Error(),
		}
	}
	changedFiles := git.GetChangedFiles()

	// Verify if the list of modified files is empty
	return len(changedFiles) > 0, nil
}

// Update the version value in the .version.json file based on the changes in Git
func (version *VersionData) UpdateVersion() (string, error) {
	if version.Version == "" || version.Commit == "" {
		return "", &VersionError{
			Message: "Error: a commit and version values have not been specified",
		}
	}

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", &VersionError{
			Message: "Error obtaining current directory",
		}
	}

	// Get the relative path to the current directory
	relativePath, err := filepath.Rel(currentDir, version.filePath)
	if err != nil {
		return "", &VersionError{
			Message: "Error obtaining relative path",
		}
	}

	// Get the base path of the file
	dirPath := filepath.Dir(relativePath)

	// Make a Git instance
	git := git.Git{
		DirPath:    dirPath,
		FromCommit: version.Commit,
	}

	// Update the data of Git
	errUpdate := git.UpdateData()
	if errUpdate != nil {
		return "", &VersionError{
			Message: "Error updating Git data: " + errUpdate.Error(),
		}
	}

	// Obtains the type of version increment based on the commit messages
	incType := determineVersionBump(git.CommitMessages)

	// Increment the current version
	currentVersion, newVersion, err := incrementVersion(version.Version, incType)
	if err != nil {
		return "", &VersionError{
			Message: "Error incrementing the current version: " + err.Error(),
		}
	}

	if incType != "none" {
		// Print the list of commit messages
		fmt.Println("Commit messages: ")
		for _, msg := range git.GetCommitMessages() {
			if strings.HasPrefix(msg, "Updated version") {
				continue
			}
			fmt.Println("+", msg)
		}
		fmt.Println()

		// Print the list of files changed in Git
		fmt.Println("Files changed: ")
		for _, file := range git.GetChangedFiles() {
			if strings.HasSuffix(file, versionFile) {
				continue
			}
			fmt.Println("+", file)
		}
		fmt.Println()

		// Report the version bump, update the version and commit values and update Git
		fmt.Println("Version bumped from " + currentVersion + " to " + newVersion)
		version.Commit = git.LastCommit
		version.Version = newVersion

		// Serializes the updated structure back to JSON
		updatedContent, err := json.MarshalIndent(version, "", "  ")
		if err != nil {
			return "", &VersionError{
				Message: "Error serializing the updated structure: " + err.Error(),
			}
		}

		// Write the updated content to the file
		err = os.WriteFile(version.filePath, updatedContent, os.ModePerm)
		if err != nil {
			return "", &VersionError{
				Message: "Error writing the updated content to the file: " + err.Error(),
			}

		}

		addFiles := []string{relativePath}
		commitMessage := "Updated version (" + version.Version + ") in " + getBaseDirFromFilePath(git.DirPath)
		tagMessage := version.Version + "_" + getBaseDirFromFilePath(git.DirPath)
		output, err := git.UpdateGit(addFiles, commitMessage, tagMessage)
		if err != nil {
			return "", &VersionError{
				Message: "Error updating Git: " + err.Error(),
			}

		}

		// Print the output of the git command
		lastLine := output[len(output)-1]
		if strings.TrimSpace(lastLine) == "" {
			output = output[:len(output)-1]
		}
		for _, file := range output {
			fmt.Println(file)
		}

	} else {
		fmt.Printf("Current version: %s (Bump skipped!)\n", currentVersion)
		version.Version = currentVersion
	}

	return version.Version, nil
}

// Private methods

// Get the commit stored in the .version.json file
func (version *VersionData) getCommitValueFromJsonFile() (string, error) {
	// Read the content of the .version.json file
	content, err := os.ReadFile(version.filePath)
	if err != nil {
		return "", &VersionError{
			Message: "Error reading file content: " + err.Error(),
		}
	}

	// Deserializes the content into a Version structure
	err = json.Unmarshal(content, version)
	if err != nil {
		return "", &VersionError{
			Message: "Error deserialize file content: " + err.Error(),
		}
	}

	// Returns the commit value
	return version.Commit, nil
}

// Get the commit stored in the .version.json file
func (version *VersionData) getCurrentVersionFromJsonFile() (string, error) {
	// Read the content of the .version.json file
	content, err := os.ReadFile(version.filePath)
	if err != nil {
		return "", &VersionError{
			Message: "Error reading file content: " + err.Error(),
		}
	}

	// Desializes the content into a Version structure
	err = json.Unmarshal(content, version)
	if err != nil {
		return "", &VersionError{
			Message: "Error deserialize file content: " + err.Error(),
		}

	}

	// Returns the version value
	return version.Version, nil
}

// Private auxiliary functions

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
		// A message contains at the beginning of the given string the following prefix "feat:", "fix:" or "BREAKING CHANGE:"
		if strings.HasPrefix(message, "BREAKING CHANGE:") ||
			strings.HasPrefix(message, "breaking change:") ||
			strings.HasPrefix(message, "Breaking change:") ||
			strings.HasPrefix(message, "bc:") ||
			strings.HasPrefix(message, "BC:") ||
			strings.HasPrefix(message, "Bc:") {
			major = true
		} else if strings.HasPrefix(message, "feat:") ||
			strings.HasPrefix(message, "Feat:") ||
			strings.HasPrefix(message, "feature:") ||
			strings.HasPrefix(message, "Feature:") ||
			strings.HasPrefix(message, "FEAT") {
			minor = true
		} else if strings.HasPrefix(message, "fix:") ||
			strings.HasPrefix(message, "Fix:") ||
			strings.HasPrefix(message, "FIX") ||
			strings.HasPrefix(message, "bug:") ||
			strings.HasPrefix(message, "Bug:") ||
			strings.HasPrefix(message, "BUG") {
			patch = true
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
