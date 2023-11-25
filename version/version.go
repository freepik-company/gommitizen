package version

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"

	"gommitizen/changelog"
	"gommitizen/git"
)

// Variables
var bcPrefix = []string{"BREAKING CHANGE:", "breaking change:", "Breaking change:", "bc:", "BC:", "Bc:"}
var featPrefix = []string{"feat:", "Feat:", "feature:", "Feature:", "FEAT"}
var fixPrefix = []string{"fix:", "Fix:", "FIX", "bug:", "Bug:", "BUG"}

// Constants

const versionFile = ".version.json"

// Personalized error type

type VersionError struct {
	Message string
}

func (e *VersionError) Error() string {
	return e.Message
}

// Manage the version information for our project

type VersionData struct {
	Version         string   `json:"version"`
	Commit          string   `json:"commit"`
	VersionFiles    []string `json:"version_files"`
	filePath        string
	git             *git.Git
	updateChangelog bool
}

func NewVersionData() *VersionData {
	return &VersionData{Version: "", Commit: "", filePath: "", git: nil, updateChangelog: false}
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

func (version *VersionData) GetGit() *git.Git {
	return version.git
}

func (version *VersionData) GetUpdateChangelog() bool {
	return version.updateChangelog
}

// Métodos setter
func (version *VersionData) SetVersion(v string) {
	version.Version = v
}

func (version *VersionData) SetCommit(c string) {
	version.Commit = c
}

func (version *VersionData) SetFilePath(fp string) {
	version.filePath = fp
}

func (version *VersionData) SetGit(g *git.Git) {
	version.git = g
}

func (version *VersionData) SetUpdateChangelog(uc bool) {
	version.updateChangelog = uc
}

// Public methods

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

		// Update the CHANGELOG.md file
		if version.updateChangelog == true {
			err = version.UpdateChangelog()
			if err != nil {
				fmt.Println("Error updating the CHANGELOG.md file:", err)
				return "", err
			}
		}

		// Update every file that contains the version in VersionFiles with the new version
		for _, file := range version.VersionFiles {
			// Split file into file name and variable by colon
			fileParts := strings.Split(file, ":")
			file := fileParts[0]
			substring := fileParts[1]
			err = updateVersionOfFiles(filepath.Join(version.git.DirPath, file), substring, newVersion)
			if err != nil {
				fmt.Println("Error updating the version in the file", file, ":", err)
				return "", err
			}
		}

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

// Private methods

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
