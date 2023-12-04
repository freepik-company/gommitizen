package version

// TODO: Improve the error handling in this package creating custom errors for each case

// VersionError Custom error type
type VersionError struct {
	Message string
}

func (e *VersionError) Error() string {
	return e.Message
}
