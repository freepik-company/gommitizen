package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gommitizen/git"
)

// Constants
const ConfigFileName = ".version.json"
const DefaultCommit = "HEAD^"
const DefaultVersionTag = "0.0.0"

// Variables
var bcPrefix = []string{"BREAKING CHANGE:", "breaking change:", "Breaking change:", "bc:", "BC:", "Bc:"}
var featPrefix = []string{"feat:", "Feat:", "feature:", "Feature:", "FEAT"}
var fixPrefix = []string{"fix:", "Fix:", "FIX", "bug:", "Bug:", "BUG"}

// VersionData Manage the version information for our project
type VersionData struct {
	Version         string   `json:"version"`
	Commit          string   `json:"commit"`
	VersionFiles    []string `json:"version_files"`
	Prefix          string   `json:"prefix"`
	filePath        string
	git             git.GitI
	updateChangelog bool
}

// Public Functions

// NewVersionData Create a new VersionData object
func NewVersionData(version string, commit string, filePath string, prefix string) *VersionData {
	var err error

	if prefix == "" {
		prefix = filepath.Base(filepath.Dir(filePath))
	}

	// Get the relative path to the current directory
	var relativePath string
	relativePath, err = getRelativePath(filePath)
	if err != nil {
		panic("Error obtaining the relative path: " + err.Error())
	}

	// Get the base path of the file
	dirPath := filepath.Dir(relativePath)

	// New Git object
	git := git.NewGit(dirPath, commit)

	thisVersion := &VersionData{
		Version:  version,
		Commit:   commit,
		filePath: filePath,
		Prefix:   prefix,
		git:      git,
	}

	return thisVersion
}

// LoadVersionData Load the version data from the .version.json file
func LoadVersionData(filePath string) *VersionData {
	_, err := os.Stat(filePath)

	if err != nil {
		panic("[WARNING] Error when reading .version.json: " + err.Error())
	}

	content, err := os.ReadFile(filePath)

	if err != nil {
		panic("[WARNING] Error when reading .version.json: " + err.Error())
	}

	version := &VersionData{}
	err = json.Unmarshal(content, version)

	if err != nil {
		panic("[WARNING] Error loading .version.json: " + err.Error())
	}

	// Get the relative path to the current directory
	var relativePath string
	relativePath, err = getRelativePath(filePath)
	if err != nil {
		panic("Error obtaining the relative path: " + err.Error())
	}

	// Get the base path of the file
	dirPath := filepath.Dir(relativePath)

	// New Git object
	git := git.NewGit(dirPath, version.Commit)

	version.SetGit(git)

	version.filePath = filePath

	if err != nil {
		panic("[WARNING] Error loading .version.json: " + err.Error())
	}

	return version
}

// EmptyVersionData Create a new empty VersionData object
func EmptyVersionData(filePath string) *VersionData {
	newVersion := NewVersionData("", "", filePath, "")
	err := newVersion.Save()

	if err != nil {
		panic("[WARNING] Error when creating .version.json: " + err.Error())
	}

	return newVersion
}

func (version *VersionData) Initialize(path string) error {
	// check .version.json does not exist
	configFile := filepath.Join(path, ConfigFileName)
	if _, err := os.Stat(configFile); err == nil {
		fmt.Println("Repository already initialized")
		os.Exit(1)
	}
	version.Commit = DefaultCommit
	version.Version = DefaultVersionTag
	version.filePath = configFile

	err := version.Save()
	if err != nil {
		return fmt.Errorf("Error saving .config.json: %s", err)
	}
	return nil
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

func (version *VersionData) GetPrefix() string {
	return version.Prefix
}

func (version *VersionData) GetGit() git.GitI {
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

func (version *VersionData) SetGit(g git.GitI) {
	version.git = g
}

func (version *VersionData) SetUpdateChangelog(uc bool) {
	version.updateChangelog = uc
}
