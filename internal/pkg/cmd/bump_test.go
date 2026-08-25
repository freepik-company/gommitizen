package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirmBumpOnBranch(t *testing.T) {
	tests := []struct {
		name       string
		branch     string
		input      string
		want       bool
		wantPrompt bool
	}{
		{name: "main continues without prompting", branch: "main", want: true},
		{name: "master continues without prompting", branch: "master", want: true},
		{name: "feature branch accepts y", branch: "feature/example", input: "y\n", want: true, wantPrompt: true},
		{name: "feature branch accepts yes case insensitively", branch: "feature/example", input: "YES\n", want: true, wantPrompt: true},
		{name: "feature branch rejects n", branch: "feature/example", input: "n\n", wantPrompt: true},
		{name: "feature branch defaults to no", branch: "feature/example", input: "\n", wantPrompt: true},
		{name: "feature branch defaults to no on eof", branch: "feature/example", wantPrompt: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			got, err := confirmBumpOnBranch(tt.branch, strings.NewReader(tt.input), &output)
			if err != nil {
				t.Fatalf("confirmBumpOnBranch() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("confirmBumpOnBranch() = %v, want %v", got, tt.want)
			}

			prompted := output.Len() > 0
			if prompted != tt.wantPrompt {
				t.Errorf("confirmBumpOnBranch() prompted = %v, want %v", prompted, tt.wantPrompt)
			}
			if tt.wantPrompt && !strings.Contains(output.String(), tt.branch) {
				t.Errorf("prompt %q does not include branch %q", output.String(), tt.branch)
			}
		})
	}
}
