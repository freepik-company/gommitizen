package git

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func GetFirstCommit() (string, error) {
	cmd := "git rev-list --max-parents=0 HEAD"
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func AddFilePath(filePath string) (string, error) {
	cmd := fmt.Sprintf("git add %s", filePath)
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func CreateTag(tag string) (string, error) {
	cmd := fmt.Sprintf("git tag %s", tag)
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func CreateCommit(message string) (string, error) {
	cmd := fmt.Sprintf(`git commit -m "%s"`, message)
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func GetLastCommit() (string, error) {
	cmd := "git rev-parse HEAD"
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// https://git-scm.com/docs/pretty-formats
func GetCommits(fromCommit string, fromPath string) ([]Commit, error) {
	pretty := `--pretty=format:'{"hash": "%H", "date": "%ad", "subject": "%s"}'`
	dateFormat := `--date=format-local:'%Y-%m-%dT%H:%M:%SZ'`
	cmd := fmt.Sprintf(`git log %s %s %s.. -- %s`, pretty, dateFormat, fromCommit, fromPath)
	slog.Debug(fmt.Sprintf("exec: %s", cmd))

	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return []Commit{}, fmt.Errorf("fail %s: %v", cmd, err)
	}

	commits := make([]Commit, 0)
	for _, line := range strings.Split(string(output), "\n") {
		if len(line) > 0 {

			var commit Commit
			err = json.Unmarshal([]byte(line), &commit)
			if err != nil {
				return []Commit{}, fmt.Errorf("fail unmarshal json %s: %v", cmd, err)
			}

			slog.Debug(fmt.Sprintf("commit: %v", commit))
			commits = append(commits, commit)
		}
	}
	return commits, nil
}
