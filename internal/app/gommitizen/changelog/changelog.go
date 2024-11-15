package changelog

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/freepik-company/gommitizen/internal/app/gommitizen/conventionalcommits"
)

type data struct {
	Version string
	Date    string

	BreakingChanges []conventionalcommits.CommitData
	Features        []conventionalcommits.CommitData
	BugFixes        []conventionalcommits.CommitData
	Refactors       []conventionalcommits.CommitData
	Miscellaneous   []conventionalcommits.CommitData
}

const changelogFileName = "CHANGELOG.md"

//go:embed template.tpl
var tplFile embed.FS

func Apply(dirPath string, version string, commits []conventionalcommits.CommitData) (string, error) {
	changelogFilePath := filepath.Join(dirPath, changelogFileName)
	groupByCommonChangeType := groupByCommonChangeType(commits)

	data := data{
		Version: version,
		Date:    time.Now().Format("2006-01-02"),

		BreakingChanges: groupByCommonChangeType[conventionalcommits.CommonNameBC],
		Features:        groupByCommonChangeType[conventionalcommits.CommonNameFeat],
		BugFixes:        groupByCommonChangeType[conventionalcommits.CommonNameFix],
		Refactors:       groupByCommonChangeType[conventionalcommits.CommonNameRefactor],
		Miscellaneous:   groupByCommonChangeType[conventionalcommits.CommonNameMiscellaneous],
	}

	tpl, err := template.ParseFS(tplFile, "template.tpl")
	if err != nil {
		return "", fmt.Errorf("fail to load template: %v", err)
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("fail to execute template: %v", err)
	}

	err = prependToFile(changelogFilePath, buf)
	if err != nil {
		return "", fmt.Errorf("fail to prepend to file: %v", err)
	}

	return changelogFilePath, nil
}

func prependToFile(changelogFilePath string, data bytes.Buffer) error {
	existingContent, err := os.ReadFile(changelogFilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read file: %v", err)
	}

	outFile, err := os.OpenFile(changelogFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %v", err)
	}
	defer outFile.Close()

	_, err = outFile.Write(append(data.Bytes(), existingContent...))
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func groupByCommonChangeType(commits []conventionalcommits.CommitData) map[string][]conventionalcommits.CommitData {
	groups := make(map[string][]conventionalcommits.CommitData)

	for _, item := range commits {
		changeType := item.CommonChangeType
		groups[changeType] = append(groups[changeType], item)
	}

	return groups
}
