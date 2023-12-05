package git

// Remove a string from a slice of strings
func RemoveStringFromSlice(slice []string, s string) []string {
	var result []string

	for _, str := range slice {
		if str != s {
			result = append(result, str)
		}
	}

	return result
}
