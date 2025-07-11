package helpers

import (
	"encoding/base64"
	"encoding/hex"
)

// ToBase64 encodes the input string into a base64-encoded string.
//
// Example:
//   encoded := ToBase64("hello world")
//   fmt.Println(encoded) // Output: "aGVsbG8gd29ybGQ="
func ToBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// FromBase64 decodes a base64-encoded string back to its original string form.
//
// It returns the decoded string and any error encountered during decoding.
//
// Example:
//   decoded, err := FromBase64("aGVsbG8gd29ybGQ=")
//   if err != nil {
//       // handle error
//   }
//   fmt.Println(decoded) // Output: "hello world"
func FromBase64(encoded string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	return string(bytes), err
}

// ToHex encodes the input string into its hexadecimal representation.
//
// Example:
//   hexStr := ToHex("hello")
//   fmt.Println(hexStr) // Output: "68656c6c6f"
func ToHex(data string) string {
	return hex.EncodeToString([]byte(data))
}
