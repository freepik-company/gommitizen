package version

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gommitizen/changelog"
)

// Constants
const versionFile = ".version.json"

// Public methods

// Initialize the .version.json file
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

// Get the version and commit values from the .version.json file
func (version *VersionData) ReadData(filePath string) error {
	version.filePath = filePath

	// Read the data from the .version.json file
	err := version.readDataFromJsonFile()
	if err != nil {
		return fmt.Errorf("Error reading data from the .version.json file: %s", err)
	}

	// Get a Git object with updated data
	version.git, err = version.returnGitObjectWithUpdatedData()
	if err != nil {
		return fmt.Errorf("Error al actualizar Git: %s", err)
	}

	return nil
}

// Returns true if some file has been modified in Git from a given commit in a given directory
func (version *VersionData) IsSomeFileModified() (bool, error) {
	// Check if the version file is initialized
	err := version.checkVersionIsInitialized()
	if err != nil {
		return false, err
	}

	changedFiles := version.git.GetChangedFiles()

	// Verify if the list of modified files is empty
	return len(changedFiles) > 0, nil
}

// Update the version value in the .version.json file based on the changes in Git
func (version *VersionData) UpdateVersion() (string, error) {
	// Check if the version file is initialized
	err := version.checkVersionIsInitialized()
	if err != nil {
		return "", err
	}

	// Determine the type of version increment based on the commit messages
	incType := determineVersionBump(version.git.CommitMessages)

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
		for _, msg := range version.git.GetCommitMessages() {
			if strings.HasPrefix(msg, "Updated version") {
				continue
			}
			fmt.Println("+", msg)
		}
		fmt.Println()

		// Print the list of changed files in Git
		fmt.Println("Changed files: ")
		for _, file := range version.git.GetChangedFiles() {
			if strings.HasSuffix(file, versionFile) ||
				strings.HasSuffix(file, "CHANGELOG.md") {
				continue
			}
			print := true
			for _, versionFile := range version.VersionFiles {
				versionFileParts := strings.Split(versionFile, ":")
				if strings.HasSuffix(file, versionFileParts[0]) {
					print = false
				}
			}
			if print == true {
				fmt.Println("+", file)
			}
		}
		fmt.Println()

		// Report the version bump, update the version and commit values and update Git
		fmt.Println("Version bumped from " + currentVersion + " to " + newVersion)
		version.Commit = version.git.LastCommit
		version.Version = newVersion

		// Save the updated version and commit values in the .version.json file
		err = version.saveVersion()
		if err != nil {
			return "", &VersionError{
				Message: "Error saving the updated version and commit values in the .version.json file: " + err.Error(),
			}
		}

		// Update file wich contains the version field we want to update
		err = version.updateVersionFiles(newVersion)
		if err != nil {
			return "", &VersionError{
				Message: "Error updating the version in the files: " + err.Error(),
			}
		}

		// Update the CHANGELOG.md file
		if version.updateChangelog == true {
			err = version.UpdateChangelog()
			if err != nil {
				fmt.Println("Error updating the CHANGELOG.md file:", err)
				return "", err
			}
		}

		// Commit the changes in Git
		err = version.commitFiles()
		if err != nil {
			fmt.Println("Error committing changes with Git:", err)
			return "", err
		}
	} else {
		fmt.Printf("Current version: %s (Bump skipped!)\n", currentVersion)
		version.Version = currentVersion
	}

	return version.Version, nil
}

// Increment the version value in the .version.json file based on the given increment type
func (version *VersionData) IncrementVersion(incrementType string) (string, error) {
	incType := strings.ToLower(incrementType)

	if incType != "major" && incType != "minor" && incType != "patch" {
		return "", &VersionError{
			Message: "Error: the increment type must be 'major', 'minor' or 'patch'",
		}
	}

	// Increment the current version
	currentVersion, newVersion, err := incrementVersion(version.Version, incType)
	if err != nil {
		return "", &VersionError{
			Message: "Error incrementing the current version: " + err.Error(),
		}
	}

	// Report the version bump, update the version and commit values and update Git
	fmt.Println("Version bumped from " + currentVersion + " to " + newVersion)

	// Update the version value of the .version.json file
	version.Version = newVersion
	version.Commit = version.git.LastCommit

	// Save the updated version value in the .version.json file
	err = version.saveVersion()
	if err != nil {
		return "", &VersionError{
			Message: "Error saving the updated version value in the .version.json file: " + err.Error(),
		}
	}

	// Update file wich contains the version field we want to update
	err = version.updateVersionFiles(newVersion)
	if err != nil {
		return "", &VersionError{
			Message: "Error updating the version in the files: " + err.Error(),
		}
	}

	// Commit the changes in Git
	err = version.commitFiles()
	if err != nil {
		fmt.Println("Error committing changes with Git:", err)
		return "", err
	}

	return version.Version, nil
}

// UpdateChangelog Update the CHANGELOG.md file based on the changes in Git
func (version *VersionData) UpdateChangelog() error {
	// Check if the version file is initialized
	err := version.checkVersionIsInitialized()
	if err != nil {
		return err
	}

	bcCommits := []string{}
	featCommits := []string{}
	fixCommits := []string{}
	for _, msg := range version.git.GetCommitMessages() {
		// Ignore the commit message that updates the version
		if strings.HasPrefix(msg, "Updated version") {
			continue
		}
		// Sort the commit messages into the corresponding category
		for _, prefix := range bcPrefix {
			if strings.HasPrefix(msg, prefix) {
				bcCommits = append(bcCommits, strings.TrimSpace(strings.TrimPrefix(msg, prefix)))
			}
		}
		for _, prefix := range featPrefix {
			if strings.HasPrefix(msg, prefix) {
				featCommits = append(featCommits, strings.TrimSpace(strings.TrimPrefix(msg, prefix)))
			}
		}
		for _, prefix := range fixPrefix {
			if strings.HasPrefix(msg, prefix) {
				fixCommits = append(fixCommits, strings.TrimSpace(strings.TrimPrefix(msg, prefix)))
			}
		}
	}

	c := changelog.New(version.Version, version.git.DirPath)
	c.BcChanges = bcCommits
	c.FeatChanges = featCommits
	c.FixChanges = fixCommits
	err = c.Write()
	if err != nil {
		return err
	}

	return nil
}
