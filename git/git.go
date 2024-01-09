package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manage Git information for our project
type Git struct {
	DirPath        string
	FromCommit     string
	gitPath        string
	lastCommit     string
	changedFiles   []string
	commitMessages []string
	filterFiles    []string
	output         []string
}

// Constructor

// NewGit Create a new Git instance
func NewGit(dirPath string, fromCommit string) *Git {
	var err error

	git := &Git{
		DirPath:    dirPath,
		FromCommit: fromCommit,
	}

	git.gitPath, err = git.getTopLevel()
	if err != nil {
		panic("Error obtaining the top level of the Git repository: " + err.Error())
	}

	return git
}

// Setters
func (git *Git) SetDirPath(dirPath string) {
	var err error

	git.gitPath, err = git.getTopLevel()
	if err != nil {
		panic("Error obtaining the top level of the Git repository: " + err.Error())
	}
	git.DirPath = dirPath
}

func (git *Git) SetGitPath(gitPath string) {
	git.gitPath = gitPath
}

func (git *Git) SetFromCommit(fromCommit string) {
	git.FromCommit = fromCommit
}

func (git *Git) SetFilterFiles(excludedFiles []string) {
	for i := 0; i < len(excludedFiles); i++ {
		git.filterFiles = append(git.filterFiles, filepath.Join(git.DirPath, excludedFiles[i]))
	}
}

// Getters
func (git *Git) GetDirPath() string {
	return git.DirPath
}

func (git *Git) GetChangedFiles() []string {
	return git.changedFiles
}

func (git *Git) GetCommitMessages() []string {
	return git.commitMessages
}

func (git *Git) GetFromCommit() string {
	return git.FromCommit
}

func (git *Git) GetLastCommit() string {
	return git.lastCommit
}

func (git *Git) GetFilterFiles() []string {
	return git.filterFiles
}

func (git *Git) GetOutput() []string {
	return git.output
}

// Public methods

// Initialize a Git directory (git init)
func (git *Git) Initialize(gitPath string) error {
	var err error
	var output []byte

	git.gitPath = gitPath
	git.DirPath = gitPath

	cmd := exec.Command("git", "init", git.gitPath)
	output, err = cmd.Output()

	for _, line := range strings.Split(string(output), "\n") {
		git.output = append(git.output, line)
	}

	if err != nil {
		return fmt.Errorf("error: the Git repository could not be initialized: %v", err)
	}

	git.SetFromCommit("HEAD")
	err = git.RetrieveLastCommit()
	if err != nil {
		return &GitError{
			Message: "Error: the data could not be retrieved: " + err.Error(),
		}
	}

	return nil
}

// RetrieveLastCommit Get the last commit in Git and store it in the lastCommit attribute
func (git *Git) RetrieveLastCommit() error {
	lastCommit, err := git.getLastCommitFromGit()
	if err != nil {
		return err
	}

	git.lastCommit = lastCommit

	return nil
}

// RetrieveData Get the list of modified files in Git from a given commit in a given directory and store them in the changedFiles attribute
// Also retrieves the commit messages and stores them in the commitMessages attribute.
func (git *Git) RetrieveData() error {
	var err error
	var changedFiles []string
	var commitMessages []string
	var lastCommit string

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

	var errorMsg = "the list of modified files could not be retrieved: %v\n\n" +
		"\t** gommitizen requires you to make an initial commit first and\n" +
		"\t   then a conventional commit to establish a reference commit.\n"

	changedFiles, err = git.getListOfModifiedFilesInGitFromAGivenCommitInDirExcludingFiles()
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	git.changedFiles = changedFiles

	commitMessages, err = git.getCommitMessages()
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	git.commitMessages = commitMessages

	lastCommit, err = git.getLastCommitFromGit()
	if err != nil {
		return err
	}
	git.lastCommit = lastCommit

	return nil
}

// ConfirmChanges Update the data of Git
func (git *Git) ConfirmChanges(files []string, commitMessage string, tagMessage string) error {
	var err error

	git.CleanOutput()

	for _, file := range files {
		err = git.Add(file)
		if err != nil {
			return fmt.Errorf("error: the file %s could not be added: %v", file, err)
		}
	}

	err = git.Commit(commitMessage)
	if err != nil {
		return fmt.Errorf("error: the commit could not be created: %v", err)
	}

	err = git.Tag(tagMessage)
	if err != nil {
		return fmt.Errorf("error: the tag could not be created: %v", err)
	}

	return nil
}

// CleanOutput Clean the output of Git
func (git *Git) CleanOutput() {
	git.output = []string{}
}

// Add a file to the Git repository
func (git *Git) Add(filePath string) error {
	cmd := exec.Command("git", "-C", git.gitPath, "add", filePath)
	output, err := cmd.Output()

	for _, line := range strings.Split(string(output), "\n") {
		git.output = append(git.output, line)
	}

	if err != nil {
		return err
	}

	return nil
}

// Commit Make a new commit in Git
func (git *Git) Commit(message string) error {
	cmd := exec.Command("git", "-C", git.gitPath, "commit", "-m", message)
	output, err := cmd.Output()

	for _, line := range strings.Split(string(output), "\n") {
		git.output = append(git.output, line)
	}

	if err != nil {
		return err
	}

	return nil
}

// Tag Make a new tag in Git
func (git *Git) Tag(tag string) error {
	cmd := exec.Command("git", "-C", git.gitPath, "tag", tag)
	output, err := cmd.Output()

	for _, line := range strings.Split(string(output), "\n") {
		git.output = append(git.output, line)
	}

	if err != nil {
		return err
	}

	return nil
}

// Private methods

// Get the list of modified files in Git from a given commit in a given directory
// Allows you to exclude files from the list of modified files
func (git *Git) getListOfModifiedFilesInGitFromAGivenCommitInDirExcludingFiles() ([]string, error) {
	cmd := exec.Command("git", "-C", git.gitPath, "diff", "--name-only", git.FromCommit, "HEAD", git.DirPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Divides the output into lines
	lines := strings.Split(string(output), "\n")

	// Remove the last line (empty)
	lines = lines[:len(lines)-1]

	// Remove excluded files from the list of modified files
	for _, excludeFile := range git.filterFiles {
		lines = RemoveStringFromSlice(lines, excludeFile)
	}

	return lines, nil
}

// Get the commit messages for the modified files in Git from a given commit in a given directory
func (git *Git) getCommitMessages() ([]string, error) {
	// Build the git log command with options to get messages and modified files
	args := append([]string{"-C", git.gitPath, "log", "--pretty=%s", git.FromCommit + "..", "--"}, git.changedFiles...)
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

// Get the current commit in Git
func (git *Git) getLastCommitFromGit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(output), "\n") {
		git.output = append(git.output, line)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetTopLevel Get the top level of the Git repository
func (git *Git) getTopLevel() (string, error) {
	if git.DirPath == "" {
		return "", &GitError{
			Message: "Error: the working directory has not been specified",
		}
	}

	cmd := exec.Command("git", "-C", git.DirPath, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error: the top level could not be retrieved: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}
