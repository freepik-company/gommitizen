package git

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
