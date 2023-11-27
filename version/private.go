package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gommitizen/git"
)

// Private methods

// Update the version in the files that contain it given by the parameter VersionFiles
func (version *VersionData) updateVersionFiles(v string) error {
	// Update every file that contains the version in VersionFiles with the new version
	for _, file := range version.VersionFiles {
		// Split file into file name and variable by colon
		fileParts := strings.Split(file, ":")
		file := fileParts[0]
		substring := fileParts[1]
		err := updateVersionOfFiles(filepath.Join(version.git.DirPath, file), substring, v)
		if err != nil {
			return fmt.Errorf("Error updating the version in the file %s: %s", file, err)
		}
	}

	return nil
}

// Save the version and commit values in the .version.json file
func (version *VersionData) saveVersion() error {
	// Serializes the updated structure back to JSON
	updatedContent, err := json.MarshalIndent(version, "", "  ")
	if err != nil {
		return &VersionError{
			Message: "Error serializing the updated structure: " + err.Error(),
		}
	}

	// Write the updated content to the file
	err = os.WriteFile(version.filePath, updatedContent, os.ModePerm)
	if err != nil {
		return &VersionError{
			Message: "Error writing the updated content to the file: " + err.Error(),
		}
	}

	return nil
}

// Get the commit stored in the .version.json file
func (version *VersionData) readDataFromJsonFile() error {
	// Read the content of the .version.json file
	content, err := os.ReadFile(version.filePath)
	if err != nil {
		return &VersionError{
			Message: "Error reading file content: " + err.Error(),
		}
	}

	// Deserializes the content into a Version structure
	err = json.Unmarshal(content, version)
	if err != nil {
		return &VersionError{
			Message: "Error deserialize file content: " + err.Error(),
		}
	}

	// Returns the commit value
	return nil
}

// Check if the version file is initialized
func (version *VersionData) checkVersionIsInitialized() error {
	if version.filePath == "" {
		return &VersionError{
			Message: "Error: A .version.json file has not been specified",
		}
	}

	if _, err := os.Stat(version.filePath); os.IsNotExist(err) {
		return &VersionError{
			Message: "Error: the .version.json file does not exist",
		}
	}

	if version.Version == "" || version.Commit == "" {
		return &VersionError{
			Message: "Error: the version and commit values have not been read",
		}
	}

	if version.git == nil {
		return &VersionError{
			Message: "Error: the Git object has not been initialized",
		}
	}

	return nil
}

// return Git object with updated data
func (version *VersionData) returnGitObjectWithUpdatedData() (*git.Git, error) {
	// Get the relative path to the current directory
	relativePath, err := getRelativePath(version.filePath)
	if err != nil {
		return nil, fmt.Errorf("Error obtaining the relative path: %s", err)
	}

	// Get the base path of the file
	dirPath := filepath.Dir(relativePath)

	// Make a Git instance
	git := git.Git{
		DirPath:    dirPath,
		FromCommit: version.Commit,
	}

	// Update Git data
	err = git.UpdateData()
	if err != nil {
		fmt.Println("Error updating Git data:", err)
		return nil, err
	}

	return &git, nil
}

// commitFiles Commit the changes in Git
func (version *VersionData) commitFiles() error {
	// Get the relative path to the current directory
	relativeFilePath, err := getRelativePath(version.filePath)
	if err != nil {
		return fmt.Errorf("Error obtaining the relative path: %s", err)
	}

	// Pay attention to the CHANGELOG.md file and those that host extra versions such as Chart.yaml
	addFiles := []string{}
	addFiles = append(addFiles, relativeFilePath)
	if version.updateChangelog == true {
		addFiles = append(addFiles, filepath.Join(version.git.DirPath, "CHANGELOG.md"))
	}
	for _, file := range version.VersionFiles {
		// Split file into file name and variable by colon
		fileParts := strings.Split(file, ":")
		file := filepath.Join(version.git.DirPath, fileParts[0])
		addFiles = append(addFiles, file)
	}
	commitMessage := "Updated version (" + version.Version + ") in " + getBaseDirFromFilePath(version.git.DirPath)
	tagMessage := version.Version + "_" + getBaseDirFromFilePath(version.git.DirPath)
	output, err := version.git.UpdateGit(addFiles, commitMessage, tagMessage)
	if err != nil {
		fmt.Println("Error updating Git:", err)
		return err
	}

	fmt.Println("")
	for _, line := range output {
		if line != "" {
			fmt.Println(line)
		}
	}

	return nil
}
