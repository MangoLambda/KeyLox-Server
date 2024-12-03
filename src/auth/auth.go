package auth

import (
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

func VerifyCredentials(storedPasswordHash string, storedServerSalt string, basicAuthorizationHeader string) (areCredentialsValid bool, username string, err error) {
	// Check if the Authorization header is in the correct format
	const basicPrefix = "Basic "
	if !strings.HasPrefix(basicAuthorizationHeader, basicPrefix) {
		return false, "", fmt.Errorf("authorization header is not in the correct format")
	}

	// Decode the base64-encoded credentials
	encodedCredentials := strings.TrimPrefix(basicAuthorizationHeader, basicPrefix)
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return false, "", fmt.Errorf("unable to decode the credentials [%v]: %v", basicAuthorizationHeader, err)
	}

	// Split the credentials into username and password
	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		return false, "", fmt.Errorf("unable to split the credentials [%v]", string(decodedCredentials))
	}

	authHeaderUsername, authHeaderPassword := credentials[0], credentials[1]

	// Decode password bytes from base64 string
	decodedAuthHeaderPassword, err := base64.StdEncoding.DecodeString(authHeaderPassword)
	if err != nil {
		return false, authHeaderUsername, fmt.Errorf("unable to decode the provided password: %v", err)
	}

	// Decode server salt bytes from base64 string
	decodedStoredServerSalt, err := base64.StdEncoding.DecodeString(storedServerSalt)
	if err != nil {
		return false, authHeaderUsername, fmt.Errorf("unable to decode the stored server salt: %v", err)
	}

	// Combine the provided password with the stored salt and hash i
	combinedKeySalt := append(decodedAuthHeaderPassword, decodedStoredServerSalt...)
	hashedKey := sha3.Sum512(combinedKeySalt)
	computedHash := base64.StdEncoding.EncodeToString(hashedKey[:])

	// Compare the computed hash with the stored password hash
	if computedHash != storedPasswordHash {
		return false, authHeaderUsername, fmt.Errorf("password does not match for user: %v", authHeaderUsername)
	}

	return true, authHeaderUsername, nil
}
