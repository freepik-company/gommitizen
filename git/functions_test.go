package git

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGit_RemoveStringFromSlice(t *testing.T) {
	t.Run("RemoveStringFromSlice", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		s := "b"

		expected := []string{"a", "c"}
		result := RemoveStringFromSlice(slice, s)
		assert.Equal(t, result, expected, "The slice should be equal")
	})
}
