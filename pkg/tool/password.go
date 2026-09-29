package tool

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// bcryptCost is the work factor for new password hashes. 12 is a good
// interactive-login trade-off (~250ms on modern hardware).
const bcryptCost = 12

// HashPassword hashes a plaintext password with bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// IsLegacyHash reports whether the stored hash is a legacy unsalted MD5 hex
// digest (32 hex chars, no bcrypt prefix).
func IsLegacyHash(hash string) bool {
	return len(hash) == 32 && !strings.HasPrefix(hash, "$2")
}

// CheckPassword verifies a plaintext password against the stored hash.
// Legacy MD5 hashes are still accepted so existing users keep access; callers
// should rehash to bcrypt after a successful verification.
func CheckPassword(hash, password string) bool {
	if IsLegacyHash(hash) {
		return MD5V([]byte(password)) == hash
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
