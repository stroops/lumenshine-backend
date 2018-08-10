package helpers

import (
	"strings"
)

// StringInSliceI checks if a string is inside a slice (caseinsensitive)
func StringInSliceI(a string, list []string) bool {
	a = strings.ToLower(a)
	for _, b := range list {
		if strings.ToLower(b) == a {
			return true
		}
	}
	return false
}

// StringInSlice checks if a string is inside a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IntInSlice checks if a int is inside a slice
func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
