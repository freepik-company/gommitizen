package git

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGit_Initialize(t *testing.T) {
	git := &Git{}

	tmpDir := t.TempDir()

	err := git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")
}

func TestGit_RetrieveLastCommit(t *testing.T) {
	git := &Git{}

	tmpDir := t.TempDir()

	err := git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	lastCommit := git.GetLastCommit()
	assert.NotEmpty(t, lastCommit, "Must have a commit")
}

func TestGit_RetrieveData(t *testing.T) {
	var err error

	git := &Git{
		FromCommit: "HEAD",
	}

	tmpDir := t.TempDir()

	err = git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	tmpFile := filepath.Join(tmpDir, "tempFile")
	err = os.WriteFile(tmpFile, []byte("temp"), 0644)
	assert.Nil(t, err, "Must not have error when creating the file")

	err = git.Add(tmpFile)
	assert.Nil(t, err, "Must not have error when adding the file")

	err = git.Commit("test")
	assert.Nil(t, err, "Must not have error when committing")

	err = git.RetrieveData()
	assert.Nil(t, err, "Must not have error when retrieving data")
}

func TestGit_ConfirmChanges(t *testing.T) {
	var err error

	git := &Git{
		FromCommit: "HEAD",
	}

	tmpDir := t.TempDir()

	err = git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	tmpFile1 := filepath.Join(tmpDir, "tempFile")
	err = os.WriteFile(tmpFile1, []byte("temp1"), 0644)
	assert.Nil(t, err, "Must not have error when creating the first file")

	tmpFile2 := filepath.Join(tmpDir, "tempFile")
	err = os.WriteFile(tmpFile2, []byte("temp2"), 0644)
	assert.Nil(t, err, "Must not have error when creating the second file")

	err = git.ConfirmChanges([]string{tmpFile1, tmpFile2}, "test", "0.0.1")
	assert.Nil(t, err, "Must not have error when confirming changes")
}

func TestGit_CleanOutput(t *testing.T) {
	git := &Git{}

	tmpDir := t.TempDir()

	err := git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	output := git.GetOutput()
	for _, line := range output {
		fmt.Println(line)
	}

	git.CleanOutput()
	assert.Equal(t, 4, len(output), "Should be 4 lines")
}

func TestGit_Add(t *testing.T) {
	var err error

	git := &Git{}

	tmpDir := t.TempDir()

	err = git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	tmpFile := filepath.Join(tmpDir, "tempFile")
	err = os.WriteFile(tmpFile, []byte("temp"), 0644)
	assert.Nil(t, err, "Must not have error when creating the file")

	err = git.Add(tmpFile)
	assert.Nil(t, err, "Must not have error when adding the file")
}

func TestGit_Commit(t *testing.T) {
	var err error

	git := &Git{}

	tmpDir := t.TempDir()

	err = git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	tmpFile := filepath.Join(tmpDir, "tempFile")
	err = os.WriteFile(tmpFile, []byte("temp"), 0644)
	assert.Nil(t, err, "Must not have error when creating the file")

	err = git.Add(tmpFile)
	assert.Nil(t, err, "Must not have error when adding the file")

	err = git.Commit("test")
	assert.Nil(t, err, "Must not have error when committing")
}

func TestGit_Tag(t *testing.T) {
	var err error

	git := &Git{}

	tmpDir := t.TempDir()

	err = git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	tmpFile := filepath.Join(tmpDir, "tempFile")
	err = os.WriteFile(tmpFile, []byte("temp"), 0644)
	assert.Nil(t, err, "No error should occur when creating the file")

	err = git.Add(tmpFile)
	assert.Nil(t, err, "No error should occur when adding the file")

	err = git.Commit("test")
	assert.Nil(t, err, "No error should occur when committing")

	err = git.Tag("v0.0.1")
	assert.Nil(t, err, "No error should occur when tagging")
}

func TestGit_SetDirPath(t *testing.T) {
	git := &Git{}

	tmpDir := t.TempDir()

	err := git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	assert.Equal(t, tmpDir, git.GetDirPath(), "The directory paths should be the same")
}

func TestGit_SetFilterFiles(t *testing.T) {
	git := &Git{}

	tmpDir := t.TempDir()

	err := git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	git.SetFilterFiles([]string{
		"file1.txt",
		"file2.txt",
	})

	assert.Equal(t, []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}, git.GetFilterFiles(), "The files should be the same")
}

