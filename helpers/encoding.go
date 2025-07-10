package helpers

import (
	"encoding/base64"
	"encoding/hex"
)

// ToBase64 encodes string to base64
func ToBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// FromBase64 decodes base64 string
func FromBase64(encoded string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	return string(bytes), err
}

// ToHex encodes string to hex
func ToHex(data string) string {
	return hex.EncodeToString([]byte(data))
}
