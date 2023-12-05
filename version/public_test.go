package version

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gommitizen/git/mockGit"
	"os"
	"path/filepath"
	"testing"
)

func TestVersionData_Save(t *testing.T) {
	version := &VersionData{
		Commit:   "abc123",
		Prefix:   "v",
		Version:  "1.0.0",
		filePath: "/tmp/.version.json",
	}

	err := version.Save()
	assert.NoError(t, err, "Must not have error when saving")

	_, err = os.Stat(version.filePath)
	assert.Equal(t, nil, err, "The config file should exist")

	fileContent, err := os.ReadFile(version.filePath)
	assert.NoError(t, err, "Must not have error when reading the file content")

	expectedJSON, _ := version.String()
	assert.Equal(t, expectedJSON, string(fileContent), "The file content should be the same as the expected JSON")
}

func TestVersionData_String(t *testing.T) {
	// Crear una instancia de VersionData para probar
	version := &VersionData{
		Commit:  "abc123",
		Prefix:  "v",
		Version: "1.0.0",
		VersionFiles: []string{
			"file1.txt",
			"file2.txt",
		},
	}

	jsonString, err := version.String()
	assert.NoError(t, err, "Must not have error when getting the JSON string")

	var parsedVersion VersionData
	err = json.Unmarshal([]byte(jsonString), &parsedVersion)
	assert.NoError(t, err, "Must not have error when parsing the JSON string")

	assert.Equal(t, version.Commit, parsedVersion.Commit, "The 'Commit' fields do not match")
	assert.Equal(t, version.Version, parsedVersion.Version, "The 'Version' fields do not match")
	assert.Equal(t, version.Prefix, parsedVersion.Prefix, "The 'Prefix' fields do not match")
	assert.Equal(t, version.VersionFiles, parsedVersion.VersionFiles, "The 'VersionFiles' fields do not match")
}

func TestVersionData_IsSomeFileModified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHandler := mockGit.NewMockGitI(ctrl)
	mockGitHandler.EXPECT().GetChangedFiles().Return([]string{
		"file1.txt",
		"file2.txt",
	})

	version := VersionData{
		filePath: ".version.json",
		Commit:   DefaultCommit,
		Version:  DefaultVersionTag,
		Prefix:   "v",

		git: mockGitHandler,
	}

	tmpDir := t.TempDir()
	err := version.Initialize(tmpDir)
	assert.NoError(t, err, "Must not have error when initializing")

	isModified, err := version.IsSomeFileModified()
	assert.NoError(t, err, "Must not have error when checking if there are modified files")
	assert.Truef(t, isModified, "Must be true when there are modified files")
}

func TestVersionData_UpdateVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHandler := mockGit.NewMockGitI(ctrl)

	testCases := []struct {
		name     string
		messages []string
		expected string
	}{
		// Fix
		{
			name: "Fix",
			messages: []string{
				"fix: fix bug",
				"Updated version (0.0.1) in tmp",
			}, expected: "0.0.1",
		},
		// Feat
		{
			name: "Feat",
			messages: []string{
				"fix: fix bug",
				"feat: add new feature",
				"Updated version (0.0.1) in tmp",
			}, expected: "0.1.0",
		},
		// BC
		{
			name: "BC",
			messages: []string{
				"fix: fix bug",
				"feat: add new feature",
				"bc: breaking change",
				"Updated version (0.0.1) in tmp",
			}, expected: "1.0.0",
		},
		// None
		{
			name: "None",
			messages: []string{
				"fix bug",
				"add new feature",
				"breaking change",
			}, expected: "0.0.0",
		},
	}

	mockGitHandler.EXPECT().GetChangedFiles().Return([]string{
		"file1.txt",
		"file2.txt",
	}).MaxTimes(len(testCases))
	mockGitHandler.EXPECT().GetFromCommit().Return("0.0.0").MaxTimes(len(testCases))
	mockGitHandler.EXPECT().GetDirPath().Return("/tmp").MaxTimes(len(testCases))
	mockGitHandler.EXPECT().RetrieveData().Return(nil).MaxTimes(len(testCases))
	mockGitHandler.EXPECT().SetFromCommit("abcdef1").MaxTimes(len(testCases))
	mockGitHandler.EXPECT().GetLastCommit().Return("abcdef1").AnyTimes()

	version := VersionData{
		filePath: "/tmp/.version.json",
		Commit:   DefaultCommit,
		Version:  DefaultVersionTag,
		Prefix:   "v",

		git: mockGitHandler,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockGitHandler.EXPECT().GetCommitMessages().Return(tc.messages)

			tmpDir := t.TempDir()
			err := version.Initialize(tmpDir)
			assert.NoError(t, err, "Must not have error when initializing")

			relativeFilePath, err := getRelativePath(version.filePath)
			assert.NoError(t, err, "Must not have error when getting the relative path")
			addFiles := append([]string{}, relativeFilePath)
			commitMessage := "Updated version (" + tc.expected + ") in tmp"
			tagMessage := tc.expected + "_tmp"
			mockGitHandler.EXPECT().ConfirmChanges(addFiles, commitMessage, tagMessage).Return(nil).AnyTimes()
			mockGitHandler.EXPECT().GetOutput().Return([]string{
				"[master abcdef1] Commit message example",
				"1 file changed, 1 insertion(+)",
			}).AnyTimes()

			newVersion, err := version.UpdateVersion()
			assert.NoError(t, err, "Must not have error when updating the version")
			assert.Equal(t, tc.expected, newVersion, "The updated version does not match the expected one")
		})
	}
}

