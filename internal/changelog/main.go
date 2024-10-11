package changelog

import (
	"gommitizen/internal/conventionalcommits"
	"path/filepath"
)

type data struct {
	version string
	commits []conventionalcommits.ConventionalCommit
}

const changelogFileName = "CHANGELOG.md"

func Apply(dirPath string, version string, cvCommits []conventionalcommits.ConventionalCommit) (string, error) {
	changelogFilePath := filepath.Join(dirPath, changelogFileName)

	// TODO

	return changelogFilePath, nil
}
