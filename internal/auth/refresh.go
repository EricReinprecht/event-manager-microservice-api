package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRefreshToken() string {

	bytes := make([]byte, 64)

	_, err := rand.Read(bytes)

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(bytes)
}
