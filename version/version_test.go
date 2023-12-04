package version

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gommitizen/git"
	"gommitizen/git/mockGit"
	"path/filepath"
	"testing"
)

func TestVersionData_NewVersionData(t *testing.T) {
	tmpDir := t.TempDir()

	gitHandler := &git.Git{}
	err := gitHandler.Initialize(tmpDir)
	assert.NoError(t, err, "Must not have error when initializing")

	version := NewVersionData("0.0.0", "HEAD^", filepath.Join(tmpDir, ".version.json"), "v")

	assert.Equal(t, filepath.Join(tmpDir, ".version.json"), version.filePath, "The config file should be the default one")
	assert.Equal(t, "HEAD^", version.Commit, "The commit should be the default one")
	assert.Equal(t, "v", version.Prefix, "The prefix should be the default one")
	assert.Equal(t, "0.0.0", version.Version, "The version should be the default one")
}

func TestVersionData_LoadVersionData(t *testing.T) {
	var err error

	tmpDir := t.TempDir()

	gitHandler := &git.Git{}
	err = gitHandler.Initialize(tmpDir)
	assert.NoError(t, err, "Must not have error when initializing")

	version := NewVersionData("0.0.0", "HEAD^", filepath.Join(tmpDir, ".version.json"), "v")

	err = version.Save()
	assert.NoError(t, err, "Must not have error when saving")

	loadedVersion := LoadVersionData(filepath.Join(tmpDir, ".version.json"))
	assert.NotNilf(t, loadedVersion, "The loaded version should not be nil")
	assert.Equal(t, version, loadedVersion, "The loaded version should be the same as the saved one")
}

func TestVersionData_EmptyVersionData(t *testing.T) {
	version := &VersionData{}

	assert.Equal(t, "", version.filePath, "The config file should be the default one")
	assert.Equal(t, "", version.Commit, "The commit should be the default one")
	assert.Equal(t, "", version.Prefix, "The prefix should be the default one")
	assert.Equal(t, "", version.Version, "The version should be the default one")
}

func TestVersionData_Initialize(t *testing.T) {
	version := &VersionData{
		Commit:   "HEAD^",
		Prefix:   "v",
		Version:  "0.0.0",
		filePath: "test",
	}

	tmpDir := t.TempDir()

	err := version.Initialize(tmpDir)
	assert.NoError(t, err, "Must not have error when initializing")

	configFile := filepath.Join(tmpDir, ConfigFileName)
	assert.Equal(t, configFile, version.filePath, "The config file should be the default one")

	expectedCommit := DefaultCommit
	assert.Equal(t, expectedCommit, version.Commit, "The commit should be the default one")

	expectedVersion := DefaultVersionTag
	assert.Equal(t, expectedVersion, version.Version, "The version should be the default one")
}

func TestVersionData_GetVersion(t *testing.T) {
	version := &VersionData{
		Version: "0.0.0",
	}

	assert.Equal(t, "0.0.0", version.GetVersion(), "The version should be the default one")
}

func TestVersionData_GetCommit(t *testing.T) {
	version := &VersionData{
		Commit: "HEAD^",
	}

	assert.Equal(t, "HEAD^", version.GetCommit(), "The commit should be the default one")
}

func TestVersionData_GetFilePath(t *testing.T) {
	version := &VersionData{
		filePath: "test",
	}

	assert.Equal(t, "test", version.GetFilePath(), "The config file should be the default one")
}

func TestVersionData_GetPrefix(t *testing.T) {
	version := &VersionData{
		Prefix: "v",
	}

	assert.Equal(t, "v", version.GetPrefix(), "The prefix should be the default one")
}

func TestVersionData_GetGit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHandler := mockGit.NewMockGitI(ctrl)

	mockGitHandler.EXPECT().GetDirPath().Return("/tmp").AnyTimes()

	version := &VersionData{
		git: mockGitHandler,
	}

	assert.NotNilf(t, version.GetGit(), "The git object should not be nil")
	assert.Equal(t, mockGitHandler.GetDirPath(), version.GetGit().GetDirPath(), "The git directory should be the default one")
}

func TestVersionData_GetUpdateChangelog(t *testing.T) {
	version := &VersionData{
		updateChangelog: true,
	}

	assert.Equal(t, true, version.GetUpdateChangelog(), "The updateChangelog should be the default one")
}

func TestVersionData_SetVersion(t *testing.T) {
	version := &VersionData{
		Version: "0.0.0",
	}

	version.SetVersion("1.0.0")
	assert.Equal(t, "1.0.0", version.Version, "The version should be the default one")
}

func TestVersionData_SetCommit(t *testing.T) {
	version := &VersionData{
		Commit: "HEAD^",
	}

	version.SetCommit("HEAD")
	assert.Equal(t, "HEAD", version.Commit, "The commit should be the default one")
}

func TestVersionData_SetFilePath(t *testing.T) {
	version := &VersionData{
		filePath: "test",
	}

	version.SetFilePath("test2")
	assert.Equal(t, "test2", version.filePath, "The config file should be the default one")
}

func TestVersionData_SetGit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHandler := mockGit.NewMockGitI(ctrl)

	mockGitHandler.EXPECT().GetDirPath().Return("/tmp").AnyTimes()

	version := &VersionData{
		git: nil,
	}

	assert.Nil(t, version.git, "The git object should be nil")

	version.SetGit(mockGitHandler)

	assert.NotNilf(t, version.GetGit(), "The git object should not be nil")
	assert.Equal(t, mockGitHandler.GetDirPath(), version.GetGit().GetDirPath(), "The git directory should be the default one")
}

func TestVersionData_SetUpdateChangelog(t *testing.T) {
	version := &VersionData{
		updateChangelog: true,
	}

	version.SetUpdateChangelog(false)
	assert.Equal(t, false, version.updateChangelog, "The updateChangelog should be the default one")
}
