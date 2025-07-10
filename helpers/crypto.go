package helpers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares hash with plain
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HMACSHA256 generates HMAC
func HMACSHA256(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GenerateOTP generates a random OTP string based on a given length
func GenerateOTP(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("GenerateOTP: length must be > 0")
	}

	if charset == "" {
		charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	}

	charsetLen := len(charset)
	otp := make([]byte, length)

	for i := 0; i < length; {
		b := make([]byte, 1)
		_, err := rand.Read(b)
		if err != nil {
			return "", fmt.Errorf("GenerateOTP: failed to read random byte: %w", err)
		}

		val := int(b[0])
		max := 256 - (256 % charsetLen)
		if val >= max {
			continue
		}

		otp[i] = charset[val%charsetLen]
		i++
	}

	return string(otp), nil
}