func TestVersionData_IncrementVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHandler := mockGit.NewMockGitI(ctrl)

	testCases := []struct {
		name     string
		incType  string
		expected string
	}{
		{
			name:     "Major",
			incType:  "major",
			expected: "1.0.0",
		},
		{
			name:     "Minor",
			incType:  "minor",
			expected: "0.1.0",
		},
		{
			name:     "Patch",
			incType:  "patch",
			expected: "0.0.1",
		},
		{
			name:     "None",
			incType:  "none",
			expected: "0.0.0",
		},
	}

	version := VersionData{
		filePath: "/tmp/.version.json",
		Commit:   DefaultCommit,
		Version:  DefaultVersionTag,
		Prefix:   "v",

		git: mockGitHandler,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			err := version.Initialize(tmpDir)
			assert.NoError(t, err, "Must not have error when initializing")

			relativeFilePath, err := getRelativePath(version.filePath)
			assert.NoError(t, err, "Must not have error when getting the relative path")
			addFiles := append([]string{}, relativeFilePath)

			if tc.incType != "none" {
				mockGitHandler.EXPECT().GetLastCommit().Return("abcdef1")
				mockGitHandler.EXPECT().GetDirPath().Return("/tmp")
				mockGitHandler.EXPECT().ConfirmChanges(addFiles, "Updated version ("+tc.expected+") in tmp", tc.expected+"_tmp").Return(nil)
				mockGitHandler.EXPECT().GetOutput().Return([]string{}).AnyTimes()

				newVersion, err := version.IncrementVersion(tc.incType)
				assert.NoError(t, err, "Must not have error when incrementing the version")
				assert.Equal(t, tc.expected, newVersion, "The incremented version does not match the expected one")
			} else {
				_, err := version.IncrementVersion(tc.incType)
				assert.Error(t, err, "Must have error when incrementing the version with an invalid increment type")
			}
		})
	}
}

func TestVersionData_UpdateChangelog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHandler := mockGit.NewMockGitI(ctrl)

	mockGitHandler.EXPECT().GetCommitMessages().Return([]string{
		"fix: fix bug",
		"feat: add new feature",
		"bc: breaking change",
		"Updated version (0.0.1) in tmp",
	})

	version := VersionData{
		filePath: "/tmp/.version.json",
		Commit:   DefaultCommit,
		Version:  DefaultVersionTag,
		Prefix:   "v",

		git: mockGitHandler,
	}

	tmpDir := t.TempDir()
	err := version.Initialize(tmpDir)
	assert.NoError(t, err, "Must not have error when initializing")

	mockGitHandler.EXPECT().GetDirPath().Return(tmpDir)

	err = version.UpdateChangelog()
	assert.NoError(t, err, "Must not have error when updating the changelog")
	assert.FileExistsf(t, filepath.Join(tmpDir, "CHANGELOG.md"), "The CHANGELOG.md file should exist in the directory "+tmpDir)
}
