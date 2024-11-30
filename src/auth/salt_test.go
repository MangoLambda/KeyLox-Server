package auth

import (
	"encoding/base64"
	"testing"
)

const SaltCount = 10000

func TestGenerateSalt(t *testing.T) {
	salt := GenerateSalt()
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		t.Fatalf("Failed to decode base64 salt: %v", err)
	}
	if len(saltBytes) != SaltSize {
		t.Errorf("Expected salt size %d, got %d", SaltSize, len(saltBytes))
	}
}

func TestGenerateSaltBytes(t *testing.T) {
	salt := GenerateSaltBytes()
	if len(salt) != SaltSize {
		t.Errorf("Expected salt size %d, got %d", SaltSize, len(salt))
	}
}

func TestGenerateSaltUniqueness(t *testing.T) {
	salts := make(map[string]struct{}, SaltCount)

	for i := 0; i < SaltCount; i++ {
		salt := GenerateSalt()
		if _, exists := salts[salt]; exists {
			t.Errorf("Duplicate salt found: %s", salt)
		}
		salts[salt] = struct{}{}
	}
}

func TestGenerateSaltBytesUniqueness(t *testing.T) {
	salts := make(map[string]struct{}, SaltCount)

	for i := 0; i < SaltCount; i++ {
		salt := GenerateSaltBytes()
		saltStr := string(salt)
		if _, exists := salts[saltStr]; exists {
			t.Errorf("Duplicate salt found: %x", salt)
		}
		salts[saltStr] = struct{}{}
	}
}
