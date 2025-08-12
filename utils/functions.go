// utils/functions.go
package utils

import (
	"strconv"
	"strings"
)

// IsInList checks if an item exists in a list of strings
func IsInList(items []string, item string) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

// ParseIntList splits a string into a slice of ints.
// Returns nil if the input is empty or no valid ints found.
func ParseIntList(s string) []int {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []int
	for _, p := range parts {
		if v, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			out = append(out, v)
		}
	}
	return out
}

// ParseStringList splits a string into a slice of strings.
// Returns nil if the input is empty.
func ParseStringList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}
