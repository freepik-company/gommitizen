package bumpmanager

import (
	"fmt"
	"gommitizen/internal/cmdgit"
	"strings"

	"github.com/Masterminds/semver"
)

var bcPrefix = []string{
	"BREAKING CHANGE:", "BREAKING CHANGE(",
	"breaking change:", "breaking change(",
	"Breaking change:", "Breaking change(",
	"bc:", "bc(",
	"BC:", "BC(",
	"Bc:", "Bc(",
}
var featPrefix = []string{
	"feat:", "feat(",
	"Feat:", "Feat(",
	"feature:", "feature(",
	"Feature:", "Feature(",
	"FEAT:", "FEAT(",
}
var fixPrefix = []string{
	"fix:", "fix(",
	"Fix:", "Fix(",
	"FIX:", "FIX(",
	"bug:", "bug(",
	"Bug:", "Bug(",
	"BUG:", "BUG(",
	"bugfix:", "bugfix(",
	"Bugfix:", "Bugfix(",
	"BUGFIX:", "BUGFIX(",
}
var refactorPrefix = []string{
	"refactor:", "refactor(",
	"Refactor:", "Refactor(",
	"REFACTOR:", "REFACTOR(",
}

func DetermineVersionBump(commitMessages []string) string {
	major := false
	minor := false
	patch := false

	for _, message := range commitMessages {
		ignoreStartLen := 10
		if strings.HasPrefix(message[ignoreStartLen:], "Updated version") {
			continue
		}
		for _, prefix := range bcPrefix {
			if strings.HasPrefix(message[ignoreStartLen:], prefix) {
				major = true
				break
			}
		}
		for _, prefix := range featPrefix {
			if strings.HasPrefix(message[ignoreStartLen:], prefix) {
				minor = true
				break
			}
		}
		for _, prefix := range fixPrefix {
			if strings.HasPrefix(message[ignoreStartLen:], prefix) {
				patch = true
				break
			}
		}
		for _, prefix := range refactorPrefix {
			if strings.HasPrefix(message[ignoreStartLen:], prefix) {
				patch = true
				break
			}
		}
	}

	if major {
		return "major"
	} else if minor {
		return "minor"
	} else if patch {
		return "patch"
	}

	return "none"
}

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
		return []string{
			"Nothing to commit",
		}, nil
	}

	message := fmt.Sprintf("Updated version (%s)", tagVersions[0])
	if len(tagVersions) > 1 {
		message = fmt.Sprintf("Updated version (%s)", strings.Join(tagVersions, ", "))
	}

	for _, filePath := range modifiedFiles {
		_, err := cmdgit.Add(filePath)
		if err != nil {
			return nil, fmt.Errorf("error adding file %s: %v", filePath, err)
		}
	}

	_, err := cmdgit.Commit(message)
	if err != nil {
		return nil, fmt.Errorf("error committing %s: %v", message, err)
	}

	for _, tagVersion := range tagVersions {
		_, err := cmdgit.Tag(tagVersion)
		if err != nil {
			return nil, fmt.Errorf("error tagging %s: %v", tagVersion, err)
		}
	}

	return []string{
		"Files added and committed",
	}, nil
}
