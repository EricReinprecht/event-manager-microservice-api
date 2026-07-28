package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateToken() string {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(bytes)
}

func HashToken(token string) string {

	hash := sha256.Sum256(
		[]byte(token),
	)

	return hex.EncodeToString(
		hash[:],
	)
}
