package auth

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/sha3"
)

func HashNewKey(b64Key string) (b64HashedKey string, b64ServerSalt string, err error) {
	b64ServerSalt = GenerateSalt()
	b64HashedKey, err = HashKey(b64Key, b64ServerSalt)

	return b64HashedKey, b64ServerSalt, err
}

func HashKey(b64Key string, b64ServerSalt string) (string, error) {
	binaryServerSalt, err := base64.StdEncoding.DecodeString(b64ServerSalt)
	if err != nil {
		// TODO: This needs to be an internal server error
		return "", err
	}

	binaryKey, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return "", fmt.Errorf("key is not a valid base64 string")
	}

	combinedKeySalt := append(binaryKey, binaryServerSalt...)
	hashedKey := sha3.Sum512(combinedKeySalt)
	b64HashedKey := base64.StdEncoding.EncodeToString(hashedKey[:])

	return b64HashedKey, nil
}
