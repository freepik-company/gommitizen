package version

import (
	"gommitizen/git"
)

// Variables
var bcPrefix = []string{"BREAKING CHANGE:", "breaking change:", "Breaking change:", "bc:", "BC:", "Bc:"}
var featPrefix = []string{"feat:", "Feat:", "feature:", "Feature:", "FEAT"}
var fixPrefix = []string{"fix:", "Fix:", "FIX", "bug:", "Bug:", "BUG"}

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

// Constructor
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

// Setters
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
