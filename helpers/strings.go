package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"
)

// IsEmpty returns true if the input string is empty or contains only whitespace.
//
// Example:
//
//	IsEmpty("   ")  // Output: true
//	IsEmpty("text") // Output: false
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Truncate shortens the string to a maximum length `max`.
// If the string exceeds `max` and `withEllipsis` is true, it appends "..." at the end.
//
// If max is less than or equal to 3, ellipsis is ignored and string is truncated strictly.
//
// Examples:
//
//	Truncate("Hello, World", 8, true)  // Output: "Hello..."
//	Truncate("Hello", 10, true)        // Output: "Hello"
func Truncate(s string, max int, withEllipsis bool) string {
	if len(s) <= max {
		return s
	}
	if withEllipsis && max > 3 {
		return s[:max-3] + "..."
	}
	return s[:max]
}

// Slugify converts a string to a URL-friendly slug:
// lowercase, alphanumeric with words separated by hyphens.
//
// Example:
//
//	Slugify("Hello, World!")  // Output: "hello-world"
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// RandomString generates a random alphanumeric string of length n.
//
// Returns an error if the random byte generation fails.
//
// Example:
//
//	str, err := RandomString(10)
//	if err == nil {
//	    fmt.Println(str) // Output: e.g. "4f2a1c9b8e"
//	}
func RandomString(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:n], nil
}

// ContainsAny returns true if the string s contains any of the provided substrings.
//
// Example:
//
//	ContainsAny("hello world", "test", "world")  // Output: true
//	ContainsAny("hello world", "test", "abc")    // Output: false
func ContainsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// RemoveWhitespace removes all whitespace characters (spaces, tabs, newlines) from the string.
//
// Example:
//
//	RemoveWhitespace("a b \n c\t")  // Output: "abc"
func RemoveWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// IsNumeric returns true if the string contains only numeric digits (0-9).
//
// Example:
//
//	IsNumeric("12345")  // Output: true
//	IsNumeric("123a5")  // Output: false
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Capitalize converts the first character of the string to uppercase,
// leaving the rest of the string unchanged.
//
// Example:
//
//	Capitalize("hello") // Output: "Hello"
//	Capitalize("")      // Output: ""
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// IsAlphaNumeric returns true if the string contains only letters (a-z, A-Z) and digits (0-9).
//
// Example:
//
//	IsAlphaNumeric("abc123")  // Output: true
//	IsAlphaNumeric("abc_123") // Output: false
func IsAlphaNumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}

// IsSlug returns true if the string is a valid URL slug containing only lowercase letters,
// digits, and hyphens.
//
// Example:
//
//	IsSlug("my-slug-123")  // Output: true
//	IsSlug("My_Slug")      // Output: false
func IsSlug(s string) bool {
	re := regexp.MustCompile(`^[a-z0-9\-]+$`)
	return re.MatchString(s)
}

// IsStrongPassword performs a basic password strength check.
//
// The password must be at least 8 characters and contain at least one uppercase letter,
// one number, and one symbol from !@#~$%^&*()+|_
//
// Example:
//
//	IsStrongPassword("Passw0rd!") // Output: true
//	IsStrongPassword("password")  // Output: false
func IsStrongPassword(p string) bool {
	length := len(p) >= 8
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(p)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(p)
	hasSymbol := regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(p)
	return length && hasUpper && hasNumber && hasSymbol
}
