package helpers

import (
	"regexp"
	"strings"
)

// IsEmailValid performs basic validation to check if the provided email string
// matches a simple email format pattern.
//
// This validation is not fully RFC compliant but is sufficient for common cases.
//
// It returns true if the email is valid, false otherwise.
//
// Example:
//
//	valid := IsEmailValid("example.user+test@gmail.com")
//	fmt.Println(valid) // Output: true
//
//	valid = IsEmailValid("invalid-email@@domain.com")
//	fmt.Println(valid) // Output: false
func IsEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(strings.ToLower(email))
}

// MaskEmail obfuscates the local part of an email address by replacing all characters
// except the first character before the '@' symbol with asterisks '*'.

// If the local part is 1 character or less, the email is returned unchanged.
//
// Example:
//
//	masked := MaskEmail("example@gmail.com")
//	fmt.Println(masked) // Output: "e******@gmail.com"
//
//	masked = MaskEmail("a@domain.com")
//	fmt.Println(masked) // Output: "a@domain.com"
func MaskEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 1 {
		return email
	}
	return email[:1] + strings.Repeat("*", at-1) + email[at:]
}

// NormalizeEmail trims leading and trailing spaces from the email string
// and converts it to lowercase.
//
// Example:
//
//	normalized := NormalizeEmail("  ExAmple@Domain.Com  ")
//	fmt.Println(normalized) // Output: "example@domain.com"
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
