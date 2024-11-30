package auth

import (
	"encoding/base64"
	"testing"
)

func TestHashNewKey(t *testing.T) {
	// Example base64-encoded key
	b64Key := base64.StdEncoding.EncodeToString([]byte("example_key"))

	// Call HashNewKey
	b64HashedKey, b64ServerSalt, err := HashNewKey(b64Key)
	if err != nil {
		t.Fatalf("HashNewKey failed: %v", err)
	}

	// Decode the hashed key and server salt
	hashedKey, err := base64.StdEncoding.DecodeString(b64HashedKey)
	if err != nil {
		t.Fatalf("Failed to decode base64 hashed key: %v", err)
	}

	serverSalt, err := base64.StdEncoding.DecodeString(b64ServerSalt)
	if err != nil {
		t.Fatalf("Failed to decode base64 server salt: %v", err)
	}

	// Verify the lengths of the hashed key and server salt
	if len(hashedKey) != 64 {
		t.Errorf("Expected hashed key length 64, got %d", len(hashedKey))
	}

	if len(serverSalt) != SaltSize {
		t.Errorf("Expected server salt length %d, got %d", SaltSize, len(serverSalt))
	}
}

func TestHashKey(t *testing.T) {
	// Example base64-encoded key and server salt
	b64Key := base64.StdEncoding.EncodeToString([]byte("example_key"))
	b64ServerSalt := base64.StdEncoding.EncodeToString([]byte("example_salt"))

	// Call HashKey
	b64HashedKey, err := HashKey(b64Key, b64ServerSalt)
	if err != nil {
		t.Fatalf("HashKey failed: %v", err)
	}

	// Decode the hashed key
	hashedKey, err := base64.StdEncoding.DecodeString(b64HashedKey)
	if err != nil {
		t.Fatalf("Failed to decode base64 hashed key: %v", err)
	}

	// Verify the length of the hashed key
	if len(hashedKey) != 64 {
		t.Errorf("Expected hashed key length 64, got %d", len(hashedKey))
	}
}
