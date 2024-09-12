package cmdinit

import (
	"fmt"
	"gommitizen/git"
	"gommitizen/internal/version"
	"log/slog"
	"os"
	"path/filepath"
)

func Run(directory, prefix string) {
	path, err := normalizePath(directory)
	if err != nil {
		panic(err)
	}

	if len(prefix) == 0 {
		prefix = filepath.Base(path)
	}

	slog.Debug(fmt.Sprintf("init gommitizen in %s with %s prefix", path, prefix))

	commit, err := git.GetFirstCommit(path)
	if err != nil {
		panic(err)
	}

	config := version.New(path, "0.0.0", commit, prefix)
	err = config.Save()
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Initializing gommitizen in %s", config.GetFileVersionPath()))
}

func normalizePath(directory string) (string, error) {
	if len(directory) > 0 {
		if isRelativePath(directory) {
			return toAbsolutePath(directory)
		} else {
			return directory, nil
		}
	}
	return getCurrentDirectory()
}

func getCurrentDirectory() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}
	return path, nil
}

func toAbsolutePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", fmt.Errorf("error converting to absolute path: %v", err)
	}
	return absPath, nil
}

func isRelativePath(path string) bool {
	return !filepath.IsAbs(path)
}
