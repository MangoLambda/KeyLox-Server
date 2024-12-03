package auth

import (
	"encoding/base64"
	"testing"
)

// TODO: Fix unit tests. The stored password hash should not be directly provided in the Authorization header.

func TestVerifyCredentials_ValidCredentials(t *testing.T) {
	storedPasswordHash := base64.StdEncoding.EncodeToString([]byte("storedPasswordHash"))
	storedServerSalt := base64.StdEncoding.EncodeToString([]byte("serverSalt"))
	basicAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:storedPasswordHash"))

	valid, username, err := VerifyCredentials(storedPasswordHash, storedServerSalt, basicAuthorizationHeader)
	if !valid {
		t.Errorf("expected valid to be true, got false")
	}
	if username != "user" {
		t.Errorf("expected username to be user, got %v", username)
	}
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}

func TestVerifyCredentials_InvalidAuthorizationHeaderFormat(t *testing.T) {
	storedPasswordHash := base64.StdEncoding.EncodeToString([]byte("storedPasswordHash"))
	storedServerSalt := base64.StdEncoding.EncodeToString([]byte("serverSalt"))
	basicAuthorizationHeader := "Bearer " + base64.StdEncoding.EncodeToString([]byte("user:password"))

	valid, username, err := VerifyCredentials(storedPasswordHash, storedServerSalt, basicAuthorizationHeader)
	if valid {
		t.Errorf("expected valid to be false, got true")
	}
	if username != "" {
		t.Errorf("expected username to be empty, got %v", username)
	}
	if err == nil {
		t.Errorf("expected error to be non-nil, got nil")
	}
}

func TestVerifyCredentials_InvalidBase64Credentials(t *testing.T) {
	storedPasswordHash := base64.StdEncoding.EncodeToString([]byte("storedPasswordHash"))
	storedServerSalt := base64.StdEncoding.EncodeToString([]byte("serverSalt"))
	basicAuthorizationHeader := "Basic invalidBase64"

	valid, username, err := VerifyCredentials(storedPasswordHash, storedServerSalt, basicAuthorizationHeader)
	if valid {
		t.Errorf("expected valid to be false, got true")
	}
	if username != "" {
		t.Errorf("expected username to be empty, got %v", username)
	}
	if err == nil {
		t.Errorf("expected error to be non-nil, got nil")
	}
}

func TestVerifyCredentials_InvalidCredentialsFormat(t *testing.T) {
	storedPasswordHash := base64.StdEncoding.EncodeToString([]byte("storedPasswordHash"))
	storedServerSalt := base64.StdEncoding.EncodeToString([]byte("serverSalt"))
	basicAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("userpassword"))

	valid, username, err := VerifyCredentials(storedPasswordHash, storedServerSalt, basicAuthorizationHeader)
	if valid {
		t.Errorf("expected valid to be false, got true")
	}
	if username != "" {
		t.Errorf("expected username to be empty, got %v", username)
	}
	if err == nil {
		t.Errorf("expected error to be non-nil, got nil")
	}
}

func TestVerifyCredentials_PasswordDoesNotMatch(t *testing.T) {
	storedPasswordHash := base64.StdEncoding.EncodeToString([]byte("storedPasswordHash"))
	storedServerSalt := base64.StdEncoding.EncodeToString([]byte("serverSalt"))
	basicAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:wrongpassword"))

	valid, username, err := VerifyCredentials(storedPasswordHash, storedServerSalt, basicAuthorizationHeader)
	if valid {
		t.Errorf("expected valid to be false, got true")
	}
	if username != "user" {
		t.Errorf("expected username to be user, got %v", username)
	}
	if err == nil {
		t.Errorf("expected error to be non-nil, got nil")
	}
}

func TestVerifyCredentials_InvalidStoredPasswordHashBase64(t *testing.T) {
	storedPasswordHash := "invalidBase64"
	storedServerSalt := base64.StdEncoding.EncodeToString([]byte("serverSalt"))
	basicAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:password"))

	valid, username, err := VerifyCredentials(storedPasswordHash, storedServerSalt, basicAuthorizationHeader)
	if valid {
		t.Errorf("expected valid to be false, got true")
	}
	if username != "" {
		t.Errorf("expected username to be empty, got %v", username)
	}
	if err == nil {
		t.Errorf("expected error to be non-nil, got nil")
	}
}

func TestVerifyCredentials_InvalidStoredServerSaltBase64(t *testing.T) {
	storedPasswordHash := base64.StdEncoding.EncodeToString([]byte("storedPasswordHash"))
	storedServerSalt := "invalidBase64"
	basicAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:password"))

	valid, username, err := VerifyCredentials(storedPasswordHash, storedServerSalt, basicAuthorizationHeader)
	if valid {
		t.Errorf("expected valid to be false, got true")
	}
	if username != "" {
		t.Errorf("expected username to be empty, got %v", username)
	}
	if err == nil {
		t.Errorf("expected error to be non-nil, got nil")
	}
}
