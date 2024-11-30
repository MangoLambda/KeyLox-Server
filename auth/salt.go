package auth

import (
	"crypto/rand"
	"encoding/base64"
)

const SaltSize = 32

func GenerateSalt() string {
	saltBytes := GenerateSaltBytes()
	return base64.StdEncoding.EncodeToString(saltBytes)
}

func GenerateSaltBytes() []byte {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		panic("failed to generate salt bytes: " + err.Error())
	}
	return salt
}
