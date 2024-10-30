package bumpmanager

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/freepik-company/gommitizen/internal/cmdgit"
)

func IncrementVersion(currentVersionStr string, incType string) (string, string, error) {
	currentVersion, err := semver.NewVersion(currentVersionStr)
	if err != nil {
		return "", "", err
	}

	var newVersionStr string
	var newVersion semver.Version
	if incType == "major" {
		newVersionStr = "Major"
		newVersion = currentVersion.IncMajor() // Increment the major (for example, from 1.2.3 to 2.0.0)
	} else if incType == "minor" {
		newVersionStr = "Minor"
		newVersion = currentVersion.IncMinor() // Increment the minor (for example, from 1.2.3 to 1.3.0)
	} else if incType == "patch" {
		newVersionStr = "Patch"
		newVersion = currentVersion.IncPatch() // Increment the patch (for example, from 1.2.3 to 1.2.4)
	} else {
		newVersionStr = ""
		newVersion = *currentVersion
	}

	return newVersion.String(), newVersionStr, nil
}

func BumpCommitAll(modifiedFiles []string, tagVersions []string) ([]string, error) {
	if len(modifiedFiles) == 0 && len(tagVersions) == 0 {
		return []string{"Nothing to commit"}, nil
	}

	for _, filePath := range modifiedFiles {
		_, err := cmdgit.AddFilePath(filePath)
		if err != nil {
			return nil, fmt.Errorf("error adding file %s: %v", filePath, err)
		}
	}

	message := fmt.Sprintf("bump: new version %s", tagVersions[0])
	if len(tagVersions) > 1 {
		message = fmt.Sprintf("bump: new versions %s", strings.Join(tagVersions, ", "))
	}

	_, err := cmdgit.CreateCommit(message)
	if err != nil {
		return nil, fmt.Errorf("error committing %s: %v", message, err)
	}

	for _, tagVersion := range tagVersions {
		_, err := cmdgit.CreateTag(tagVersion)
		if err != nil {
			return nil, fmt.Errorf("error tagging %s: %v", tagVersion, err)
		}
	}

	return []string{"Files added and committed"}, nil
}
