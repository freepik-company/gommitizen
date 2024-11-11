package git

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
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
	pretty := `--pretty=format:'%H||%ad||%s'`
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
			hash, date, subject := strings.Split(line, "||")[0], strings.Split(line, "||")[1], strings.Split(line, "||")[2]

			dateTime, err := time.Parse("2006-01-02T15:04:05Z", date)
			if err != nil {
				return []Commit{}, fmt.Errorf("fail parsing date %s: %v", date, err)
			}
			commit := Commit{Hash: hash, Date: dateTime, Subject: subject}

			slog.Debug(fmt.Sprintf("commit: %v", commit))
			commits = append(commits, commit)
		}
	}
	return commits, nil
}
