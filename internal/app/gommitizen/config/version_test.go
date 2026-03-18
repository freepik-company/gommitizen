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

func TestUpdateVersionOfFiles(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		substring   string
		newVersion  string
		wantContent string
		wantModified bool
	}{
		{
			name:        "unquoted YAML value",
			content:     "image: myapp\n  tag: 0.8.2\nreplicas: 3\n",
			substring:   "tag",
			newVersion:  "0.9.0",
			wantContent: "image: myapp\n  tag: 0.9.0\nreplicas: 3\n",
			wantModified: true,
		},
		{
			name:        "double-quoted YAML value",
			content:     "image: myapp\n  tag: \"0.8.2\"\nreplicas: 3\n",
			substring:   "tag",
			newVersion:  "0.9.0",
			wantContent: "image: myapp\n  tag: \"0.9.0\"\nreplicas: 3\n",
			wantModified: true,
		},
		{
			name:        "single-quoted YAML value",
			content:     "image: myapp\n  tag: '0.8.2'\nreplicas: 3\n",
			substring:   "tag",
			newVersion:  "0.9.0",
			wantContent: "image: myapp\n  tag: '0.9.0'\nreplicas: 3\n",
			wantModified: true,
		},
		{
			name:        "no match returns false",
			content:     "image: myapp\nreplicas: 3\n",
			substring:   "tag",
			newVersion:  "0.9.0",
			wantContent: "image: myapp\nreplicas: 3\n",
			wantModified: false,
		},
		{
			name:        "same version returns false",
			content:     "  tag: 0.9.0\n",
			substring:   "tag",
			newVersion:  "0.9.0",
			wantContent: "  tag: 0.9.0\n",
			wantModified: false,
		},
		{
			name:        "equals sign separator",
			content:     "version = 1.2.3\n",
			substring:   "version",
			newVersion:  "1.3.0",
			wantContent: "version = 1.3.0\n",
			wantModified: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), "testfile.yaml")
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}

			modified, err := updateVersionOfFiles(tmpFile, tt.substring, tt.newVersion)
			if err != nil {
				t.Fatalf("updateVersionOfFiles() error = %v", err)
			}
			if modified != tt.wantModified {
				t.Errorf("updateVersionOfFiles() modified = %v, want %v", modified, tt.wantModified)
			}

			got, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("failed to read result file: %v", err)
			}
			if string(got) != tt.wantContent {
				t.Errorf("file content = %q, want %q", string(got), tt.wantContent)
			}
		})
	}
}

func TestUpdateVersionWithVersionFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create a YAML file with a quoted version
	yamlContent := "image: myapp\n  tag: \"1.0.0\"\n"
	yamlPath := filepath.Join(tempDir, "values.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	// Create a ConfigVersion pointing to the YAML file
	cv := &ConfigVersion{
		dirPath:      tempDir,
		Version:      "1.0.0",
		Commit:       "abc123",
		VersionFiles: []string{"values.yaml:tag"},
		Alias:        "test",
	}
	// Save the .version.json so UpdateVersion can re-save it
	data, err := json.MarshalIndent(cv, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, defaultFileName), data, 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	modifiedFiles, err := cv.UpdateVersion("2.0.0", "def456")
	if err != nil {
		t.Fatalf("UpdateVersion() error = %v", err)
	}

	// Should include both .version.json and values.yaml
	if len(modifiedFiles) != 2 {
		t.Fatalf("expected 2 modified files, got %d: %v", len(modifiedFiles), modifiedFiles)
	}

	// Verify the YAML file was actually updated
	got, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("failed to read yaml file: %v", err)
	}
	expected := "image: myapp\n  tag: \"2.0.0\"\n"
	if string(got) != expected {
		t.Errorf("yaml content = %q, want %q", string(got), expected)
	}
}

func TestUpdateVersionNoMatchExcludedFromModified(t *testing.T) {
	tempDir := t.TempDir()

	// Create a YAML file WITHOUT a matching key
	yamlContent := "image: myapp\nreplicas: 3\n"
	yamlPath := filepath.Join(tempDir, "values.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	cv := &ConfigVersion{
		dirPath:      tempDir,
		Version:      "1.0.0",
		Commit:       "abc123",
		VersionFiles: []string{"values.yaml:tag"},
		Alias:        "test",
	}
	data, err := json.MarshalIndent(cv, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, defaultFileName), data, 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	modifiedFiles, err := cv.UpdateVersion("2.0.0", "def456")
	if err != nil {
		t.Fatalf("UpdateVersion() error = %v", err)
	}

	// Should only include .version.json, NOT values.yaml
	if len(modifiedFiles) != 1 {
		t.Fatalf("expected 1 modified file, got %d: %v", len(modifiedFiles), modifiedFiles)
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
