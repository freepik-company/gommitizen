package git

// TODO: Improve the error handling in this package creating custom errors for each case

// Managed Custom Errors
type GitError struct {
	Message string
}

func (e *GitError) Error() string {
	return e.Message
}
