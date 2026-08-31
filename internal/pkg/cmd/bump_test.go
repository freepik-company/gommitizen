package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirmBumpOnBranch(t *testing.T) {
	tests := []struct {
		name               string
		branch             string
		assumeYes          bool
		input              string
		want               bool
		wantErr            string
		wantOutputContains string
		wantPrompt         bool
	}{
		{name: "main continues without prompting", branch: "main", want: true},
		{name: "master continues without prompting", branch: "master", want: true},
		{name: "feature branch accepts y", branch: "feature/example", input: "y\n", want: true, wantOutputContains: "feature/example", wantPrompt: true},
		{name: "feature branch accepts y at eof", branch: "feature/example", input: "y", want: true, wantOutputContains: "feature/example", wantPrompt: true},
		{name: "feature branch accepts yes case insensitively", branch: "feature/example", input: "YES\n", want: true, wantOutputContains: "feature/example", wantPrompt: true},
		{name: "feature branch rejects n", branch: "feature/example", input: "n\n", wantOutputContains: "feature/example", wantPrompt: true},
		{name: "feature branch defaults to no", branch: "feature/example", input: "\n", wantOutputContains: "feature/example", wantPrompt: true},
		{name: "feature branch accepts yes flag", branch: "feature/example", assumeYes: true, want: true, wantOutputContains: "feature/example"},
		{name: "feature branch reports eof", branch: "feature/example", wantErr: "rerun with --yes", wantOutputContains: "feature/example", wantPrompt: true},
		{name: "detached HEAD prompts safely", input: "n\n", wantOutputContains: "detached HEAD", wantPrompt: true},
		{name: "detached HEAD is identified", assumeYes: true, want: true, wantOutputContains: "detached HEAD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			got, err := confirmBumpOnBranch(tt.branch, tt.assumeYes, strings.NewReader(tt.input), &output)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("confirmBumpOnBranch() unexpected error = %v", err)
			}
			if tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Fatalf("confirmBumpOnBranch() error = %v, want error containing %q", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("confirmBumpOnBranch() = %v, want %v", got, tt.want)
			}

			prompted := strings.Contains(output.String(), "Continue with the bump?")
			if prompted != tt.wantPrompt {
				t.Errorf("confirmBumpOnBranch() prompted = %v, want %v", prompted, tt.wantPrompt)
			}
			if !strings.Contains(output.String(), tt.wantOutputContains) {
				t.Errorf("output %q does not contain %q", output.String(), tt.wantOutputContains)
			}
		})
	}
}
