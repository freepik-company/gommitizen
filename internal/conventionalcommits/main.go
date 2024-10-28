package conventionalcommits

import (
	"fmt"
	"gommitizen/internal/cmdgit"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

// https://www.conventionalcommits.org/en/v1.0.0/

type ChangeType struct {
	Order      int
	CommonName string
	Prefixes   []string
}

type ConventionalCommit struct {
	ShortHash string
	Hash      string
	Date      time.Time

	CommonChangeType string
	ChangeType       string
	Scope            string
	Subject          string
}

var (
	commonNameBC       = "Breaking changes"
	commonNameFeat     = "Features"
	commonNameFix      = "Bug Fixes"
	commonNameRefactor = "Miscellaneous"

	changeTypes = []ChangeType{
		{
			CommonName: commonNameBC,
			Prefixes:   []string{"bc", "breaking change"},
		},
		{
			CommonName: commonNameFeat,
			Prefixes:   []string{"feat", "feature"},
		},
		{
			CommonName: commonNameFix,
			Prefixes:   []string{"fix", "bug", "bugfix"},
		},
		{
			CommonName: commonNameRefactor,
			Prefixes:   []string{"refactor", "perf", "performance", "test", "tests", "chore", "ci", "build", "docs", "style"},
		},
	}
)

func (cc ConventionalCommit) String() string {
	if cc.Scope == "" {
		return fmt.Sprintf("%s: %s #%s", cc.ChangeType, cc.Subject, cc.ShortHash)
	} else {
		return fmt.Sprintf("%s(%s): %s #%s", cc.ChangeType, cc.Scope, cc.Subject, cc.ShortHash)
	}
}

func ReadConventionalCommits(commits []cmdgit.Commit) []ConventionalCommit {
	cvcommits := make([]ConventionalCommit, 0)
	for _, commit := range commits {
		re := regexp.MustCompile(`(?P<change_type>\w+)(\((?P<scope>\w+)\))?:\s*(?P<subject>.+?)\s*$`)
		match := re.FindStringSubmatch(commit.Subject)
		if match == nil {
			slog.Debug(fmt.Sprintf("ignore commit, no cc by pattern: %s", commit.Subject))
			continue
		}

		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		commonChangeType := determinateCommonChangeType(result["change_type"])
		if commonChangeType == "unknown" {
			slog.Debug(fmt.Sprintf("ignore commit, no cc by common: %s", result["change_type"]))
			continue
		}

		cc := ConventionalCommit{
			ShortHash: commit.AbbreviationHash(),
			Hash:      string(commit.Hash),
			Date:      time.Time(commit.Date),

			CommonChangeType: commonChangeType,
			ChangeType:       result["change_type"],
			Scope:            result["scope"],
			Subject:          result["subject"],
		}

		slog.Debug(fmt.Sprintf("cccommit: %v", cc))
		cvcommits = append(cvcommits, cc)
	}
	return cvcommits
}

func determinateCommonChangeType(changeType string) string {
	for _, ct := range changeTypes {
		for _, prefix := range ct.Prefixes {
			if strings.ToLower(changeType) == prefix {
				return ct.CommonName
			}
		}
	}
	return "unknown"
}

func DetermineIncrementType(commits []ConventionalCommit) string {
	var hasMinor, hasPatch bool

	for _, commit := range commits {
		switch commit.CommonChangeType {
		case commonNameBC:
			return "major"
		case commonNameFeat:
			hasMinor = true
		case commonNameFix, commonNameRefactor:
			hasPatch = true
		}
	}

	switch {
	case hasMinor:
		return "minor"
	case hasPatch:
		return "patch"
	default:
		return "none"
	}
}
