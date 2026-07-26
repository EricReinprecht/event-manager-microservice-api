package auth

import (
	"crypto/rand"
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
