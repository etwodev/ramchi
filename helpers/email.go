package helpers

import (
	"regexp"
	"strings"
)

// IsEmailValid performs basic validation
// not RFC perfect, but solid
func IsEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(strings.ToLower(email))
}

// MaskEmail hides part of the email
func MaskEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 1 {
		return email
	}
	return email[:1] + strings.Repeat("*", at-1) + email[at:]
}

// NormalizeEmail trims spaces and converts to lowercase.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
