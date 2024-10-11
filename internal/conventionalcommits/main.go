package conventionalcommits

import (
	"fmt"
	"gommitizen/internal/cmdgit"
	"regexp"
)

// https://www.conventionalcommits.org/en/v1.0.0/

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

type ConventionalCommit struct {
	ShortCommit string
	Date        string

	Type        string
	Scope       string
	Description string
}

func (cc ConventionalCommit) String() string {
	if cc.Scope == "" {
		return fmt.Sprintf("%s: %s #%s", cc.Type, cc.Description, cc.ShortCommit)
	} else {
		return fmt.Sprintf("%s(%s): %s #%s", cc.Type, cc.Scope, cc.Description, cc.ShortCommit)
	}
}

func FilterAndParse(commits []cmdgit.Commit) []ConventionalCommit {
	conventionalcommits := make([]ConventionalCommit, 0)
	for _, commit := range commits {

		re := regexp.MustCompile(`(?P<type>\w+)(\((?P<scope>\w+)\))?:(?P<message>.+)`)
		match := re.FindStringSubmatch(commit.Title)
		if match == nil {
			continue
		}

		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		cc := ConventionalCommit{
			ShortCommit: commit.ShortCommit,
			Date:        commit.Date,

			Type:        result["type"],
			Scope:       result["scope"],
			Description: result["message"],
		}
		conventionalcommits = append(conventionalcommits, cc)
	}
	return conventionalcommits
}

func DetermineIncrementType(conventionalCommits []ConventionalCommit) string {
	major := false
	minor := false
	patch := false

	for _, cc := range conventionalCommits {
		for _, prefix := range bcPrefix {
			if cc.Type == prefix {
				major = true
				break
			}
		}
		for _, prefix := range featPrefix {
			if cc.Type == prefix {
				minor = true
				break
			}
		}
		for _, prefix := range fixPrefix {
			if cc.Type == prefix {
				patch = true
				break
			}
		}
		for _, prefix := range refactorPrefix {
			if cc.Type == prefix {
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
