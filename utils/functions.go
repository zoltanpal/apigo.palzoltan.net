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

type TSMode string

const (
	TSOr  TSMode = "OR"
	TSAnd TSMode = "AND"
)

// BuildToTSQuery builds a to_tsquery string like "w1 & w2" (AND) or "w1 | w2" (OR).
// If prefix==true, appends :* to each token for right-truncation (prefix match).
// Returns ok=false if nothing usable remains after trimming.
func BuildToTSQuery(words []string, mode TSMode, prefix bool) (q string, ok bool) {
	toks := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		// light sanitation: strip tsquery operator chars to avoid syntax errors
		w = strings.Map(func(r rune) rune {
			switch r {
			case '&', '|', '!', '(', ')', ':':
				return ' '
			default:
				return r
			}
		}, w)
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		if prefix {
			w += ":*"
		}
		toks = append(toks, w)
	}
	if len(toks) == 0 {
		return "", false
	}
	sep := " | "
	if mode == TSAnd {
		sep = " & "
	}
	return strings.Join(toks, sep), true
}

// SanitizeTSWord cleans a single string for to_tsquery
func SanitizeTSWord(w string) string {
	w = strings.TrimSpace(w)
	w = strings.Map(func(r rune) rune {
		switch r {
		case '&', '|', '!', '(', ')', ':':
			return ' '
		default:
			return r
		}
	}, w)
	return strings.TrimSpace(w)
}
