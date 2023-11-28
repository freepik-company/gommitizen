package version

// VersionError Custom error type
type VersionError struct {
	Message string
}

func (e *VersionError) Error() string {
	return e.Message
}
