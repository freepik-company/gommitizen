package version

import "testing"

func TestVersionData_Error(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		e := &VersionError{
			Message: "Test error",
		}

		expected := "Test error"
		if e.Error() != expected {
			t.Errorf("Error() = %v, want %v", e.Error(), expected)
		}
	})
}
