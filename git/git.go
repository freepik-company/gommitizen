package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Managa custom errors
type GitError struct {
	Message string
}

func (e *GitError) Error() string {
	return e.Message
}

// Manage Git information for our project
type Git struct {
	DirPath        string
	FromCommit     string
	LastCommit     string
	ChangedFiles   []string
	CommitMessages []string
	ExcludedFiles  []string
}

// Setters
func (git *Git) SetDirPath(dirPath string) {
	git.DirPath = dirPath
}

func (git *Git) SetFromCommit(fromCommit string) {
	git.FromCommit = fromCommit
}

func (git *Git) setExcludedFiles(excludedFiles []string) {
	for i := 0; i < len(excludedFiles); i++ {
		git.ExcludedFiles = append(git.ExcludedFiles, filepath.Join(git.DirPath, excludedFiles[i]))
	}
}

func (git *Git) setLastCommit() error {
	lastCommit, err := getLastCommitFromGit()
	if err != nil {
		return err
	}

	git.LastCommit = lastCommit

	return nil
}

// Getters
func (git *Git) GetChangedFiles() []string {
	return git.ChangedFiles
}

func (git *Git) GetCommitMessages() []string {
	return git.CommitMessages
}

func (git *Git) GetFromCommit() string {
	return git.FromCommit
}

func (git *Git) GetLastCommit() string {
	return git.LastCommit
}

func (git *Git) GetExcludedFiles() []string {
	return git.ExcludedFiles
}

// Public methods

// Get the list of modified files in Git from a given commit in a given directory and store them in the ChangedFiles attribute
// Also retrieves the commit messages and stores them in the CommitMessages attribute.
func (git *Git) UpdateData() error {
	if git.DirPath == "" {
		return &GitError{
			Message: "Error: the working directory has not been specified",
		}
	}

	if git.FromCommit == "" {
		return &GitError{
			Message: "Error: the commit value has not been specified",
		}
	}

	changedFiles, err := git.getListOfModifiedFilesInGitFromAGivenCommitInDirExcludingFiles()
	if err != nil {
		return err
	}
	git.ChangedFiles = changedFiles

	commitMesages, err := git.getCommitMessages()
	if err != nil {
		return err
	}
	git.CommitMessages = commitMesages

	lastCommit, err := getLastCommitFromGit()
	if err != nil {
		return err
	}
	git.LastCommit = lastCommit

	return nil
}

// Update the data of Git
func (git *Git) UpdateGit(files []string, commitMessage string, tagMessage string) ([]string, error) {
	output := []string{}

	for _, file := range files {
		outputAdd, errAdd := add(file)
		if errAdd != nil {
			return nil, errAdd
		}
		output = append(output, outputAdd)
	}

	outputCommit, errCommit := commit(commitMessage)
	if errCommit != nil {
		return nil, errCommit
	}
	output = append(output, outputCommit)

	outputTag, errTag := tag(tagMessage)
	if errTag != nil {
		return nil, errTag
	}
	output = append(output, outputTag)

	return output, nil
}

// Private methods

// Get the list of modified files in Git from a given commit in a given directory
// Allows you to exclude files from the list of modified files
func (git *Git) getListOfModifiedFilesInGitFromAGivenCommitInDirExcludingFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", git.FromCommit, "HEAD", git.DirPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Divides the output into lines
	lines := strings.Split(string(output), "\n")

	// Remove the last line (empty)
	lines = lines[:len(lines)-1]

	// Remvoe excluded files from the list of modified files
	for _, excludeFile := range git.ExcludedFiles {
		lines = removeStringFromSlice(lines, excludeFile)
	}

	return lines, nil
}

// Get the commit messages for the modified files in Git from a given commit in a given directory
func (git *Git) getCommitMessages() ([]string, error) {
	// Build the git log command with options to get messages and modified files
	args := append([]string{"log", "--pretty=%s", git.FromCommit + "..", "--"}, git.ChangedFiles...)
	cmd := exec.Command("git", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Divide the output into lines and remove empty lines
	lines := strings.Split(string(output), "\n")
	var CommitMessages []string

	// Loop over the lines to build the list of commit messages
	for i := 0; i < len(lines); i++ {
		message := lines[i]
		if strings.TrimSpace(message) != "" {
			// Add the message to the list
			CommitMessages = append(CommitMessages, message)
		}
	}

	return CommitMessages, nil
}

// Private functions

// Add a file to the Git repository
func add(filePath string) (string, error) {
	cmd := exec.Command("git", "add", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Make a new commit in Git
func commit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Make a new tag in Git
func tag(tag string) (string, error) {
	cmd := exec.Command("git", "tag", tag)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Get the current commit in Git
func getLastCommitFromGit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// Auxiliar functions

// Remove a string from a slice of strings
func removeStringFromSlice(slice []string, s string) []string {
	var result []string

	for _, str := range slice {
		if str != s {
			result = append(result, str)
		}
	}

	return result
}
