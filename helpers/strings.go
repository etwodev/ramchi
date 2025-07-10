package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"
)

// IsEmpty returns true if the trimmed string is empty.
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Truncate shortens a string to a max length with optional ellipsis.
func Truncate(s string, max int, withEllipsis bool) string {
	if len(s) <= max {
		return s
	}
	if withEllipsis && max > 3 {
		return s[:max-3] + "..."
	}
	return s[:max]
}

// Slugify creates a URL-safe slug (lowercase, hyphens, alphanumeric).
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// RandomString generates a random alphanumeric string of n bytes.
func RandomString(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:n], nil
}

// ContainsAny checks if a string contains any of the substrings.
func ContainsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// RemoveWhitespace removes all space, tab, newline characters.
func RemoveWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// IsNumeric checks if string only contains digits.
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Capitalize capitalizes the first letter of a string.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// IsAlphaNumeric checks if string is alphanumeric
func IsAlphaNumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}

// IsSlug checks if string is URL slug friendly
func IsSlug(s string) bool {
	re := regexp.MustCompile(`^[a-z0-9\-]+$`)
	return re.MatchString(s)
}

// IsStrongPassword performs basic password strength check
func IsStrongPassword(p string) bool {
	length := len(p) >= 8
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(p)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(p)
	hasSymbol := regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(p)
	return length && hasUpper && hasNumber && hasSymbol
}
