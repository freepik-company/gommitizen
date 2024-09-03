package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Private methods

// Update the Version in the files that contain it given by the parameter VersionFiles
func (version *VersionData) updateVersionFiles(v string) error {
	// Update every file that contains the version in VersionFiles with the new version
	for _, versionFile := range version.VersionFiles {
		// Split file into file name and variable by colon

		index := strings.Index(versionFile, ":")
		if index == -1 {
			fmt.Printf("warning, `%s` is not a valid format", versionFile)
			continue
		}

		fileName := versionFile[:index]
		substring := versionFile[index+1:]

		err := updateVersionOfFiles(filepath.Join(version.git.GetDirPath(), fileName), substring, v)
		if err != nil {
			return fmt.Errorf("Error updating the version in the file %s: %s", fileName, err)
		}
	}

	return nil
}

// Save the Version and Commit values in the .version.json file
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

// Get the Commit stored in the .version.json file
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

// Check if the Version file is initialized
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
func (version *VersionData) loadGitObjectWithUpdatedData() error {
	var err error

	// Update Git data
	err = version.git.RetrieveData()
	if err != nil {
		fmt.Println("Error updating Git data:", err)
		return err
	}

	return nil
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
	dirPath := version.git.GetDirPath()
	if version.updateChangelog == true {
		addFiles = append(addFiles, filepath.Join(dirPath, "CHANGELOG.md"))
	}
	for _, file := range version.VersionFiles {
		// Split file into file name and variable by colon
		fileParts := strings.Split(file, ":")
		file := filepath.Join(version.git.GetDirPath(), fileParts[0])
		addFiles = append(addFiles, file)
	}
	prj := getBaseDirFromFilePath(dirPath)
	var commitMessage string
	var tagMessage string
	if prj == "." { // root project
		commitMessage = "Updated version (" + version.Version + ")"
		tagMessage = version.Version
	} else { // subproject
		commitMessage = "Updated version (" + version.Version + ") in " + prj
		tagMessage = version.Version + "_" + prj
	}
	err = version.git.ConfirmChanges(addFiles, commitMessage, tagMessage)
	if err != nil {
		fmt.Println("Error updating Git:", err)
		return err
	}

	fmt.Println("")
	output := version.git.GetOutput()
	for _, line := range output {
		if line != "" {
			fmt.Println(line)
		}
	}

	return nil
}
