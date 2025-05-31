package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	encodedSalt := base64.StdEncoding.EncodeToString(salt)

	return encodedSalt, nil
}
