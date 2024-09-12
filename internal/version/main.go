package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Version struct {
	path string

	Version      string   `json:"version"`
	Commit       string   `json:"commit"`
	VersionFiles []string `json:"version_files"`
	Prefix       string   `json:"prefix"`
}

const (
	defaultFileName = ".version.json"
)

func New(path string, version string, commit string, prefix string) *Version {
	return &Version{
		path: path,

		Version:      version,
		Commit:       commit,
		VersionFiles: make([]string, 0),
		Prefix:       prefix,
	}
}

func Read(path string) (*Version, error) {
	fileVersionPath := filepath.Join(path, defaultFileName)
	data, err := os.ReadFile(fileVersionPath)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %v", fileVersionPath, err)
	}

	var version Version
	err = json.Unmarshal(data, &version)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %v", err)
	}

	version.path = path
	return &version, nil
}

func (v Version) Save() error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("parse struct to json: %v", err)
	}

	err = os.WriteFile(v.GetFileVersionPath(), data, 0644)
	if err != nil {
		return fmt.Errorf("write file %s: %v", v.GetFileVersionPath(), err)
	}

	return nil
}

func (v Version) GetFileVersionPath() string {
	return filepath.Join(v.path, defaultFileName)
}
