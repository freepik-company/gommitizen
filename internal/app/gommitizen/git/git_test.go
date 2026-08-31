package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetCurrentBranch(t *testing.T) {
	repoPath := t.TempDir()
	runGitForBranchTest(t, repoPath, "init", "-b", "feature/example")

	trackedFile := filepath.Join(repoPath, "tracked.txt")
	if err := os.WriteFile(trackedFile, []byte("test\n"), 0o600); err != nil {
		t.Fatalf("write tracked file: %v", err)
	}
	runGitForBranchTest(t, repoPath, "add", "tracked.txt")
	runGitForBranchTest(
		t,
		repoPath,
		"-c", "user.name=Gommitizen Test",
		"-c", "user.email=gommitizen@example.com",
		"-c", "commit.gpgsign=false",
		"commit", "-m", "initial commit",
	)

	branch, err := GetCurrentBranch(repoPath)
	if err != nil {
		t.Fatalf("GetCurrentBranch() unexpected error = %v", err)
	}
	if branch != "feature/example" {
		t.Fatalf("GetCurrentBranch() = %q, want %q", branch, "feature/example")
	}

	runGitForBranchTest(t, repoPath, "checkout", "--detach")
	branch, err = GetCurrentBranch(repoPath)
	if err != nil {
		t.Fatalf("GetCurrentBranch() on detached HEAD unexpected error = %v", err)
	}
	if branch != "" {
		t.Fatalf("GetCurrentBranch() on detached HEAD = %q, want empty branch", branch)
	}
}

func runGitForBranchTest(t *testing.T, repoPath string, args ...string) string {
	t.Helper()

	commandArgs := append([]string{"-C", repoPath}, args...)
	output, err := exec.Command("git", commandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(commandArgs, " "), err, output)
	}

	return strings.TrimSpace(string(output))
}
