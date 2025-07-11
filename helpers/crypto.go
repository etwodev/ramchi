package helpers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the provided plaintext password using bcrypt algorithm.
//
// It returns the hashed password string and any error encountered during hashing.
//
// Example:
//
//	hash, err := HashPassword("mysecretpassword")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(hash) // Output: $2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36G1e15Ny5rQmj.LrZIvVbG
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent.
//
// Returns true if the password matches the hash, false otherwise.
//
// Example:
//
//	valid := CheckPassword("$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36G1e15Ny5rQmj.LrZIvVbG", "mysecretpassword")
//	fmt.Println(valid) // Output: true
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HMACSHA256 generates a base64-encoded HMAC using SHA-256 for the given data and secret key.
//
// It returns the resulting HMAC string.
//
// Example:
//
//	hmac := HMACSHA256("mysecretkey", "data to protect")
//	fmt.Println(hmac) // Output: "XUFAKrxLKna5cZ2REBfFkg=="
func HMACSHA256(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GenerateOTP generates a random OTP (One-Time Password) string of the specified length,
// using characters from the provided charset.
//
// If charset is empty, it defaults to "23456789ABCDEFGHJKLMNPQRSTUVWXYZ".
//
// Returns an error if length is zero or negative, or if random generation fails.
//
// Example:
//
//	otp, err := GenerateOTP(6, "")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(otp) // Output: "7G4K2H"
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
