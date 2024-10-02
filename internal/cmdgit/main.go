package cmdgit

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func GetFirstCommit(path string) (string, error) {
	cmd := "git rev-list --max-parents=0 HEAD"
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func Add(filePath string) (string, error) {
	cmd := fmt.Sprintf("git add %s", filePath)
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func Tag(tag string) (string, error) {
	cmd := fmt.Sprintf("git tag %s", tag)
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func Commit(message string) (string, error) {
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

func GetCommitMessages(fromCommit string, path string) ([]string, error) {
	cmd := fmt.Sprintf(`git log --pretty="%%h - %%s" %s.. -- %s`, fromCommit, path)
	slog.Debug(fmt.Sprintf("exec: %s", cmd))
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return []string{}, fmt.Errorf("fail %s: %v", cmd, err)
	}
	commitMessages := make([]string, 0)
	for _, commitMessage := range strings.Split(string(output), "\n") {
		if len(commitMessage) > 0 {
			commitMessages = append(commitMessages, commitMessage)
		}
	}
	return commitMessages, nil
}
