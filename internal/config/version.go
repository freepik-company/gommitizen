package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigVersion struct {
	path string

	Version      string   `json:"version"`
	Commit       string   `json:"commit"`
	VersionFiles []string `json:"version_files"`
	Prefix       string   `json:"prefix"`
}

func NewConfigVersion(path string, version string, commit string, prefix string) *ConfigVersion {
	nPath, err := normalizePath(path)
	if err != nil {
		panic(fmt.Errorf("normalizePath %s: %v", path, err))
	}

	nPrefix := prefix
	if len(prefix) == 0 {
		nPrefix = filepath.Base(path)
	}

	return &ConfigVersion{
		path: nPath,

		Version:      version,
		Commit:       commit,
		VersionFiles: make([]string, 0),
		Prefix:       nPrefix,
	}
}

func Read(path string) (*ConfigVersion, error) {
	fileVersionPath := filepath.Join(path, defaultFileName)
	data, err := os.ReadFile(fileVersionPath)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %v", fileVersionPath, err)
	}

	var version ConfigVersion
	err = json.Unmarshal(data, &version)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %v", err)
	}

	version.path = path
	return &version, nil
}

func (v ConfigVersion) Save() error {
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

func (v ConfigVersion) GetFileVersionPath() string {
	return filepath.Join(v.path, defaultFileName)
}
