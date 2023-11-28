package version

import (
	"encoding/json"
	"fmt"
	"gommitizen/git"
	"os"
	"path/filepath"
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
	git             *git.Git
	updateChangelog bool
}

func NewVersionData(version string, commit string, filePath string, prefix string) *VersionData {
	var err error

	// if file already exists, raise an error
	_, err = os.Stat(filePath)
	if err == nil {
		panic("[WARNING] .version.json already exists")
	}

	if prefix == "" {
		prefix = filepath.Base(filepath.Dir(filePath))
	}

	thisVersion := &VersionData{Version: version, Commit: commit, filePath: filePath, Prefix: prefix}

	if err != nil {
		panic("[WARNING] Error when creating .version.json: " + err.Error())
	}

	return thisVersion
}

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

	err = version.load()

	if err != nil {
		panic("[WARNING] Error loading .version.json: " + err.Error())
	}

	return version
}

func EmptyVersionData(filePath string) *VersionData {
	newVersion := NewVersionData("", "", filePath, "")
	err := newVersion.Save()

	if err != nil {
		panic("[WARNING] Error when creating .version.json: " + err.Error())
	}

	return newVersion
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
func (version *VersionData) GetGit() *git.Git {
	return version.git
}

func (version *VersionData) GetUpdateChangelog() bool {
	return version.updateChangelog
}

// Funciones públicas

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

// private methods
func (version *VersionData) load() error {
	var err error
	version.git, err = version.returnGitObjectWithUpdatedData()

	return err
}
