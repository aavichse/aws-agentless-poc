package infra

// Helper function to return a pointer to a string
func StrPtr(s string) *string {
	return &s
}

// Helper function to return a pointer to a slice of strings
func SlicePtr(s []string) *[]string {
	return &s
}
