package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewConfigVersion(t *testing.T) {
	dirPath := "/tmp"
	version := "1.0.0"
	commit := "abc123"
	alias := "v"

	v := NewConfigVersion(dirPath, version, commit, alias)

	if v.dirPath != dirPath {
		t.Errorf("expected path %s, got %s", dirPath, v.dirPath)
	}
	if v.Version != version {
		t.Errorf("expected version %s, got %s", version, v.Version)
	}
	if v.Commit != commit {
		t.Errorf("expected commit %s, got %s", commit, v.Commit)
	}
	if v.Alias != alias {
		t.Errorf("expected prefix %s, got %s", alias, v.Alias)
	}
	if len(v.VersionFiles) != 0 {
		t.Errorf("expected empty VersionFiles, got %v", v.VersionFiles)
	}
}

func TestRead(t *testing.T) {
	// Crear un archivo JSON temporal
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, defaultFileName)
	versionData := ConfigVersion{
		Version:      "1.0.0",
		Commit:       "abc123",
		VersionFiles: []string{"file1", "file2"},
		Alias:        "v",
	}
	data, err := json.Marshal(versionData)
	if err != nil {
		t.Fatalf("failed to marshal version data: %v", err)
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Leer el archivo JSON usando la función Read
	v, err := ReadConfigVersion(tempDir)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	// Verificar los resultados
	expected := &ConfigVersion{
		dirPath:      tempDir,
		Version:      "1.0.0",
		Commit:       "abc123",
		VersionFiles: []string{"file1", "file2"},
		Alias:        "v",
	}
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("Read() = %v, want %v", v, expected)
	}
}