func TestGit_SetFromCommit(t *testing.T) {
	git := &Git{}

	git.SetFromCommit("HEAD")

	assert.Equal(t, "HEAD", git.GetFromCommit(), "The commits should be the same")
}

func TestGit_GetChangedFiles(t *testing.T) {
	var err error

	git := &Git{}

	tmpDir := t.TempDir()

	err = git.Initialize(tmpDir)
	assert.Nil(t, err, "No error should occur when initializing")

	git.SetFromCommit("HEAD^")

	tmpFile := filepath.Join(tmpDir, "tempFile")

	// First Commit
	err = os.WriteFile(tmpFile, []byte("temp1"), 0644)
	assert.Nil(t, err, "No error should occur when creating the file")
	err = git.Add(tmpFile)
	assert.Nil(t, err, "No error should occur when adding the file")
	err = git.Commit("test1")
	assert.Nil(t, err, "No error should occur when committing")

	// Second Commit
	err = os.WriteFile(tmpFile, []byte("temp2"), 0644)
	assert.Nil(t, err, "No error should occur when creating the file")
	err = git.Add(tmpFile)
	assert.Nil(t, err, "No error should occur when adding the file")
	err = git.Commit("test2")
	assert.Nil(t, err, "No error should occur when committing")

	err = git.RetrieveData()
	assert.Nil(t, err, "No error should occur when retrieving data")

	changedFiles := git.GetChangedFiles()
	assert.Equal(t, []string{strings.TrimPrefix(tmpFile, tmpDir+"/")}, changedFiles, "The files should be the same")
}

func TestGit_GetCommitMessages(t *testing.T) {
	var err error

	git := &Git{}

	tmpDir := t.TempDir()

	err = git.Initialize(tmpDir)
	assert.Nil(t, err, "No error should occur when initializing")

	git.SetFromCommit("HEAD^")

	tmpFile := filepath.Join(tmpDir, "tempFile")

	// First Commit
	err = os.WriteFile(tmpFile, []byte("temp1"), 0644)
	assert.Nil(t, err, "No error should occur when creating the file")
	err = git.Add(tmpFile)
	assert.Nil(t, err, "No error should occur when adding the file")
	err = git.Commit("Fist commit: This commit message won't be retrieved")
	assert.Nil(t, err, "No error should occur when committing")

	// Second Commit
	err = os.WriteFile(tmpFile, []byte("temp2"), 0644)
	assert.Nil(t, err, "No error should occur when creating the file")
	err = git.Add(tmpFile)
	assert.Nil(t, err, "No error should occur when adding the file")
	err = git.Commit("The actual commit message: This commit message will be retrieved")
	assert.Nil(t, err, "No error should occur when committing")

	err = git.RetrieveData()
	assert.Nil(t, err, "No error should occur when retrieving data")

	commitMessages := git.GetCommitMessages()
	assert.Equal(t, []string{"The actual commit message: This commit message will be retrieved"}, commitMessages, "The commit messages should be the same")
}

func TestGit_GetFromCommit(t *testing.T) {
	git := &Git{}

	git.SetFromCommit("HEAD")

	assert.Equal(t, "HEAD", git.GetFromCommit(), "The hashes should be the same")
}

func TestGit_GetLastCommit(t *testing.T) {
	git := &Git{}

	tmpDir := t.TempDir()

	err := git.Initialize(tmpDir)
	assert.Nil(t, err, "No error should occur when initializing")

	lastCommit := git.GetLastCommit()
	assert.NotEmpty(t, lastCommit, "Must have a commit")
}

func TestGit_GetFilterFiles(t *testing.T) {
	git := &Git{}

	git.SetFilterFiles([]string{
		"file1.txt",
		"file2.txt",
	})

	assert.Equal(t, []string{
		"file1.txt",
		"file2.txt",
	}, git.GetFilterFiles(), "The files should be the same")
}

func TestGit_GetOutput(t *testing.T) {
	git := &Git{}

	tmpDir := t.TempDir()

	err := git.Initialize(tmpDir)
	assert.Nil(t, err, "Must not have error when initializing")

	output := git.GetOutput()
	assert.Equal(t, 4, len(output), "Should be 4 lines")
}
