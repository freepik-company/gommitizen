package conventionalcommits

// TODO: simplify this package using the new parser structure of package go-conventionalcommits

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"

	"github.com/freepik-company/gommitizen/internal/git"
)

// https://www.conventionalcommits.org/en/v1.0.0/

type ChangeType struct {
	Order      int
	CommonName string
	Prefixes   []string
}

type CommitData struct {
	ShortHash string
	Hash      string
	Date      time.Time

	CommonChangeType string
	ChangeType       string
	Scope            string
	Subject          string
}

const (
	CommonNameBC            = "Breaking changes"
	CommonNameFeat          = "Features"
	CommonNameFix           = "Bug Fixes"
	CommonNameRefactor      = "Refactors"
	CommonNameMiscellaneous = "Miscellaneous"
)

var (
	changeTypes = []ChangeType{
		{
			CommonName: CommonNameBC,
			Prefixes:   []string{"bc", "breaking change"},
		},
		{
			CommonName: CommonNameFeat,
			Prefixes:   []string{"feat", "feature"},
		},
		{
			CommonName: CommonNameFix,
			Prefixes:   []string{"fix", "bug", "bugfix"},
		},
		{
			CommonName: CommonNameRefactor,
			Prefixes:   []string{"refactor"},
		},
		{
			CommonName: CommonNameMiscellaneous,
			Prefixes:   []string{"perf", "performance", "test", "tests", "chore", "ci", "build", "docs", "style"},
		},
	}
)

func (cc CommitData) String() string {
	if cc.Scope == "" {
		return fmt.Sprintf("%s: %s #%s", cc.ChangeType, cc.Subject, cc.ShortHash)
	} else {
		return fmt.Sprintf("%s(%s): %s #%s", cc.ChangeType, cc.Scope, cc.Subject, cc.ShortHash)
	}
}

func ReadConventionalCommits(commits []git.Commit) []CommitData {
	cvcommits := make([]CommitData, 0)

	opts := []conventionalcommits.MachineOption{
		parser.WithTypes(conventionalcommits.TypesConventional),
		parser.WithBestEffort(),
	}

	for _, commit := range commits {
		res, err := parser.NewMachine(opts...).Parse([]byte(commit.Subject))
		if err != nil || !res.Ok() {
			slog.Debug(fmt.Sprintf("ignore commit, no cc by parser: %s", commit.Subject))
			continue
		}
		ccData := res.(*conventionalcommits.ConventionalCommit)

		commonChangeType := determinateCommonChangeType(ccData.Type)
		if commonChangeType == "unknown" {
			slog.Debug(fmt.Sprintf("ignore commit, no cc by common: %s", ccData.Type))
			continue
		}

		scope := ""
		if ccData.Scope != nil {
			scope = *ccData.Scope
		}

		cc := CommitData{
			ShortHash: commit.AbbreviationHash(),
			Hash:      commit.Hash,
			Date:      commit.Date,

			CommonChangeType: commonChangeType,
			ChangeType:       ccData.Type,
			Scope:            scope,
			Subject:          ccData.Description,
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

func DetermineIncrementType(commits []CommitData) string {
	var hasMinor, hasPatch bool

	for _, commit := range commits {
		switch commit.CommonChangeType {
		case CommonNameBC:
			return "major"
		case CommonNameFeat:
			hasMinor = true
		case CommonNameFix, CommonNameRefactor:
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
