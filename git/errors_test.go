package git

import "testing"

func TestGit_Error(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		e := &GitError{
			Message: "Test error",
		}

		expected := "Test error"
		if e.Error() != expected {
			t.Errorf("Error() = %v, want %v", e.Error(), expected)
		}
	})
}
