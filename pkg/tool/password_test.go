package tool

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("s3cret-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("expected bcrypt hash, got %q", hash)
	}
	if hash == "s3cret-password" || strings.Contains(hash, "s3cret") {
		t.Fatal("hash must not contain the plaintext password")
	}

	// Same password hashes differently (random salt).
	hash2, err := HashPassword("s3cret-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == hash2 {
		t.Fatal("bcrypt should produce a different hash per call (salted)")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !CheckPassword(hash, "correct-horse") {
		t.Fatal("CheckPassword should accept the correct password")
	}
	if CheckPassword(hash, "wrong-password") {
		t.Fatal("CheckPassword should reject a wrong password")
	}
}

func TestCheckPasswordLegacyMD5(t *testing.T) {
	legacy := MD5V([]byte("old-md5-password"))

	if !IsLegacyHash(legacy) {
		t.Fatalf("32-char hex digest should be detected as legacy: %q", legacy)
	}
	if IsLegacyHash("$2a$10$abcdefghijklmnopqrstuvwxyz0123456789012345678901234567890") {
		t.Fatal("bcrypt hash should not be detected as legacy")
	}
	if IsLegacyHash("") || IsLegacyHash("tooshort") {
		t.Fatal("short/empty strings should not be detected as legacy")
	}

	if !CheckPassword(legacy, "old-md5-password") {
		t.Fatal("legacy MD5 hash should verify against the matching password")
	}
	if CheckPassword(legacy, "not-the-password") {
		t.Fatal("legacy MD5 hash should reject a wrong password")
	}
}
